package main

import (
	"flag"
	"hash/fnv"
	"log"
	"net"
	"os"
	"strings"

	paxos "dat520/lab5/gorumspaxos"
)

func main() {
	var (
		localAddr = flag.String("laddr", "localhost:8080", "local address to listen on")
		srvAddrs  = flag.String("addrs", "", "all other remaining replica addresses separated by ','")
	)
	flag.Usage = func() {
		log.Printf("Usage: %s [OPTIONS]\nOptions:", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	addrs := strings.Split(*srvAddrs, ",")
	if len(addrs) == 0 {
		log.Fatalln("no server addresses provided")
	}
	l, err := net.Listen("tcp", *localAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Printf("Waiting for requests at %s", l.Addr().String())
	nodeMap := make(map[string]uint32)
	addrs = append(addrs, *localAddr)
	for _, addr := range addrs {
		id := calculateHash(addr)
		nodeMap[addr] = uint32(id)
	}
	replica := paxos.NewPaxosReplica(calculateHash(*localAddr), nodeMap)
	replica.Serve(l)
}

// calculateHash calculates an integer hash for the address of the node
func calculateHash(address string) int {
	h := fnv.New32a()
	h.Write([]byte(address))
	return int(h.Sum32())
}
