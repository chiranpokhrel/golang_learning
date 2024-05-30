package gorumspaxos

import "testing"

func TestMyIndex(t *testing.T) {
	tests := []struct {
		name     string
		nodeMap  map[string]uint32
		myID     int
		expected int
	}{
		{name: "Single node", nodeMap: map[string]uint32{"node1": 1}, myID: 1, expected: 0},
		{name: "Multiple nodes, target present", nodeMap: map[string]uint32{"node1": 1, "node2": 2, "node3": 3}, myID: 2, expected: 1},
		{name: "Multiple nodes, target absent", nodeMap: map[string]uint32{"node1": 1, "node2": 2, "node3": 3}, myID: 4, expected: -1},
		{name: "Unsorted input", nodeMap: map[string]uint32{"node1": 3, "node2": 1, "node3": 2}, myID: 3, expected: 2},
		{name: "Empty map", nodeMap: map[string]uint32{}, myID: 1, expected: -1},
		{name: "5 nodes, target 5", nodeMap: map[string]uint32{"node1": 1, "node2": 2, "node3": 3, "node4": 4, "node5": 5}, myID: 5, expected: 4},
		{name: "5 nodes, target 4", nodeMap: map[string]uint32{"node1": 1, "node2": 2, "node3": 3, "node4": 4, "node5": 5}, myID: 4, expected: 3},
		{name: "5 nodes, target 3", nodeMap: map[string]uint32{"node1": 1, "node2": 2, "node3": 3, "node4": 4, "node5": 5}, myID: 3, expected: 2},
		{name: "5 nodes, target 2", nodeMap: map[string]uint32{"node1": 1, "node2": 2, "node3": 3, "node4": 4, "node5": 5}, myID: 2, expected: 1},
		{name: "5 nodes, target 1", nodeMap: map[string]uint32{"node1": 1, "node2": 2, "node3": 3, "node4": 4, "node5": 5}, myID: 1, expected: 0},
		{name: "5 nodes, target absent", nodeMap: map[string]uint32{"node1": 1, "node2": 2, "node3": 3, "node4": 4, "node5": 5}, myID: 6, expected: -1},
		{name: "5 nodes, large ids, target absent", nodeMap: map[string]uint32{"node1": 12345, "node2": 23456, "node3": 34567, "node4": 45678, "node5": 56789}, myID: 123456, expected: -1},
		{name: "5 nodes, large ids, target  12345", nodeMap: map[string]uint32{"node1": 12345, "node2": 23456, "node3": 34567, "node4": 45678, "node5": 56789}, myID: 12345, expected: 0},
		{name: "5 nodes, large ids, target  23456", nodeMap: map[string]uint32{"node1": 12345, "node2": 23456, "node3": 34567, "node4": 45678, "node5": 56789}, myID: 23456, expected: 1},
		{name: "5 nodes, large ids, target  34567", nodeMap: map[string]uint32{"node1": 12345, "node2": 23456, "node3": 34567, "node4": 45678, "node5": 56789}, myID: 34567, expected: 2},
		{name: "5 nodes, large ids, target  45678", nodeMap: map[string]uint32{"node1": 12345, "node2": 23456, "node3": 34567, "node4": 45678, "node5": 56789}, myID: 45678, expected: 3},
		{name: "5 nodes, large ids, target  56789", nodeMap: map[string]uint32{"node1": 12345, "node2": 23456, "node3": 34567, "node4": 45678, "node5": 56789}, myID: 56789, expected: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := myIndex(tt.myID, tt.nodeMap)
			if got != tt.expected {
				t.Errorf("myIndex(%d, %v) = %d; want %d", tt.myID, tt.nodeMap, got, tt.expected)
			}
		})
	}
}
