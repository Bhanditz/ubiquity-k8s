// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.ibm.com/almaden-containers/ubiquity/remote/mounter"
)

type FakeMounter struct {
	MountStub        func(mountpoint string, volumeConfig map[string]interface{}) error
	mountMutex       sync.RWMutex
	mountArgsForCall []struct {
		mountpoint   string
		volumeConfig map[string]interface{}
	}
	mountReturns struct {
		result1 error
	}
	UnmountStub        func(volumeConfig map[string]interface{}) error
	unmountMutex       sync.RWMutex
	unmountArgsForCall []struct {
		volumeConfig map[string]interface{}
	}
	unmountReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMounter) Mount(mountpoint string, volumeConfig map[string]interface{}) error {
	fake.mountMutex.Lock()
	fake.mountArgsForCall = append(fake.mountArgsForCall, struct {
		mountpoint   string
		volumeConfig map[string]interface{}
	}{mountpoint, volumeConfig})
	fake.recordInvocation("Mount", []interface{}{mountpoint, volumeConfig})
	fake.mountMutex.Unlock()
	if fake.MountStub != nil {
		return fake.MountStub(mountpoint, volumeConfig)
	}
	return fake.mountReturns.result1
}

func (fake *FakeMounter) MountCallCount() int {
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	return len(fake.mountArgsForCall)
}

func (fake *FakeMounter) MountArgsForCall(i int) (string, map[string]interface{}) {
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	return fake.mountArgsForCall[i].mountpoint, fake.mountArgsForCall[i].volumeConfig
}

func (fake *FakeMounter) MountReturns(result1 error) {
	fake.MountStub = nil
	fake.mountReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMounter) Unmount(volumeConfig map[string]interface{}) error {
	fake.unmountMutex.Lock()
	fake.unmountArgsForCall = append(fake.unmountArgsForCall, struct {
		volumeConfig map[string]interface{}
	}{volumeConfig})
	fake.recordInvocation("Unmount", []interface{}{volumeConfig})
	fake.unmountMutex.Unlock()
	if fake.UnmountStub != nil {
		return fake.UnmountStub(volumeConfig)
	}
	return fake.unmountReturns.result1
}

func (fake *FakeMounter) UnmountCallCount() int {
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	return len(fake.unmountArgsForCall)
}

func (fake *FakeMounter) UnmountArgsForCall(i int) map[string]interface{} {
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	return fake.unmountArgsForCall[i].volumeConfig
}

func (fake *FakeMounter) UnmountReturns(result1 error) {
	fake.UnmountStub = nil
	fake.unmountReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMounter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeMounter) recordInvocation(key string, args []interface{}) {
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

var _ mounter.Mounter = new(FakeMounter)