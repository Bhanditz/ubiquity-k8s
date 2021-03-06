package controller

/**
	These tests are not in controller_test.go because they have to be inside the package of 
	the module they are testing so that it will be possible to test not exporeted functions.
**/

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/IBM/ubiquity/utils/logs"
	"github.com/IBM/ubiquity/fakes"
)



var _ = Describe("controller_internal_tests", func() {
	Context(".getK8sBaseDir", func() {
		It("should succeed if path is correct", func() {
			res, err := getK8sPodsBaseDir("/var/lib/kubelet/pods/1f94f1d9-8f36-11e8-b227-005056a4d4cb/volumes/ibm~ubiquity-k8s-flex/pvc-123")
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal("/var/lib/kubelet/pods"))
		})
		It("should succeed if path is correct and not default", func() {
			res, err := getK8sPodsBaseDir("/tmp/kubelet/pods/1f94f1d9-8f36-11e8-b227-005056a4d4cb/volumes/ibm~ubiquity-k8s-flex/pvc-123")
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal("/tmp/kubelet/pods"))
		})
		It("should fail if path is not of correct structure", func() {
			k8smountpoint := "/tmp/kubelet/soemthing/1f94f1d9-8f36-11e8-b227-005056a4d4cb/volumes/ibm~ubiquity-k8s-flex/pvc-123"
			res, err := getK8sPodsBaseDir(k8smountpoint)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(&WrongK8sDirectoryPathError{k8smountpoint}))
			Expect(res).To(Equal(""))
		})
		It("should fail if path is not of correct structure", func() {
			k8smountpoint := "/tmp/kubelet/soemthing/pods/1f94f1d9-8f36-11e8-b227-005056a4d4cb/volumes/ibm~ubiquity-k8s-flex"
			res, err := getK8sPodsBaseDir(k8smountpoint)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(&WrongK8sDirectoryPathError{k8smountpoint}))
			Expect(res).To(Equal(""))
		})
	})
	Context(".checkSlinkAlreadyExistsOnMountPoint", func() {
		var (
			fakeExecutor *fakes.FakeExecutor
			mountPoint string
		)
		BeforeEach(func() {
			fakeExecutor = new(fakes.FakeExecutor)
			mountPoint = "/tmp/kubelet/pods/1f94f1d9-8f36-11e8-b227-005056a4d4cb/volumes/ibm~ubiquity-k8s-flex/pvc-123"
		})
		It("should return no error if it is the first volume", func() {
			fakeExecutor.GetGlobFilesReturns([]string{}, nil)
			err := checkSlinkAlreadyExistsOnMountPoint("mountPoint", mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).To(BeNil())
		})
		It("should return no error if there are no other links", func() {
			fakeExecutor.GetGlobFilesReturns([]string{"/tmp/file1", "/tmp/file2"}, nil)
			fakeExecutor.IsSameFileReturns(false)
			err:= checkSlinkAlreadyExistsOnMountPoint("mountPoint",mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).To(BeNil())
		})
		It("should return an error if this mountpoint already has links", func() {
			file := "/tmp/file1"
			fakeExecutor.GetDeviceForFileStatReturnsOnCall(0, 12)
			fakeExecutor.GetDeviceForFileStatReturnsOnCall(1, 15)
			fakeExecutor.GetGlobFilesReturns([]string{file}, nil)
			fakeExecutor.IsSameFileReturns(true)
			err:= checkSlinkAlreadyExistsOnMountPoint("mountPoint", mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(&PVIsAlreadyUsedByAnotherPod{"mountPoint", []string{file}}))
		})
		It("should not return an error if this mountpoint already has links that are not mounted", func() {
			file := "/tmp/file1"
			fakeExecutor.GetDeviceForFileStatReturnsOnCall(0, 12)
			fakeExecutor.GetDeviceForFileStatReturnsOnCall(1, 12)
			fakeExecutor.GetGlobFilesReturns([]string{file}, nil)
			fakeExecutor.IsSameFileReturns(true)
			err:= checkSlinkAlreadyExistsOnMountPoint("mountPoint", mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).To(BeNil())
		})
		It("should return no errors if this mountpoint has only one links and it is the current pvc", func() {
			file := mountPoint
			fakeExecutor.GetGlobFilesReturns([]string{file}, nil)
			fakeExecutor.IsSameFileReturns(true)
			err := checkSlinkAlreadyExistsOnMountPoint("mountPoint", file, logs.GetLogger(), fakeExecutor)
			Expect(err).To(BeNil())
		})
		It("should return error if getk8sbaseDir returns an error", func() {
			k8sMountPoint := "/tmp/kubelet/something"
			err := checkSlinkAlreadyExistsOnMountPoint("mountPoint", k8sMountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).To(Equal(&WrongK8sDirectoryPathError{k8sMountPoint}))
		})
		It("should return error if glob  returns an error", func() {
			errstrObj := fmt.Errorf("An error ooccured")
			fakeExecutor.GetGlobFilesReturns(nil, errstrObj)
			err := checkSlinkAlreadyExistsOnMountPoint("mountPoint", mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).To(Equal(errstrObj))

		})
		It("should return error if stat function on the mountpoint returns an error", func() {
			errstrObj := fmt.Errorf("An error ooccured")
			fakeExecutor.GetGlobFilesReturns([]string{"/tmp/file1"}, nil)
			fakeExecutor.StatReturns(nil, errstrObj)
			err:= checkSlinkAlreadyExistsOnMountPoint("mountPoint", mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).To(Equal(errstrObj))
		})
		It("should continue if stat on a link returns an error", func() {
			errstrObj := fmt.Errorf("An error ooccured")
			fakeExecutor.GetGlobFilesReturns([]string{"/tmp/file1"}, nil)
			fakeExecutor.StatReturnsOnCall(1, nil, errstrObj)
			err := checkSlinkAlreadyExistsOnMountPoint("mountPoint", mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	Context(".checkMountPointIsMounted", func() {
		var (
			fakeExecutor *fakes.FakeExecutor
			mountPoint string
			errstrObj error
		)
		BeforeEach(func() {
			fakeExecutor = new(fakes.FakeExecutor)
			mountPoint = "/ubiquity/6001738CFC9035E8000000000091E219"
			errstrObj = fmt.Errorf("An error ooccured")
		})
		It("should return error if fails to get stat of mountpoint", func() {
			fakeExecutor.StatReturnsOnCall(0, nil, errstrObj)
			_, err := checkMountPointIsMounted(mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errstrObj))
		})
		It("should return error if fails to get lstat of mountpoint/..", func() {
			fakeExecutor.LstatReturnsOnCall(0, nil, errstrObj)
			_, err := checkMountPointIsMounted(mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errstrObj))
		})
		It("should return true if directory is mounted", func() {
			fakeExecutor.GetDeviceForFileStatReturnsOnCall(0, 12)
			fakeExecutor.GetDeviceForFileStatReturnsOnCall(1, 15)
			res, err := checkMountPointIsMounted(mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeTrue())
		})
		It("should return false if directory is not mounted", func() {
			fakeExecutor.GetDeviceForFileStatReturnsOnCall(0, 12)
			fakeExecutor.GetDeviceForFileStatReturnsOnCall(1, 12)
			res, err := checkMountPointIsMounted(mountPoint, logs.GetLogger(), fakeExecutor)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeFalse())
		})
	})
})

