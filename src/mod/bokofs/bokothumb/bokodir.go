package bokothumb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/webdav"
	"imuslab.com/bokofs/bokofsd/mod/renderer"
)

/*
	bokodir.go

	The bokodir implements a disk based file system from the webdav.FileSystem interface
	A file in this implementation corrisponding to a real file on disk
*/

type Resolutions struct {
	Width  int
	Height int
}

type RouterDir struct {
	Prefix     string //Path prefix to trim, usually is the root path of the worker
	ThumbStore string //Path to the thumbnail store
	FsPath     string //Disk path for the corrisponding file system to create thumbnail

	/* Private Properties */
	renderer *renderer.RenderHandler
	dir      webdav.Dir
}

// CreateThumbnailRenderer creates a new thumbnail renderer from a directory
func CreateThumbnailRenderer(thumbDir string, sourceFsDir string, prefix string, readonly bool) (*RouterDir, error) {
	if _, err := os.Stat(sourceFsDir); os.IsNotExist(err) {
		//Check if the sourceFsDir is a valid directory
		return nil, err
	}

	//Initiate the dir
	fs := webdav.Dir(thumbDir)

	//Create the thumbnail store if it does not exist
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return nil, err
	}

	//Create the renderer
	thumbrRenderer := renderer.NewRenderHandler()

	return &RouterDir{
		Prefix:     prefix,
		ThumbStore: thumbDir,
		FsPath:     sourceFsDir,
		renderer:   thumbrRenderer,
		dir:        fs,
	}, nil
}

func (r *RouterDir) cleanPrefix(name string) string {
	name = filepath.ToSlash(filepath.Clean(name))
	fmt.Println("[Bokothumb]", r.Prefix, name, strings.TrimPrefix(name, r.Prefix))
	return strings.TrimPrefix(name, r.Prefix)
}

func (r *RouterDir) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	// Implement the Mkdir method
	return webdav.ErrForbidden
}

func (r *RouterDir) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// Implement the OpenFile method
	name = r.cleanPrefix(name)
	ext := filepath.Ext(name)
	if ext == "" {

		//Requested a folder path. Render the content
		contents, err := os.ReadDir(filepath.Join(r.FsPath, name))
		if err != nil {
			return nil, err
		}

		//Start thumbnail rendering in background
		outputFolder := filepath.Join(r.ThumbStore, name)
		for _, entry := range contents {
			if entry.IsDir() {
				os.MkdirAll(filepath.Join(outputFolder, entry.Name()), 0755)
				continue
			}
			go func() {
				r.renderer.RenderThumbnail(filepath.Join(r.FsPath, name, entry.Name()), outputFolder)
			}()
		}
	} else {
		//Requested a file path. Render the thumbnail
		outputFolder := filepath.Join(r.ThumbStore, filepath.Dir(name))
		if err := os.MkdirAll(filepath.Dir(outputFolder), 0755); err != nil {
			return nil, err
		}
		//The RenderThumbnail function will immeidiately return if the thumbnail already exists
		//If the thumbnail is not present, it will be generated in real time
		r.renderer.RenderThumbnail(filepath.Join(r.FsPath, name), outputFolder)
	}

	fmt.Println("[Bokothumb]", "OpenFile called to "+name)
	// Check if the file is being opened with write permissions
	if flag&(os.O_WRONLY|os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC) != 0 {
		return nil, webdav.ErrForbidden
	}
	return r.dir.OpenFile(ctx, name, flag, perm)
}

func (r *RouterDir) RemoveAll(ctx context.Context, name string) error {
	// Implement the RemoveAll method
	name = r.cleanPrefix(name)
	fmt.Println("[Bokothumb]", "RemoveAll called to "+name)
	return r.dir.RemoveAll(ctx, name)
}

func (r *RouterDir) Rename(ctx context.Context, oldName, newName string) error {
	// Implement the Rename method
	return webdav.ErrForbidden
}

func (r *RouterDir) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	// Implement the Stat method
	name = r.cleanPrefix(name)
	fmt.Println("[Bokothumb]", "Stat called to "+name)
	return r.dir.Stat(ctx, name)
}

// Ensure RouterDir implements the FileSystem interface
var _ webdav.FileSystem = (*RouterDir)(nil)
