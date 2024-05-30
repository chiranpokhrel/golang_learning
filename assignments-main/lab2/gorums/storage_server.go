package gorums

import (
	"sync"
)

// The storage server should implement the server interface defined in the protobuf files
type StorageServer struct {
	sync.RWMutex
	data []string
	// TODO: Add fields if necessary

}

// Creates a new StorageServer.
func NewStorageServer() *StorageServer {
	//TODO(student): implement NewStorageServer
	return &StorageServer{}
}

// Start the server listening on the provided address string
// The function should be non-blocking
// Returns the full listening address of the server as string
// Hint: Use go routine to start the server.
func (s *StorageServer) StartServer(addr string) string {
	//TODO(student): implement StartServer
	return ""
}

// Returns the data slice on this server
func (s *StorageServer) GetData() []string {
	s.RLock()
	defer s.RUnlock()
	return s.data
}

// Sets the data slice to a value
func (s *StorageServer) SetData(data []string) {
	s.Lock()
	defer s.Unlock()
	s.data = data
}
