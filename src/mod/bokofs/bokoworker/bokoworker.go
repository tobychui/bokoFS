package bokoworker

import (
	"os"
	"path/filepath"
	"strings"

	"imuslab.com/bokofs/bokofsd/mod/bokofs/bokofile"
	"imuslab.com/bokofs/bokofsd/mod/bokofs/bokothumb"
)

/*
Boko Worker

A boko worker is an instance of WebDAV file server that serves a specific
disk partition or subpath in which the user can interact with the disk
via WebDAV interface
*/
type Options struct {
	NodeName       string //The node name (also the id) of the directory tree, e.g. disk1
	ServePath      string // The actual path to serve, e.g. /media/disk1/mydir
	ThumbnailStore string // The path to the thumbnail store, e.g. /media/disk1/thumbs
}

type Worker struct {
	/* Worker Properties */
	NodeName  string //The node name (also the id) of the directory tree, e.g. disk1
	ServePath string // The actual path to serve, e.g. /media/disk1/mydir

	/* Runtime Properties */
	Filesystem *bokofile.RouterDir  //The file system to serve
	Thumbnails *bokothumb.RouterDir //Thumbnail interface for this worker
}

// NewFSWorker creates a new file system worker from a directory
func NewFSWorker(options *Options) (*Worker, error) {
	nodeName := options.NodeName
	mountPath := options.ServePath
	thumbnailStore := options.ThumbnailStore

	if !strings.HasPrefix(nodeName, "/") {
		nodeName = "/" + nodeName
	}

	mountPath, _ = filepath.Abs(mountPath)
	fs, err := bokofile.CreateRouterFromDir(mountPath, nodeName, false)
	if err != nil {
		return nil, err
	}

	//Create the thumbnail store if it does not exist
	os.MkdirAll(thumbnailStore, 0755)
	thumbrender, err := bokothumb.CreateThumbnailRenderer(thumbnailStore, mountPath, nodeName, false)
	if err != nil {
		return nil, err
	}

	return &Worker{
		NodeName:  nodeName,
		ServePath: mountPath,

		Filesystem: fs,
		Thumbnails: thumbrender,
	}, nil
}
