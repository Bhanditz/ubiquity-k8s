package model

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//go:generate counterfeiter -o ../fakes/fake_storage_client.go . StorageClient
type StorageClient interface {
	Activate() error
	CreateVolume(name string, opts map[string]interface{}) error
	RemoveVolume(name string, forceDelete bool) error
	ListVolumes() ([]VolumeMetadata, error)
	GetVolume(name string) (volumeMetadata VolumeMetadata, volumeConfigDetails SpectrumConfig, err error)
	Attach(name string) (string, error)
	Detach(name string) error
	GetPluginName() string
}

type CreateRequest struct {
	Name string
	Opts map[string]interface{}
}

type RemoveRequest struct {
	Name        string
	ForceDelete bool
}

type AttachRequest struct {
	Name string
}

type DetachRequest struct {
	Name string
}

type ActivateResponse struct {
	Implements []string
}

func (r *ActivateResponse) WriteResponse(w http.ResponseWriter) {
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Errorf("Error marshalling response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(data))
}

type GenericResponse struct {
	Err string
}

func (r *GenericResponse) WriteResponse(w http.ResponseWriter) {
	if r.Err != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Errorf("Error marshalling response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(data))
}

type GenericRequest struct {
	Name string
}

type MountResponse struct {
	Mountpoint string
	Err        string
}

func (r *MountResponse) WriteResponse(w http.ResponseWriter) {
	if r.Err != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Errorf("Error marshalling Get response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(data))
}

type VolumeMetadata struct {
	Name       string
	Mountpoint string
}

type SpectrumConfig struct {
	FilesetId  string `json:"fileset"`
	Filesystem string `json:"filesystem"`
}
type GetResponse struct {
	Volume VolumeMetadata
	Err    string
	Config SpectrumConfig
}

func (r *GetResponse) WriteResponse(w http.ResponseWriter) {
	if r.Err != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Errorf("Error marshalling Get response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(data))
}

type ListResponse struct {
	Volumes []VolumeMetadata
	Err     string
}

func (r *ListResponse) WriteResponse(w http.ResponseWriter) {
	if r.Err != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Errorf("Error marshalling Get response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(data))
}

type FlexVolumeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Device  string `json:"device"`
}

type FlexVolumeMountRequest struct {
	MountPath   string                 `json:"mountPath"`
	MountDevice string                 `json:"name"`
	Opts        map[string]interface{} `json:"opts"`
}

type FlexVolumeAttachRequest struct {
	VolumeId   string `json:"volumeID"`
	Filesystem string `json:"filesystem"`
	Size       string `json:"size"`
	Path       string `json:"path"`
	Fileset    string `json:"fileset"`
}

type FlexVolumeUnmountRequest struct {
	MountPath string `json:"mountPath"`
}

type FlexVolumeDetachRequest struct {
Name string `json:"name"`
}
