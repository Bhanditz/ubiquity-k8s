// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.ibm.com/almaden-containers/spectrum-common.git/core"
	"github.ibm.com/almaden-containers/spectrum-common.git/models"
)

type FakeSpectrumClient struct {
	CreateStub        func(name string, opts map[string]interface{}) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		name string
		opts map[string]interface{}
	}
	createReturns struct {
		result1 error
	}
	RemoveStub        func(name string) error
	removeMutex       sync.RWMutex
	removeArgsForCall []struct {
		name string
	}
	removeReturns struct {
		result1 error
	}
	AttachStub        func(name string) (string, error)
	attachMutex       sync.RWMutex
	attachArgsForCall []struct {
		name string
	}
	attachReturns struct {
		result1 string
		result2 error
	}
	DetachStub        func(name string) error
	detachMutex       sync.RWMutex
	detachArgsForCall []struct {
		name string
	}
	detachReturns struct {
		result1 error
	}
	ListStub        func() ([]models.VolumeMetadata, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct{}
	listReturns     struct {
		result1 []models.VolumeMetadata
		result2 error
	}
	GetStub        func(name string) (*models.VolumeMetadata, *models.SpectrumConfig, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		name string
	}
	getReturns struct {
		result1 *models.VolumeMetadata
		result2 *models.SpectrumConfig
		result3 error
	}
	IsMountedStub        func() (bool, error)
	isMountedMutex       sync.RWMutex
	isMountedArgsForCall []struct{}
	isMountedReturns     struct {
		result1 bool
		result2 error
	}
	MountStub        func() error
	mountMutex       sync.RWMutex
	mountArgsForCall []struct{}
	mountReturns     struct {
		result1 error
	}
	GetFileSetForMountPointStub        func(mountPoint string) (string, error)
	getFileSetForMountPointMutex       sync.RWMutex
	getFileSetForMountPointArgsForCall []struct {
		mountPoint string
	}
	getFileSetForMountPointReturns struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSpectrumClient) Create(name string, opts map[string]interface{}) error {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		name string
		opts map[string]interface{}
	}{name, opts})
	fake.recordInvocation("Create", []interface{}{name, opts})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(name, opts)
	} else {
		return fake.createReturns.result1
	}
}

func (fake *FakeSpectrumClient) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeSpectrumClient) CreateArgsForCall(i int) (string, map[string]interface{}) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].name, fake.createArgsForCall[i].opts
}

func (fake *FakeSpectrumClient) CreateReturns(result1 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSpectrumClient) Remove(name string) error {
	fake.removeMutex.Lock()
	fake.removeArgsForCall = append(fake.removeArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("Remove", []interface{}{name})
	fake.removeMutex.Unlock()
	if fake.RemoveStub != nil {
		return fake.RemoveStub(name)
	} else {
		return fake.removeReturns.result1
	}
}

func (fake *FakeSpectrumClient) RemoveCallCount() int {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return len(fake.removeArgsForCall)
}

func (fake *FakeSpectrumClient) RemoveArgsForCall(i int) string {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return fake.removeArgsForCall[i].name
}

func (fake *FakeSpectrumClient) RemoveReturns(result1 error) {
	fake.RemoveStub = nil
	fake.removeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSpectrumClient) Attach(name string) (string, error) {
	fake.attachMutex.Lock()
	fake.attachArgsForCall = append(fake.attachArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("Attach", []interface{}{name})
	fake.attachMutex.Unlock()
	if fake.AttachStub != nil {
		return fake.AttachStub(name)
	} else {
		return fake.attachReturns.result1, fake.attachReturns.result2
	}
}

func (fake *FakeSpectrumClient) AttachCallCount() int {
	fake.attachMutex.RLock()
	defer fake.attachMutex.RUnlock()
	return len(fake.attachArgsForCall)
}

func (fake *FakeSpectrumClient) AttachArgsForCall(i int) string {
	fake.attachMutex.RLock()
	defer fake.attachMutex.RUnlock()
	return fake.attachArgsForCall[i].name
}

func (fake *FakeSpectrumClient) AttachReturns(result1 string, result2 error) {
	fake.AttachStub = nil
	fake.attachReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeSpectrumClient) Detach(name string) error {
	fake.detachMutex.Lock()
	fake.detachArgsForCall = append(fake.detachArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("Detach", []interface{}{name})
	fake.detachMutex.Unlock()
	if fake.DetachStub != nil {
		return fake.DetachStub(name)
	} else {
		return fake.detachReturns.result1
	}
}

func (fake *FakeSpectrumClient) DetachCallCount() int {
	fake.detachMutex.RLock()
	defer fake.detachMutex.RUnlock()
	return len(fake.detachArgsForCall)
}

func (fake *FakeSpectrumClient) DetachArgsForCall(i int) string {
	fake.detachMutex.RLock()
	defer fake.detachMutex.RUnlock()
	return fake.detachArgsForCall[i].name
}

func (fake *FakeSpectrumClient) DetachReturns(result1 error) {
	fake.DetachStub = nil
	fake.detachReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSpectrumClient) List() ([]models.VolumeMetadata, error) {
	fake.listMutex.Lock()
	fake.listArgsForCall = append(fake.listArgsForCall, struct{}{})
	fake.recordInvocation("List", []interface{}{})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub()
	} else {
		return fake.listReturns.result1, fake.listReturns.result2
	}
}

func (fake *FakeSpectrumClient) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeSpectrumClient) ListReturns(result1 []models.VolumeMetadata, result2 error) {
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 []models.VolumeMetadata
		result2 error
	}{result1, result2}
}

func (fake *FakeSpectrumClient) Get(name string) (*models.VolumeMetadata, *models.SpectrumConfig, error) {
	fake.getMutex.Lock()
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("Get", []interface{}{name})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(name)
	} else {
		return fake.getReturns.result1, fake.getReturns.result2, fake.getReturns.result3
	}
}

func (fake *FakeSpectrumClient) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeSpectrumClient) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].name
}

func (fake *FakeSpectrumClient) GetReturns(result1 *models.VolumeMetadata, result2 *models.SpectrumConfig, result3 error) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 *models.VolumeMetadata
		result2 *models.SpectrumConfig
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeSpectrumClient) IsMounted() (bool, error) {
	fake.isMountedMutex.Lock()
	fake.isMountedArgsForCall = append(fake.isMountedArgsForCall, struct{}{})
	fake.recordInvocation("IsMounted", []interface{}{})
	fake.isMountedMutex.Unlock()
	if fake.IsMountedStub != nil {
		return fake.IsMountedStub()
	} else {
		return fake.isMountedReturns.result1, fake.isMountedReturns.result2
	}
}

func (fake *FakeSpectrumClient) IsMountedCallCount() int {
	fake.isMountedMutex.RLock()
	defer fake.isMountedMutex.RUnlock()
	return len(fake.isMountedArgsForCall)
}

func (fake *FakeSpectrumClient) IsMountedReturns(result1 bool, result2 error) {
	fake.IsMountedStub = nil
	fake.isMountedReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeSpectrumClient) Mount() error {
	fake.mountMutex.Lock()
	fake.mountArgsForCall = append(fake.mountArgsForCall, struct{}{})
	fake.recordInvocation("Mount", []interface{}{})
	fake.mountMutex.Unlock()
	if fake.MountStub != nil {
		return fake.MountStub()
	} else {
		return fake.mountReturns.result1
	}
}

func (fake *FakeSpectrumClient) MountCallCount() int {
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	return len(fake.mountArgsForCall)
}

func (fake *FakeSpectrumClient) MountReturns(result1 error) {
	fake.MountStub = nil
	fake.mountReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSpectrumClient) GetFileSetForMountPoint(mountPoint string) (string, error) {
	fake.getFileSetForMountPointMutex.Lock()
	fake.getFileSetForMountPointArgsForCall = append(fake.getFileSetForMountPointArgsForCall, struct {
		mountPoint string
	}{mountPoint})
	fake.recordInvocation("GetFileSetForMountPoint", []interface{}{mountPoint})
	fake.getFileSetForMountPointMutex.Unlock()
	if fake.GetFileSetForMountPointStub != nil {
		return fake.GetFileSetForMountPointStub(mountPoint)
	} else {
		return fake.getFileSetForMountPointReturns.result1, fake.getFileSetForMountPointReturns.result2
	}
}

func (fake *FakeSpectrumClient) GetFileSetForMountPointCallCount() int {
	fake.getFileSetForMountPointMutex.RLock()
	defer fake.getFileSetForMountPointMutex.RUnlock()
	return len(fake.getFileSetForMountPointArgsForCall)
}

func (fake *FakeSpectrumClient) GetFileSetForMountPointArgsForCall(i int) string {
	fake.getFileSetForMountPointMutex.RLock()
	defer fake.getFileSetForMountPointMutex.RUnlock()
	return fake.getFileSetForMountPointArgsForCall[i].mountPoint
}

func (fake *FakeSpectrumClient) GetFileSetForMountPointReturns(result1 string, result2 error) {
	fake.GetFileSetForMountPointStub = nil
	fake.getFileSetForMountPointReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeSpectrumClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	fake.attachMutex.RLock()
	defer fake.attachMutex.RUnlock()
	fake.detachMutex.RLock()
	defer fake.detachMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.isMountedMutex.RLock()
	defer fake.isMountedMutex.RUnlock()
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	fake.getFileSetForMountPointMutex.RLock()
	defer fake.getFileSetForMountPointMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeSpectrumClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ core.SpectrumClient = new(FakeSpectrumClient)
