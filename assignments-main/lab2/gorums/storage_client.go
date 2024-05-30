package gorums

type StorageClient struct {
	// TODO: Add fields if necessary
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string) *StorageClient {
	//TODO(student): Implement NewStorageClient
	return &StorageClient{}
}

// Writes the provided value to a random server
func (sc *StorageClient) WriteValue(value string) error {
	//TODO(student): Implement WriteValue
	return nil
}

// Returns a slice of values stored on all servers
func (sc *StorageClient) ReadValues() ([]string, error) {
	//TODO(student): Implement ReadValues
	return []string{}, nil
}
