package main

import (
	"context"
	"net"
	"time"

	"dat520/lab3/gorumsfd"
	pb "dat520/lab3/gorumsfd/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	managerDialTimeout = 5 * time.Second
)

type Replica struct {
	// server portion of the replica
	lis net.Listener
	srv *gorums.Server

	// client portion of the replica
	mgr *pb.Manager
	cfg *pb.Configuration
	fd  gorumsfd.FailureDetector
}

// NewReplica creates a new replica with the given address and failure detector
// implementation. It returns an error if the server cannot be started.
func NewReplica(addr string, fd gorumsfd.FailureDetector) (*Replica, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := gorums.NewServer()
	pb.RegisterFailureDetectorServer(srv, fd)
	return &Replica{
		lis: lis,
		srv: srv,
		fd:  fd,
	}, nil
}

// Serve starts the server and blocks until the server is stopped.
// It returns an error if the server cannot be started.
// This should be called in a goroutine and call log.Fatal if it returns an error.
func (r *Replica) Serve() error {
	return r.srv.Serve(r.lis)
}

// Stop stops the replica, closing the server and the manager if it exists.
func (r *Replica) Stop() {
	r.srv.Stop()
	r.fd.Stop()
}

// Addr returns the address of the server.
func (r *Replica) Addr() string {
	return r.lis.Addr().String()
}

// Start starts the replica with the given nodes. It returns an error if the
// configuration cannot be created.
func (r *Replica) Start(nodes gorums.NodeListOption) error {
	r.mgr = pb.NewManager(
		gorums.WithDialTimeout(managerDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
		),
	)
	cfg, err := r.mgr.NewConfiguration(nodes)
	if err != nil {
		return err
	}
	r.cfg = cfg
	hbSender := func(hb *pb.HeartBeat) {
		r.cfg.Heartbeat(context.Background(), hb)
	}
	r.fd.Start(hbSender)
	return nil
}

func (r *Replica) Configuration() *pb.Configuration {
	return r.cfg
}
