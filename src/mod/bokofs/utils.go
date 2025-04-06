package bokofs

import "imuslab.com/bokofs/bokofsd/mod/bokofs/bokoworker"

// GetRegisteredRootFolders returns all the registered root folders
// by loaded bokoFS workers. This will be shown when the client
// request the root path (/) of this bokoFS server
func (s *Server) GetRegisteredRootFolders() ([]string, error) {
	var rootFolders []string
	s.LoadedWorkers.Range(func(key, value interface{}) bool {
		thisWorker := value.(*bokoworker.Worker)
		rootFolders = append(rootFolders, thisWorker.NodeName)
		return true
	})
	return rootFolders, nil
}
