package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	paxos "dat520/lab5/gorumspaxos"
	pb "dat520/lab5/gorumspaxos/proto"

	gorums "github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		srvAddrs      = flag.String("addrs", "", "server addresses separated by ','")
		clientRequest = flag.String("clientRequest", "", "client requests separated by ','")
		clientId      = flag.String("clientId", "", "Client Id, different for each client")
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

	clientRequests := strings.Split(*clientRequest, ",")
	if len(clientRequests) == 0 {
		log.Fatalln("no client requests are provided")
	}
	// start a initial proposer
	ClientStart(addrs, clientRequests, clientId)
}

// ClientStart creates the configuration with the list of replicas addresses, which are read from the
// command line. From the list of clientRequests, send each request to the configuration and
// wait for the reply. Upon receiving the reply send the next request.
func ClientStart(addrs []string, clientRequests []string, clientId *string) {
	log.Printf("Connecting to %d Paxos replicas: %v", len(addrs), addrs)
	config, mgr := createConfiguration(addrs)
	defer mgr.Close()
	for index, request := range clientRequests {
		req := pb.Value{ClientID: *clientId, ClientSeq: uint32(index), ClientCommand: request}
		resp := doSendRequest(config, &req)
		log.Printf("response: %v\t for the client request: %v", resp, &req)
	}
}

// Internal: doSendRequest can send requests to paxos servers by quorum call and
// for the response from the quorum function.
func doSendRequest(config *pb.Configuration, value *pb.Value) *pb.Response {
	waitTimeForRequest := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeForRequest)
	defer cancel()
	resp, err := config.ClientHandle(ctx, value)
	if err != nil {
		log.Fatalf("ClientHandle quorum call error: %v", err)
	}
	if resp == nil {
		log.Println("response is nil")
	}
	return resp
}

// createConfiguration creates the gorums configuration with the list of addresses.
func createConfiguration(addrs []string) (configuration *pb.Configuration, manager *pb.Manager) {
	mgr := pb.NewManager(gorums.WithDialTimeout(5*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(), // block until connections are made
			grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
		),
	)
	qspec := paxos.NewPaxosQSpec(len(addrs))
	config, err := mgr.NewConfiguration(qspec, gorums.WithNodeList(addrs))
	if err != nil {
		log.Fatalf("Error in forming the configuration: %v\n", err)
		return nil, nil
	}
	return config, mgr
}
