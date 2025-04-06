package bokofile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/webdav"
)

/*
	bokodir.go

	The bokodir implements a disk based file system from the webdav.FileSystem interface
	A file in this implementation corrisponding to a real file on disk
*/

type RouterDir struct {
	Prefix   string //Path prefix to trim, usually is the root path of the worker
	DiskPath string //Disk path to create a file system from
	ReadOnly bool   //Label this worker as read only

	/* Private Properties */
	dir webdav.Dir
}

// Create a routerdir from a directory
func CreateRouterFromDir(dir string, prefix string, readonly bool) (*RouterDir, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, err
	}

	//Initiate the dir
	fs := webdav.Dir(dir)

	return &RouterDir{
		Prefix:   prefix,
		DiskPath: dir,
		ReadOnly: readonly,
		dir:      fs,
	}, nil
}

func (r *RouterDir) cleanPrefix(name string) string {
	name = filepath.ToSlash(filepath.Clean(name)) + "/"
	fmt.Println("[Bokodir]", r.Prefix, name, strings.TrimPrefix(name, r.Prefix))
	return strings.TrimPrefix(name, r.Prefix)
}

func (r *RouterDir) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	// Implement the Mkdir method
	name = r.cleanPrefix(name)
	fmt.Println("[Bokodir]", "Mkdir called to "+name)
	return r.dir.Mkdir(ctx, name, perm)
}

func (r *RouterDir) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// Implement the OpenFile method
	name = r.cleanPrefix(name)
	fmt.Println("[Bokodir]", "OpenFile called to "+name)
	return r.dir.OpenFile(ctx, name, flag, perm)
}

func (r *RouterDir) RemoveAll(ctx context.Context, name string) error {
	// Implement the RemoveAll method
	name = r.cleanPrefix(name)
	fmt.Println("[Bokodir]", "RemoveAll called to "+name)
	return r.dir.RemoveAll(ctx, name)
}

func (r *RouterDir) Rename(ctx context.Context, oldName, newName string) error {
	// Implement the Rename method
	oldName = r.cleanPrefix(oldName)
	newName = r.cleanPrefix(newName)
	fmt.Println("[Bokodir]", "Rename called from "+oldName+" to "+newName)
	return r.dir.Rename(ctx, oldName, newName)
}

func (r *RouterDir) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	// Implement the Stat method
	name = r.cleanPrefix(name)
	fmt.Println("[Bokodir]", "Stat called to "+name)
	return r.dir.Stat(ctx, name)
}

// Ensure RouterDir implements the FileSystem interface
var _ webdav.FileSystem = (*RouterDir)(nil)
