package main

import (
	"flag"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"dat520/lab3/gorumsfd"
	"dat520/lab3/leaderdetector"

	"github.com/relab/gorums"
)

func main() {
	addr := flag.String("addrs", "", "comma separated addresses of the replicas (the first address is the local replica)")
	flag.Parse()

	addrs := strings.Split(*addr, ",")
	if len(addrs) == 0 {
		flag.Usage()
		log.Fatalln("No server addresses provided")
	}

	nodeMap := make(map[string]uint32)
	for _, addr := range addrs {
		port := strings.Split(addr, ":")[1]
		id, err := strconv.Atoi(port)
		if err != nil {
			log.Fatalf("Unable to convert port to int: %v\n", err)
		}
		id &= 0x0f // Only use the last 4 bits to avoid long IDs
		nodeMap[addr] = uint32(id)
	}
	localAddr := addrs[0]
	ld := leaderdetector.NewMonLeaderDetector(nodeIDs(nodeMap))
	fd := gorumsfd.NewGorumsFailureDetector(nodeMap[localAddr], ld, time.Second)
	replica, err := NewReplica(localAddr, fd)
	if err != nil {
		log.Fatalf("Unable to create a replica with address %v: %v\n", localAddr, err)
	}
	defer replica.Stop()

	go func() {
		if err := replica.Serve(); err != nil {
			log.Fatalf("Failed to start replica: %v", err)
		}
	}()
	// Wait for this and other servers to start
	time.Sleep(1 * time.Second)

	err = replica.Start(gorums.WithNodeMap(nodeMap))
	if err != nil {
		log.Fatalf("Unable to create configuration with addresses %v: %v\n", addrs, err)
	}
	// Run until killed
	select {}
}

// nodeIDs converts a map of node IDs to a slice of node IDs.
// This is annoying, but should be fixed next year, by making the leaderdetector use uint32 or generics.
func nodeIDs(nodeMap map[string]uint32) []int {
	nodeIDs := make([]int, len(nodeMap))
	i := 0
	for _, id := range nodeMap {
		nodeIDs[i] = int(id)
		i++
	}
	// We need to sort the node IDs since map iteration order is non-deterministic.
	slices.Sort(nodeIDs)
	return nodeIDs
}
