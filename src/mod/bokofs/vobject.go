package bokofs

/*
	vObjects

	The vObjects accept and forward the request of the WebDAV server to the
	underlying workers that might be serving a file system or other
	system management functions

	vObjects and its definations shall only be used on the root layer of each
	of the bokoFS instance.
*/

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/net/webdav"
)

type vObjectProperties struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

type vObjectFileInfo struct {
	properties *vObjectProperties
	sys        interface{}
}

type vObject struct {
	properties *vObjectProperties
	parent     *RootRouter
}

// newVirtualObject creates a new virtual object
func (p *RootRouter) newVirtualObject(properties *vObjectProperties) *vObject {
	return &vObject{
		properties: properties,
		parent:     p,
	}
}

func (r *vObject) GetFileInfo() os.FileInfo {
	return &vObjectFileInfo{
		properties: r.properties,
		sys:        nil,
	}
}

/* File Info Interface */
func (r *vObjectFileInfo) IsDir() bool {
	return r.properties.isDir
}

func (r *vObjectFileInfo) ModTime() time.Time {
	return r.properties.modTime
}

func (r *vObjectFileInfo) Mode() os.FileMode {
	return r.properties.mode
}

func (r *vObjectFileInfo) Name() string {
	return r.properties.name
}

func (r *vObjectFileInfo) Size() int64 {
	return r.properties.size
}

func (r *vObjectFileInfo) Sys() interface{} {
	return r.sys
}

/* File Interface */
func (r *vObject) Close() error {
	//No need to implement this method as the vObject is not a file
	//as there will be no file descriptor opened for the vObject
	return nil
}

func (r *vObject) Read(p []byte) (n int, err error) {
	//No need to implement this method as the vObject is not a file
	//It is a virtual object that serves as a directory or a file system root
	return 0, nil
}

func (r *vObject) Readdir(count int) ([]os.FileInfo, error) {
	// Generate a emulated folder structure from worker registered paths
	fmt.Println("Readdir called")
	rootFolders, err := r.parent.parent.GetRegisteredRootFolders()
	if err != nil {
		return nil, err
	}

	// Generate the folder structure
	var folderList []os.FileInfo
	for _, folder := range rootFolders {
		thisVirtualObject := r.parent.newVirtualObject(&vObjectProperties{
			name:    folder,
			size:    0,
			mode:    os.ModeDir,
			modTime: time.Now(),
			isDir:   true,
		})

		folderList = append(folderList, thisVirtualObject.GetFileInfo())
	}
	return folderList, nil
}

func (r *vObject) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (r *vObject) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write operation not allowed: this part of the file system is read-only")
}

func (r *vObject) Stat() (os.FileInfo, error) {
	return r.GetFileInfo(), nil
}

// Ensure vObject implements the File interface
var _ webdav.File = (*vObject)(nil)
