# Installing the Protobuf Compiler and Plugins

The protobuf compiler is called `protoc`, and you will need to install it for this assignment.
Most Linux distributions provides a `protobuf` package.

On Ubuntu/Debian:

```shell
sudo apt install -y protobuf-compiler
```

On Archlinux:

```shell
sudo pacman -S protobuf
```

On macOS, if you have installed homebrew, you can simply run:

```shell
brew install protobuf
```

If you do not use a package manager with your OS, you should download the appropriate package from the [official release page of the Protobuf compiler](https://github.com/protocolbuffers/protobuf/releases).

Once you have installed the make sure to test that the installation is working by running:

```shell
$ protoc --version
libprotoc 25.1
```

Next, you need to install the plugins that are needed to generate Go protobuf code and gRPC code.
This can be done using these `go install` commands:

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

This will install the `protoc-gen-go` and `protoc-gen-go-grpc` commands in your `$GOPATH/bin` folder.
To test whether or not you can use these plugins, run:

```shell
$ protoc-gen-go --version
protoc-gen-go v1.32.0
$ protoc-gen-go-grpc --version
protoc-gen-go-grpc 1.3.0
```

## Compiling .proto Files

The proto file [`kv.proto`](./proto/kv.proto) needs to be compiled using `protoc` in order to generate the `kv.pb.go` and `kv_grpc.pb.go` files which are used by the Go client/server implementations in this assignment.
To compile the proto file, run the command below from the `lab2/grpc` directory:

```shell
cd lab2/grpc
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. proto/kv.proto
```

Or if you are in the `proto` directory:

```shell
cd lab2/grpc/proto
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. kv.proto
```

## Troubleshooting

If the plugins are not found, then you may need to add the following line to your shell's configuration file to make the plugins binaries available in your default environment:

```shell
export PATH="$PATH:$(go env GOPATH)/bin"
```

- If you are using `zsh`, add the line above to your `$HOME/.zshrc` file.
- If you are using `bash`, add the line above to your `$HOME/.bashrc` file.
