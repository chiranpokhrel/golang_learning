# Remote Procedure Calls (RPCs) with gRPC

A popular way to design distributed applications is by means of remote procedure calls.
A remote procedure call allows a client application to invoke a server's procedure or method almost as if the server's method was local to the client application.
The main difference between local procedure calls and remote procedure calls is that clients (and servers) need to perform some setup before procedures can be called, and clients must be equipped to handle errors due to network delays or disconnection.

gRPC is an open source framework for working with remote procedure calls that can interact across different languages and operating systems.
Before getting started you should read through the following documents:

- [Introduction to gRPC](https://grpc.io/docs/what-is-grpc/introduction/):
  This gives a brief overview of how gRPC works and how to work with protocol buffers.
- [gRPC Quick Start (Go)](https://grpc.io/docs/languages/go/quickstart/):
  Guides you through setting up the gRPC working environment in Go and installing the `protoc` protocol buffer compiler.
- [gRPC Basics Tutorial (Go)](https://grpc.io/docs/languages/go/basics/):
  This document contains several important pieces of information:

  - Overview of service definitions and code generation.
  - Distinction between simple RPCs, server-side streaming RPCs, client-side streaming RPCs and bidirectional streaming RPCs.
  - Creating and starting the server.
    We have already set up the skeleton code for the server implementation [here](server/server.go).
  - Creating the client.

- [Protocol buffers](https://developers.google.com/protocol-buffers/) and [Protocol Buffer Basics: Go](https://developers.google.com/protocol-buffers/docs/gotutorial):
  You should get familiar with what protocol buffers are and how to implement them.
  These links introduce protocol buffers, why they are useful, how to define and compile them, and how to use the generated message types in Go.

Other useful resources:

- [API (`package grpc`)](https://pkg.go.dev/google.golang.org/grpc?tab=doc):
  This is the API of the `grpc` package in Go, which is used by the generated code as well as the portions of your application that interact with gRPC.
- [Language Guide (proto3)](https://developers.google.com/protocol-buffers/docs/proto3)
  Detailed specification of the `proto3` language syntax.

## Installing the Protobuf Compiler and Plugins

Please read the [this](proto.md) documentation to install the protobuf compiler and gRPC plugins necessary for this lab.

## Building your own RPC-based key-value storage

This repository contains code for a very simple key-value storage
service, where the keys and values are strings.
It offers the following three methods:

- `Insert`: inserts a key-value pair. It returns a bool indicating success or failure.
- `Lookup`: returns the value of the given key.
- `Keys`: returns a slice of all the keys in the storage.

Please look at, **but do not change**, the `kv.pb.go` file.
This file is generated using `protoc` and contains important APIs and message definitions needed for these exercises.

```go
// Messages
type InsertRequest struct {
    Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
    Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

type InsertResponse struct {
    Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

type LookupRequest struct {
    Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

type LookupResponse struct {
    Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

type KeysRequest struct {
}

type KeysResponse struct {
    Keys []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
}

// Methods
// Client API for KeyValueService service
type KeyValueServiceClient interface {
	Insert(ctx context.Context, in *InsertRequest, opts ...grpc.CallOption) (*InsertResponse, error)
	Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error)
	Keys(ctx context.Context, in *KeysRequest, opts ...grpc.CallOption) (*KeysResponse, error)
}

// Server API for KeyValueService service
type KeyValueServiceServer interface {
	Insert(context.Context, *InsertRequest) (*InsertResponse, error)
	Lookup(context.Context, *LookupRequest) (*LookupResponse, error)
	Keys(context.Context, *KeysRequest) (*KeysResponse, error)
	mustEmbedUnimplementedKeyValueServiceServer()
}
```

## Tasks

1. Create a client and connect to the server.

2. In the client, call the `Insert()` gRPC for a number of key/value pairs.

3. In the server, implement the `Lookup()` gRPC, which should return the value of the requested key.

4. In the server, implement the `Keys()` gRPC, which should return a slice of the keys
   (not the values) of the map back to the client.

5. In the client, call the `Lookup()` gRPC for each of the key/value pairs inserted and verify
   the result returned from the `Lookup()` gRPC matches the value inserted for the corresponding key.

6. In the client, call the `Keys()` gRPC and verify that the number of keys returned matches
   the expected number.

7. Several clients connecting to the server may read and write concurrently from the shared key-value map.
   This will eventually cause inconsistencies in the map, unless some form of protection is instituted.
   Implement locking at the appropriate locations in the code. See [pkg/sync](https://pkg.go.dev/sync).

### Extras

- Explain why the clients can access the map at the server concurrently.
- If you run your server without protection on the map, are you able to provoke inconsistencies in the map?

### Troubleshooting

In case of compiler errors related to `kv.pb.go`, it may help to recompile the proto file.
Please refer to [this](proto.md#Compiling_.proto_Files) documentation to see how to compile the proto file.
