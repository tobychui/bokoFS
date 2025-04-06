package bokofs

import (
	"context"
	"os"

	"golang.org/x/net/webdav"
)

/*
FlowRouter

This interface is used to define the flow of the file system
*/
type FlowRouter interface {
	Mkdir(ctx context.Context, name string, perm os.FileMode) error
	OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error)
	RemoveAll(ctx context.Context, name string) error
	Rename(ctx context.Context, oldName, newName string) error
	Stat(ctx context.Context, name string) (os.FileInfo, error)
}
