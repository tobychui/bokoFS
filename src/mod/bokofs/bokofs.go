package bokofs

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/webdav"
	"imuslab.com/bokofs/bokofsd/mod/bokofs/bokoworker"
)

/*
	WebDAV module

	This module handle the interfacing to the underlying file system
	through the middle ware
*/

type Server struct {
	LoadedWorkers sync.Map   //Storing uuid to bokoworker pointer (*bokoworker.Worker)
	FsRouter      FlowRouter //The file system router
	ThumbRouter   FlowRouter //The thumbnail router
	fsprefix      string
	thumbprefix   string
}

/* NewWebdavInterfaceServer creates a new WebDAV server instance */
func NewWebdavInterfaceServer(fsPrefix string, thumbPrefix string) (*Server, error) {
	//Make sure the prefix has a prefix and a trailing slash
	if fsPrefix == "" || thumbPrefix == "" {
		return nil, os.ErrInvalid
	}

	if !strings.HasPrefix(fsPrefix, "/") {
		fsPrefix = "/" + fsPrefix
	}

	if !strings.HasSuffix(fsPrefix, "/") {
		fsPrefix = fsPrefix + "/"
	}

	if !strings.HasPrefix(thumbPrefix, "/") {
		thumbPrefix = "/" + thumbPrefix
	}

	if !strings.HasSuffix(thumbPrefix, "/") {
		thumbPrefix = thumbPrefix + "/"
	}

	thisServer := Server{
		LoadedWorkers: sync.Map{},
		fsprefix:      fsPrefix,
		thumbprefix:   thumbPrefix,
	}

	//Initiate the root router file system
	fsRouter, err := NewRootRouter(&thisServer, fsPrefix, RouterType_FS)
	if err != nil {
		return nil, err
	}
	thisServer.FsRouter = fsRouter

	thumbRouter, err := NewRootRouter(&thisServer, thumbPrefix, RouterType_Thumb)
	if err != nil {
		return nil, err
	}
	thisServer.ThumbRouter = thumbRouter
	return &thisServer, nil
}

func (s *Server) AddWorker(worker *bokoworker.Worker) error {
	if worker.Filesystem == nil || worker.Thumbnails == nil {
		return errors.New("missing resources router")
	}

	//Check if the worker root path is already loaded
	if _, ok := s.LoadedWorkers.Load(worker.NodeName); ok {
		return os.ErrExist
	}
	s.LoadedWorkers.Store(worker.NodeName, worker)
	return nil
}

func (s *Server) RemoveWorker(workerRootPath string) {
	s.LoadedWorkers.Delete(workerRootPath)
}

func (s *Server) FsHandler() http.Handler {
	srv := &webdav.Handler{
		FileSystem: s.FsRouter,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
			}
		},
	}
	return srv
}

func (s *Server) ThumbHandler() http.Handler {
	srv := &webdav.Handler{
		FileSystem: s.ThumbRouter,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("THUMB [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				log.Printf("THUMB [%s]: %s \n", r.Method, r.URL)
			}
		},
	}
	return srv
}
