package bokofs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/webdav"
	"imuslab.com/bokofs/bokofsd/mod/bokofs/bokoworker"
)

type RootRouter struct {
	pathPrefix string
	routerType RouterType
	parent     *Server
}

type RouterType int

const (
	RouterType_FS RouterType = iota
	RouterType_Thumb
)

type RouterDirHandler interface {
	Mkdir(ctx context.Context, name string, perm os.FileMode) error
	OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error)
	RemoveAll(ctx context.Context, name string) error
	Rename(ctx context.Context, oldName, newName string) error
	Stat(ctx context.Context, name string) (os.FileInfo, error)
}

// RootRouter implements the webdav.FileSystem interface
// It serves as the root of the file system, routing requests to the appropriate worker
func NewRootRouter(p *Server, prefix string, rType RouterType) (*RootRouter, error) {
	return &RootRouter{
		pathPrefix: prefix,
		routerType: rType,
		parent:     p,
	}, nil
}

/*
	Router Internal Implementation
*/

// fixpath fix the path to be relative to the root path of this router
func (r *RootRouter) fixpath(name string) string {
	if name == r.pathPrefix || name == "" {
		return "/"
	}
	//Trim off the prefix path
	name = strings.TrimPrefix(name, r.pathPrefix)
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}
	return name
}

// getRootDir returns the root directory of the request
func (r *RootRouter) getRootDir(name string) string {
	if name == "" || name == "/" {
		return "/"
	}

	name = filepath.ToSlash(filepath.Clean(name))
	pathChunks := strings.Split(name, "/")
	reqRootPath := "/" + pathChunks[1]
	fmt.Println("Requesting Root Path: ", reqRootPath)
	name = strings.TrimSuffix(reqRootPath, "/")
	return name
}

// GetFileSystemFromWorker returns the file system from the worker
func (r *RootRouter) getWorkerByPath(name string) (*bokoworker.Worker, error) {
	reqRootPath := r.getRootDir(name)
	targetWorker, ok := r.parent.LoadedWorkers.Load(reqRootPath)
	if !ok {
		return nil, os.ErrNotExist
	}

	return targetWorker.(*bokoworker.Worker), nil
}

// getFileSystemFromWorker returns the file system from the worker
func (r *RootRouter) getFileSystemFromWorker(worker *bokoworker.Worker) (RouterDirHandler, error) {
	if r.routerType == RouterType_FS {
		return worker.Filesystem, nil
	} else if r.routerType == RouterType_Thumb {
		return worker.Thumbnails, nil
	}
	return nil, errors.New("Invalid router type")
}

/*
	WebDAV FileSystem Interface Implementation
*/

func (r *RootRouter) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	// Implement the Mkdir method
	name = r.fixpath(name)
	fmt.Println("Mkdir called to " + name)
	return nil
}

func (r *RootRouter) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// Implement the OpenFile method
	name = r.fixpath(name)
	fmt.Println("OpenFile called to " + name)
	if filepath.ToSlash(filepath.Base(name)) == "/" {
		//Request to the vObject base path
		thisVirtualObject := r.newVirtualObject(&vObjectProperties{
			name:    name,
			size:    0,
			mode:    os.ModeDir,
			modTime: time.Now(),
			isDir:   true,
		})

		return thisVirtualObject, nil
	}

	targetWorker, err := r.getWorkerByPath(name)
	if err != nil {
		return nil, err
	}

	targetFileSystem, err := r.getFileSystemFromWorker(targetWorker)
	if err != nil {
		return nil, err
	}
	return targetFileSystem.OpenFile(ctx, name, flag, perm)
}

func (r *RootRouter) RemoveAll(ctx context.Context, name string) error {
	// Implement the RemoveAll method
	name = r.fixpath(name)
	fmt.Println("RemoveAll called to " + name)
	return nil
}

func (r *RootRouter) Rename(ctx context.Context, oldName, newName string) error {
	// Implement the Rename method
	oldName = r.fixpath(oldName)
	newName = r.fixpath(newName)
	fmt.Println("Rename called from " + oldName + " to " + newName)
	return nil
}

func (r *RootRouter) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	// Implement the Stat method
	name = r.fixpath(name)
	fmt.Println("Stat called to " + name)
	if filepath.ToSlash(filepath.Base(name)) == "/" {
		//Create an emulated file system to serve the mounted workers
		thisVirtualObject := r.newVirtualObject(&vObjectProperties{
			name:    name,
			size:    0,
			mode:    os.ModeDir,
			modTime: time.Now(),
			isDir:   true,
		})

		thisVirtualObjectFileInfo := thisVirtualObject.GetFileInfo()

		return thisVirtualObjectFileInfo, nil
	}

	//Load the target worker from the path
	targetWorker, err := r.getWorkerByPath(name)
	fmt.Println("Target Worker: ", targetWorker, name)
	if err != nil {
		return nil, err
	}

	targetFileSystem, err := r.getFileSystemFromWorker(targetWorker)
	if err != nil {
		return nil, err
	}

	return targetFileSystem.Stat(ctx, name)
}

// Ensure RootRouter implements the FileSystem interface
var _ webdav.FileSystem = (*RootRouter)(nil)
