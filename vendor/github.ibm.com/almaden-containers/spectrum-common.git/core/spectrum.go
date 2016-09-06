package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.ibm.com/almaden-containers/spectrum-common.git/models"
)

//go:generate counterfeiter -o ../fakes/fake_spectrum_client.go . SpectrumClient

type SpectrumClient interface {
	Activate() error
	Create(name string, opts map[string]interface{}) error
	CreateWithoutProvisioning(name string, opts map[string]interface{}) error
	Remove(name string) error
	Attach(name string) (string, error)
	Detach(name string) error
	ExportNfs(name string, clientCIDR string) (string, error)
	UnexportNfs(name string) error
	List() ([]models.VolumeMetadata, error)
	Get(name string) (*models.VolumeMetadata, *models.SpectrumConfig, error)
	IsMounted() (bool, error)
	Mount() error
	RemoveWithoutDeletingVolume(string) error
	GetFileSetForMountPoint(mountPoint string) (string, error)
}

type Fileset struct {
	Name             string
	Mountpoint       string
	DockerVolumeName string
}

type MappingConfig struct {
	Mappings map[string]Fileset
}

func NewSpectrumClient(logger *log.Logger, filesystem, mountpoint string, dbclient *DatabaseClient) SpectrumClient {
	return &MMCliFilesetClient{log: logger, Filesystem: filesystem, Mountpoint: mountpoint, DbClient: dbclient,
		filelock: NewFileLock(logger, filesystem, mountpoint)}
}

type MMCliFilesetClient struct {
	Filesystem  string
	Mountpoint  string
	log         *log.Logger
	DbClient    *DatabaseClient
	isMounted   bool
	isActivated bool
	filelock    *FileLock
}

func (m *MMCliFilesetClient) Activate() (err error) {
	m.log.Println("MMCliFilesetClient: Activate start")
	defer m.log.Println("MMCliFilesetClient: Activate end")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	if m.isActivated {
		return nil
	}

	clusterId, err := getClusterId()

	if err != nil {
		m.log.Println(err.Error())
		return err
	}

	if len(clusterId) == 0 {
		clusterIdErr := errors.New("Unable to retrieve clusterId: clusterId is empty")
		m.log.Println(clusterIdErr.Error())
		return clusterIdErr
	}

	m.DbClient.ClusterId = clusterId

	err = m.DbClient.CreateVolumeTable()

	if err != nil {
		m.log.Println(err.Error())
		return err
	}
	m.isActivated = true
	return nil
}

func (m *MMCliFilesetClient) Create(name string, opts map[string]interface{}) (err error) {
	m.log.Println("MMCliFilesetClient: create start")
	defer m.log.Println("MMCliFilesetClient: create end")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	volExists, err := m.DbClient.VolumeExists(name)

	if err != nil {
		m.log.Println(err.Error())
		return err
	}

	if volExists {
		return fmt.Errorf("Volume already exists")
	}

	userSpecifiedFileset, exists := opts["fileset"]
	if exists == true {
		return m.updateDBWithExistingFileset(name, userSpecifiedFileset.(string))
	} else {
		return m.create(name, opts)
	}

}

func (m *MMCliFilesetClient) CreateWithoutProvisioning(name string, opts map[string]interface{}) (err error) {
	m.log.Println("MMCliFilesetClient: CreateWithoutProvisioning start")
	defer m.log.Println("MMCliFilesetClient: createWithoutProvisioning end")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	volExists, err := m.DbClient.VolumeExists(name)

	if err != nil {
		m.log.Println(err.Error())
		return err
	}

	if volExists {
		return fmt.Errorf("Volume already exists")
	}
	userSpecifiedFileset, exists := opts["fileset"]
	if exists == true {
		return m.updateDBWithExistingFileset(name, userSpecifiedFileset.(string))
	} else {

		err := m.filesetExists(name)
		if err != nil {
			m.log.Printf("Fileset not found %#v", err)
			return err
		}

		err = m.DbClient.InsertFilesetVolume(userSpecifiedFileset.(string), name)

		if err != nil {
			m.log.Printf("Error persisting mapping %#v", err)
			return err
		}

	}
	return nil
}

func (m *MMCliFilesetClient) filesetExists(name string) error {
	m.log.Println("MMCliFilesetClient:  fileset exists start")
	defer m.log.Println("MMCliFilesetClient: fileset exists end")
	m.log.Printf("filesetExists: %s\n", name)

	spectrumCommand := "/usr/lpp/mmfs/bin/mmlsfileset"
	args := []string{m.Filesystem, name, "-Y"}
	cmd := exec.Command(spectrumCommand, args...)

	_, err := cmd.Output()
	if err != nil {
		m.log.Printf("error checking fileset %#v", err)
		return err
	}
	var line string
	scanner := bufio.NewScanner(cmd.Stdin)
	for scanner.Scan() {
		line = (scanner.Text())
		lineSlice := strings.Split(line, " ")
		if lineSlice[0] == name {
			return nil
		}
	}
	m.log.Println("fileset not found")
	return fmt.Errorf("volume not found in the filesystem")
}

func (m *MMCliFilesetClient) updateDBWithExistingFileset(name, userSpecifiedFileset string) error {
	m.log.Println("MMCliFilesetClient:  updateDBWithExistingFileset start")
	defer m.log.Println("MMCliFilesetClient: updateDBWithExistingFileset end")
	m.log.Printf("User specified fileset: %s\n", userSpecifiedFileset)

	spectrumCommand := "/usr/lpp/mmfs/bin/mmlsfileset"
	args := []string{m.Filesystem, userSpecifiedFileset, "-Y"}
	cmd := exec.Command(spectrumCommand, args...)
	_, err := cmd.Output()
	if err != nil {
		m.log.Println(err)
		return err
	}

	err = m.DbClient.InsertFilesetVolume(userSpecifiedFileset, name)

	if err != nil {
		m.log.Println(err.Error())
		return err
	}
	return nil
}

func (m *MMCliFilesetClient) updateMappingWithExistingFileset(name, userSpecifiedFileset string, mappingConfig MappingConfig) error {
	m.log.Println("MMCliFilesetClient:  updateMappingWithExistingFileset start")
	defer m.log.Println("MMCliFilesetClient: updateMappingWithExistingFileset end")
	m.log.Printf("User specified fileset: %s\n", userSpecifiedFileset)

	spectrumCommand := "/usr/lpp/mmfs/bin/mmlsfileset"
	args := []string{m.Filesystem, userSpecifiedFileset, "-Y"}
	cmd := exec.Command(spectrumCommand, args...)
	_, err := cmd.Output()
	if err != nil {
		m.log.Printf("error updating mapping with existing fileset %#v", err)
		return err
	}
	mappingConfig.Mappings[name] = Fileset{Name: userSpecifiedFileset, DockerVolumeName: name}
	// persist mapping config
	err = m.persistMappingConfig(mappingConfig)
	if err != nil {
		return err
	}
	return nil
}

func (m *MMCliFilesetClient) create(name string, opts map[string]interface{}) error {
	m.log.Println("MMCliFilesetClient: createNew start")
	defer m.log.Println("MMCliFilesetClient: createNew end")

	filesetName := generateFilesetName()
	m.log.Printf("creating a new fileset: %s\n", filesetName)
	// create fileset
	spectrumCommand := "/usr/lpp/mmfs/bin/mmcrfileset"
	args := []string{m.Filesystem, filesetName}
	cmd := exec.Command(spectrumCommand, args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Failed to create fileset")
	}
	m.log.Printf("Createfileset output: %s\n", string(output))

	err = m.DbClient.InsertFilesetVolume(filesetName, name)

	if err != nil {
		m.log.Println(err.Error())
		return err
	}
	return nil
}

func (m *MMCliFilesetClient) Remove(name string) (err error) {
	m.log.Println("MMCliFilesetClient: remove start")
	defer m.log.Println("MMCliFilesetClient: remove end")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	volExists, err := m.DbClient.VolumeExists(name)

	if err != nil {
		m.log.Println(err.Error())
		return err
	}

	if volExists {

		existingVolume, err := m.DbClient.GetVolume(name)

		if err != nil {
			m.log.Println(err.Error())
			return err
		}

		spectrumCommand := "/usr/lpp/mmfs/bin/mmdelfileset"
		args := []string{m.Filesystem, existingVolume.Fileset, "-f"}
		cmd := exec.Command(spectrumCommand, args...)
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("Failed to remove fileset")
		}
		m.log.Printf("MMCliFilesetClient: Deletefileset output: %s\n", string(output))

		err = m.DbClient.DeleteVolume(name)

		if err != nil {
			m.log.Println(err.Error())
			return err
		}
	}
	return nil
}

func (m *MMCliFilesetClient) RemoveWithoutDeletingVolume(name string) error {
	m.log.Println("MMCliFilesetClient: RemoveWithoutDeletingVolume start")
	defer m.log.Println("MMCliFilesetClient: RemoveWithoutDeletingVolume end")
	err := m.DbClient.DeleteVolume(name)
	if err != nil {
		m.log.Printf("error retrieving mapping %#v", err)
		return err
	}
	return nil
}

func (m *MMCliFilesetClient) Attach(name string) (Mountpoint string, err error) {
	m.log.Println("MMCliFilesetClient: attach start")
	defer m.log.Println("MMCliFilesetClient: attach end")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	volExists, err := m.DbClient.VolumeExists(name)

	if err != nil {
		m.log.Println(err.Error())
		return "", err
	}

	if !volExists {
		return "", fmt.Errorf("fileset couldn't be located")
	}

	existingVolume, err := m.DbClient.GetVolume(name)

	if err != nil {
		m.log.Println(err.Error())
		return "", err
	}

	if existingVolume.Mountpoint != "" {
		Mountpoint = existingVolume.Mountpoint
		return Mountpoint, nil
	}

	spectrumCommand := "/usr/lpp/mmfs/bin/mmlinkfileset"
	filesetPath := path.Join(m.Mountpoint, existingVolume.Fileset)
	args := []string{m.Filesystem, existingVolume.Fileset, "-J", filesetPath}
	cmd := exec.Command(spectrumCommand, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to link fileset")
	}
	m.log.Printf("MMCliFilesetClient: Linkfileset output: %s\n", string(output))

	//hack for now
	args = []string{"-R", "777", filesetPath}
	cmd = exec.Command("chmod", args...)
	output, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to set permissions for fileset")
	}

	err = m.DbClient.UpdateVolumeMountpoint(name, filesetPath)

	if err != nil {
		m.log.Println(err.Error())
		return "", fmt.Errorf("internal error updating mapping")
	}

	Mountpoint = filesetPath
	return Mountpoint, nil
}

func (m *MMCliFilesetClient) Detach(name string) (err error) {
	m.log.Println("MMCliFilesetClient: detach start")
	defer m.log.Println("MMCliFilesetClient: detach end")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	volExists, err := m.DbClient.VolumeExists(name)

	if err != nil {
		m.log.Println(err.Error())
		return err
	}

	if !volExists {
		return fmt.Errorf("fileset couldn't be located")
	}

	existingVolume, err := m.DbClient.GetVolume(name)

	if err != nil {
		m.log.Println(err.Error())
		return err
	}

	if existingVolume.Mountpoint == "" {
		return fmt.Errorf("fileset not linked")
	}

	spectrumCommand := "/usr/lpp/mmfs/bin/mmunlinkfileset"
	args := []string{m.Filesystem, existingVolume.Fileset}
	cmd := exec.Command(spectrumCommand, args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Failed to unlink fileset")
	}
	m.log.Printf("MMCliFilesetClient: unLinkfileset output: %s\n", string(output))

	err = m.DbClient.UpdateVolumeMountpoint(name, "")

	if err != nil {
		m.log.Println(err.Error())
		return fmt.Errorf("internal error updating mapping")
	}
	return nil
}

func (m *MMCliFilesetClient) ExportNfs(name string, clientCIDR string) (string, error) {
	m.log.Println("MMCliFilesetClient: ExportNfs start")
	defer m.log.Println("MMCliFilesetClient: ExportNfs end")
	mappingConfig, err := m.retrieveMappingConfig()
	if err != nil {
		return "", err
	}
	mapping, ok := mappingConfig.Mappings[name]
	if ok == false {
		m.log.Println("MMCliFilesetClient ExportNfs: fileset not found")
		return "", fmt.Errorf("fileset couldn't be located")
	}
	if mapping.Mountpoint == "" {
		m.log.Println("MMCliFilesetClient ExportNfs: fileset not linked")
		return "", fmt.Errorf("fileset not linked")
	}
	spectrumCommand := "/usr/lpp/mmfs/bin/mmnfs"
	filesetPath := path.Join(m.Mountpoint, mapping.Name)
	args := []string{"export", "add", filesetPath, "--client", fmt.Sprintf("%s(Access_Type=RW,Protocols=3:4,Squash=no_root_squash)", clientCIDR)}
	cmd := exec.Command(spectrumCommand, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to export fileset via NFS: %s", err.Error())
	}
	m.log.Printf("MMCliFilesetClient: ExportNfs output: %s\n", string(output))
	return filesetPath, nil
}

func (m *MMCliFilesetClient) UnexportNfs(name string) error {
	m.log.Println("MMCliFilesetClient: UnexportNfs start")
	defer m.log.Println("MMCliFilesetClient: UnexportNfs end")
	mappingConfig, err := m.retrieveMappingConfig()
	if err != nil {
		return err
	}
	mapping, ok := mappingConfig.Mappings[name]
	if ok == false {
		m.log.Println("MMCliFilesetClient UnexportNfs: fileset not found")
		return fmt.Errorf("fileset couldn't be located")
	}
	spectrumCommand := "/usr/lpp/mmfs/bin/mmnfs"
	filesetPath := path.Join(m.Mountpoint, mapping.Name)
	args := []string{"export", "remove", filesetPath, "--force"}
	cmd := exec.Command(spectrumCommand, args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Failed to unexport fileset via NFS: %s", err.Error())
	}
	m.log.Printf("MMCliFilesetClient: UnexportNfs output: %s\n", string(output))

	mapping.Mountpoint = ""
	mappingConfig.Mappings[name] = mapping
	err = m.persistMappingConfig(mappingConfig)
	if err != nil {
		return fmt.Errorf("internal error updating mapping")
	}
	return nil
}

func (m *MMCliFilesetClient) List() (volumeList []models.VolumeMetadata, err error) {
	m.log.Println("MMCliFilesetClient: list start")
	defer m.log.Println("MMCliFilesetClient: list end")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	volumesInDb, err := m.DbClient.ListVolumes()

	if err != nil {
		m.log.Println(err.Error())
		return nil, err
	}

	var volumes []models.VolumeMetadata
	for _, volume := range volumesInDb {
		volumes = append(volumes, models.VolumeMetadata{Name: volume.VolumeName, Mountpoint: volume.Mountpoint})
	}
	volumeList = volumes
	return volumeList, nil
}

func (m *MMCliFilesetClient) Get(name string) (volumeMetadata *models.VolumeMetadata, volumeConfigDetails *models.SpectrumConfig, err error) {
	m.log.Println("MMCliFilesetClient: get start")
	defer m.log.Println("MMCliFilesetClient: get finish")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	volExists, err := m.DbClient.VolumeExists(name)

	if err != nil {
		m.log.Println(err.Error())
		return nil, nil, err
	}

	if volExists {

		existingVolume, err := m.DbClient.GetVolume(name)

		if err != nil {
			m.log.Println(err.Error())
			return nil, nil, err
		}

		volumeMetadata = &models.VolumeMetadata{Name: existingVolume.VolumeName, Mountpoint: existingVolume.Mountpoint}
		volumeConfigDetails = &models.SpectrumConfig{FilesetId: existingVolume.Fileset, Filesystem: m.Filesystem}
		return volumeMetadata, volumeConfigDetails, nil
	}
	return nil, nil, fmt.Errorf("Cannot find info")
}

func (m *MMCliFilesetClient) IsMounted() (isMounted bool, err error) {
	m.log.Println("MMCliFilesetClient: isMounted start")
	defer m.log.Println("MMCliFilesetClient: isMounted end")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	if m.isMounted == true {
		isMounted = true
		return isMounted, nil
	}

	spectrumCommand := "/usr/lpp/mmfs/bin/mmlsmount"
	args := []string{m.Filesystem, "-L", "-Y"}
	cmd := exec.Command(spectrumCommand, args...)
	outputBytes, err := cmd.Output()
	if err != nil {
		m.log.Printf("Error running command\n")
		m.log.Println(err)
		return false, err
	}
	mountedNodes := extractMountedNodes(string(outputBytes))
	if len(mountedNodes) == 0 {
		//not mounted anywhere
		isMounted = false
		return isMounted, nil
	} else {
		// checkif mounted on current node -- compare node name
		currentNode, _ := os.Hostname()
		m.log.Printf("MMCliFilesetClient: node name: %s\n", currentNode)
		for _, node := range mountedNodes {
			if node == currentNode {
				m.isMounted = true
				isMounted = true
				return isMounted, nil
			}
		}
	}
	isMounted = false
	return isMounted, nil
}

func (m *MMCliFilesetClient) Mount() (err error) {
	m.log.Println("MMCliFilesetClient: mount start")
	defer m.log.Println("MMCliFilesetClient: mount end")

	m.filelock.Lock()
	defer func() {
		lockErr := m.filelock.Unlock()
		if lockErr != nil && err == nil {
			err = lockErr
		}
	}()

	if m.isMounted == true {
		return nil
	}

	spectrumCommand := "/usr/lpp/mmfs/bin/mmmount"
	args := []string{m.Filesystem, m.Mountpoint}
	cmd := exec.Command(spectrumCommand, args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Failed to mount filesystem")
	}
	m.log.Println(output)
	m.isMounted = true
	return nil
}

func extractMountedNodes(spectrumOutput string) []string {
	var nodes []string
	lines := strings.Split(spectrumOutput, "\n")
	if len(lines) == 1 {
		return nodes
	}
	for _, line := range lines[1:] {
		tokens := strings.Split(line, ":")
		if len(tokens) > 10 {
			if tokens[11] != "" {
				nodes = append(nodes, tokens[11])
			}
		}
	}
	return nodes
}

func getClusterId() (string, error) {

	var clusterId string

	spectrumCommand := "/usr/lpp/mmfs/bin/mmlscluster"
	cmd := exec.Command(spectrumCommand)
	outputBytes, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error running command: %s", err.Error())
	}
	spectrumOutput := string(outputBytes)

	lines := strings.Split(spectrumOutput, "\n")
	tokens := strings.Split(lines[4], ":")

	if len(tokens) == 2 {
		if strings.TrimSpace(tokens[0]) == "GPFS cluster id" {
			clusterId = strings.TrimSpace(tokens[1])
		}
	}
	return clusterId, nil
}

func (m *MMCliFilesetClient) retrieveMappingConfig() (MappingConfig, error) {
	m.log.Println("MMCliFilesetClient: retrieveMappingConfig start")
	defer m.log.Println("MMCliFilesetClient: retrieveMappingConfig end")
	mappingFile, err := os.Open(path.Join(m.Mountpoint, ".docker.json"))
	if err != nil {
		m.log.Println(err.Error())
		if os.IsNotExist(err) == true {
			m.log.Println("file does not exist")
			mappingConfig := MappingConfig{Mappings: map[string]Fileset{}}
			err = m.persistMappingConfig(mappingConfig)
			if err != nil {
				return MappingConfig{}, fmt.Errorf("error initializing config file (%s)", err.Error())
			}
			return mappingConfig, nil
		} else {
			return MappingConfig{}, fmt.Errorf("error opening config file (%s)", err.Error())
		}
	}
	jsonParser := json.NewDecoder(mappingFile)
	var mappingConfig MappingConfig
	if err = jsonParser.Decode(&mappingConfig); err != nil {
		return MappingConfig{}, fmt.Errorf("error parsing config file (%s)", err.Error())
	}
	return mappingConfig, nil
}

func (m *MMCliFilesetClient) GetFileSetForMountPoint(mountPoint string) (string, error) {

	volume, err := m.DbClient.GetVolumeForMountPoint(mountPoint)

	if err != nil {
		m.log.Println(err.Error())
		return "", err
	}
	return volume, nil
}

func (m *MMCliFilesetClient) persistMappingConfig(mappingConfig MappingConfig) error {
	m.log.Println("MMCliFilesetClient: persisteMappingConfig start")
	defer m.log.Println("MMCliFilesetClient: persisteMappingConfig end")
	data, err := json.Marshal(&mappingConfig)
	if err != nil {
		return fmt.Errorf("Error marshalling mapping config to file: %s", err.Error())
	}
	err = ioutil.WriteFile(path.Join(m.Mountpoint, ".docker.json"), data, 0644)
	if err != nil {
		return fmt.Errorf("Error writing json spec: %s", err.Error())
	}
	return nil
}
func generateFilesetName() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
