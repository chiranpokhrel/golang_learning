package gorums

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStorageService(t *testing.T) {
	numServers := 3
	addresses := []string{}
	servers := []*StorageServer{}
	for i := 0; i < numServers; i++ {
		srv := NewStorageServer()
		addr := srv.StartServer("")
		if diff := cmp.Diff(srv.GetData(), []string{}); diff != "" {
			t.Errorf("Unexpected data stored on server: %v, (-want +got):\n%s", addr, diff)
		}
		addresses = append(addresses, addr)
		servers = append(servers, srv)
	}

	client := NewStorageClient(addresses)

	for _, values := range storageServiceTest {
		for _, value := range values.writes {
			client.WriteValue(value)
		}

		resp, err := client.ReadValues()
		if err != nil {
			t.Errorf("Unable to read values: %v", err)
		}

		// Sort both slices, to ensure that they have the same order
		sort.Slice(resp, func(i, j int) bool {
			return resp[i] < resp[j]
		})
		sort.Slice(values.writes, func(i, j int) bool {
			return values.writes[i] < values.writes[j]
		})
		if diff := cmp.Diff(resp, values.writes); diff != "" {
			t.Errorf("Unexpected returned value:(-want +got):\n%s", diff)
		}
		for _, srv := range servers {
			srv.SetData([]string{})
		}
	}
}

func TestStorageRead(t *testing.T) {
	addresses := []string{}
	servers := []*StorageServer{}
	for i := 0; i < 3; i++ {
		srv := NewStorageServer()
		addr := srv.StartServer("")
		if diff := cmp.Diff(srv.GetData(), []string{}); diff != "" {
			t.Errorf("Unexpected data stored on server: %v:(-want +got):\n%s", addr, diff)
		}
		addresses = append(addresses, addr)
		servers = append(servers, srv)
	}

	client := NewStorageClient(addresses)

	for _, test := range storageReadTest {
		for i, val := range test.startData {
			servers[i].SetData(val)
		}

		resp, err := client.ReadValues()
		if err != nil {
			t.Errorf("Unable to read values: %v", err)
		}

		// Sort both slices, to ensure that they have the same order
		sort.Slice(resp, func(i, j int) bool {
			return resp[i] < resp[j]
		})
		sort.Slice(test.expected, func(i, j int) bool {
			return test.expected[i] < test.expected[j]
		})
		if diff := cmp.Diff(resp, test.expected); diff != "" {
			t.Errorf("Unexpected returned value:(-want +got):\n%s", diff)
		}
	}
}

var storageServiceTest = []struct {
	writes []string
}{
	{
		[]string{},
	},
	{
		[]string{""},
	},
	{
		[]string{"Value 1", "Value 2", "Value 3", "value 4", "Value 5"},
	},
	{
		[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
	},
}

var storageReadTest = []struct {
	startData [][]string
	expected  []string
}{
	{
		[][]string{{}, {}, {}},
		[]string{},
	},
	{
		[][]string{
			{"Value 1"},
			{"Value 2"},
			{"Value 3"},
		},
		[]string{"Value 1", "Value 2", "Value 3"},
	},
	{
		[][]string{
			{"Value 1", "Value 2", "Value 4"},
			{"Value 2"},
			{"Value 3"},
		},
		[]string{"Value 1", "Value 2", "Value 2", "Value 3", "Value 4"},
	},
}
