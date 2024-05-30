# Gorums

Gorums is a framework for building fault tolerant distributed systems.
The two main features provided by gorums are an abstraction for managing what processes are a part of a system, called the configuration, and an abstraction for managing replies received from a quorum call.
Gorums is built on top of gRPC and uses protobuf to define quorum calls.
The gorums repository can be found [here](https://github.com/relab/gorums).
A guide on how to install gorums and a user guide can be found [here](https://github.com/relab/gorums/blob/master/doc/user-guide.md).
This guide will take examples from the user guide.

## Managing the configuration

Gorums uses `configuration` types to manage which processes are a part of a system.
To create a configuration you first need to create a manager.
The manager maintains the connections to the nodes and you can pass several `grpc.DialOptions` to change the options used for the connections.
The code snippet below creates a gorums manager with a dial timeout of 500 milliseconds and with the following `grpc.DialOptions`:

- [`WithBlock()`](https://pkg.go.dev/google.golang.org/grpc#WithBlock) which blocks until the connection is up and
- [`WithTransportCredentials()`](https://pkg.go.dev/google.golang.org/grpc#WithTransportCredentials) which set the credentials to be used, in this case no credentials are used.

The manager provides some methods for managing the nodes, details can be found [here](https://github.com/relab/gorums/blob/master/mgr.go).

```go
mgr := NewManager(
   gorums.WithDialTimeout(500*time.Millisecond),
   gorums.WithGrpcDialOptions(
         grpc.WithBlock(),
         grpc.WithTransportCredentials(insecure.NewCredentials(),
      ),
   ),
)
```

The manager can then be used to create a new configuration using the `NewConfiguration` method.
The nodes in the configuration can be specified in several ways, for example with the `gorums.WithNodeList(addrs)` option which takes a list of addresses of the nodes.
Other options also exists, some examples and how to use them can be found [here](https://github.com/relab/gorums/blob/master/doc/user-guide.md#working-with-configurations).
The `QSpec` mentioned here is used in quorum calls and can be ignored for now.
The code snippet below creates a slice of three addresses and creates a new configuration with these three servers.

```go
// Get all all available node ids, 3 nodes
addrs := []string{
   "127.0.0.1:8080",
   "127.0.0.1:8081",
   "127.0.0.1:8082",
}
// Create a configuration including all nodes
allNodesConfig, err := mgr.NewConfiguration(gorums.WithNodeList(addrs))
if err != nil {
   log.Fatal("Failed to create config:", err)
}
```

The configuration type provides several ways of reading information about the configuration, details can be found [here](https://github.com/relab/gorums/blob/master/config.go).
When a configuration has been created these methods can be used to, for example, iterate over all the nodes and perform RPC calls on each of them, in a similar way to what we did in [gRPC](/lab2/grpc/):

```go
for _, node := range allNodesConfig.Nodes() {
   reply, err := node.Write(context.Background(), &WriteMessage)
   if err != nil {
      log.Fatal("Write call returned error:", err)
   }
}
```

However the main strength of gorums lies in its ability to manage requests, and the corresponding responses, to the entire configuration simultaneously.
These are called quorum calls and we will use them in the next section.

## Quorum calls

Quorum calls are a way for gorums to send multiple RPC calls in parallel to all nodes in the configuration.
It then uses a quorum function to manage the received responses and return a single response based on a defined function.
This makes them useful when you only require responses from a subset of the nodes.
For example, when you only require the same response from a majority of the nodes.
An other use case is when you want to select which one of several responses to use, e.g., you only want the most recently updated data.

To create a quorum call you first need to define an RPC service method as a quorum call in the protobuf definition.
Then you need to implement the servers that provide a response, and finally define a `QSpec` type to handle the incoming responses.

### Defining quorum calls in protobuf

To define quorum calls in protobuf the `gorums.proto` package must first be imported in the protobuf definition:

```protobuf
import "gorums.proto";
```

After that a service and its RPCs can be defined similarly as described in [protobuf](/lab2/protobuf/).
To define an RPC as a quorum call the `gorums.quorumcall` option must be specified.

```protobuf
service QCStorage {
   rpc Read(ReadRequest) returns (State) {
      option (gorums.quorumcall) = true;
   }
   rpc Write(State) returns (WriteResponse) {
      option (gorums.quorumcall) = true;
   }
}
```

To compile a protobuf file use the following command.
Replace \<filename\> with the full name of the .proto file:

```shell
protoc -I=$(go list -m -f {{.Dir}} github.com/relab/gorums):. \
  --go_out=paths=source_relative:. \
  --gorums_out=paths=source_relative:. \
  <filename>
```

### Creating the gorums server

To create a gorums server, a type implementing the interface defined in the compiled protobuf file must be implemented.
In the example above, that interface looks like this:

```go
type QCStorage interface {
   Read(ctx context.Context, in *ReadRequest) (*State, error)
   Write(ctx context.Context, in *State) (*WriteResponse, error)
}
```

This is similar to how we implement gRPC servers, and is described in [gRPC](/lab2/grpc/).
For some more information about how to implement the server interface see [here](https://github.com/relab/gorums/blob/master/doc/user-guide.md#implementing-the-storageserver).

Starting a gorums server is also similar to what you did in the [gRPC](/lab2/grpc/) task.
In addition a `gorums.Server` type must be created and registered like this:

```go
lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
if err != nil {
   log.Fatal(err)
}
gorumsSrv := gorums.NewServer()
srv := storageSrv{state: &State{}}
RegisterStorageServer(gorumsSrv, &srv)
gorumsSrv.Serve(lis)
```

Remember that `gorumsSrv.Serve(lis)` is a blocking call.

### Implementing the QuorumSpec interface

When creating a quorum call a `QuorumSpec` interface is also created.
The `QuorumSpec` interface contains a quorum function for each quorum call.
A quorum function has the same name as its quorum call, but with the `QF` added to the end.
The quorum function takes the received responses and returns a single response based on these.
The `QuorumSpec` interface for the service defined above looks like this:

```go
type QuorumSpec interface {
  // ReadQF is the quorum function for the Read quorum call method.
  ReadQF(req *ReadRequest, replies map[uint32]*State) (*State, bool)

  // WriteQF is the quorum function for the Write quorum call method.
  WriteQF(req *WriteRequest, replies map[uint32]*WriteResponse) (*WriteResponse, bool)
}
```

The quorum functions are called several times, with varying number of replies.
The `replies` map always contains all the replies received at this point, which means that some replies will be processed several times.
The quorum function returns two values.
The first is the response it has decided on based on the responses from the clients.
This can be any response, e.g., the first received response or the response returned by the largest number of clients.
The second value is a boolean.
The value should be `true` if the quorum function has decided on a value, and should be `false` if the quorum function has not decided on a value.
The code snippet below shows an implementation of the `ReadQF` quorum function where we decide on the first value received:

```go
type QSpec struct {
   quorumSize int
}

// ReadQF is the quorum function for the Read RPC method.
func (qs *QSpec) ReadQF(_ *ReadResponse, replies map[uint32]*State) (*State, bool) {
   for key, value := replies {
      // Return the first value in the loop along with true.
      // No other value will be returned after
      return value, true
   }
   // If there are no value in the map, return false.
   // The ReadQF quorum function will then be executed again later
   return nil, false
}
```

A more complex quorum function, where we first wait to receive all replies, then return the reply with the most recent state can be found below:

```go
// ReadQF is the quorum function for the Read RPC method.
func (qs *QSpec) ReadQF(_ *ReadResponse, replies map[uint32]*State) (*State, bool) {
   if len(replies) < qs.quorumSize {
      return nil, false
   }
   return newestState(replies), true
}

func newestState(replies map[uint32]*State) *State {
   var newest *State
   for _, s := range replies {
      if s.GetTimestamp() >= newest.GetTimestamp() {
         newest = s
      }
   }
   return newest
}
```

The type implementing the `QuorumSpec` interface contains the variables used in the quorum functions.
It should be passed when creating the configuration.
In the code snippet below a `QSpec` type with a `quorumSize` of `2` is created and passed to the configuration.

```go
allNodesConfig, err := mgr.NewConfiguration(
   &QSpec{2},
   gorums.WithNodeList(addrs),
)
```

After the configuration is created with a `QuorumSpec` the quorum call can be invoked by calling the corresponding methods on the configuration.
The code snippet below invokes the `Read` quorum call on the configuration that was created.

```go
// Invoke read quorum call:
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
reply, err := allNodesConfig.Read(ctx, &ReadRequest{})
if err != nil {
   log.Fatal("Read call returned error:", err)
}
```

## Task

Using gorums, create a storage service that stores values on individual nodes and returns all the values stored on the nodes.
The server should provide two RPC calls:

- `Write`: A regular gRPC call that is called on a single node and stores a value in a slice on that node.

- `Read`: A quorum call that is called on all nodes in the configuration and returns a slice of all the values in the slices returned by the nodes.
  If a value is stored on multiple nodes it will appear multiple times.

You should add a service to the [`storage.proto`](proto/storage.proto) file.
You should also define the RPCs used by the service and the messages.

Some skeleton code has been provided in the [`storage_client.go`](/lab2/gorums/storage_client.go) and the [`storage_server.go`](/lab2/gorums/storage_server.go) files.
The skeleton codes includes some unimplemented functions that should be implemented, as well as some implemented helper functions used for testing.
The signature of these functions should not be changed.
You must also add functions to the `StorageServer` so that it implements the interface defined in the compiled [`storage.proto`](proto/storage.proto) files.
