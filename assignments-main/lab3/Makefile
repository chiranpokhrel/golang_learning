binaries := bin
replica := $(binaries)/replica
ifeq ($(OS),Windows_NT)
    replica = $(binaries)/replica.exe
endif
gorum_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
proto_src := gorumsfd/proto/fd.proto
proto_go := $(proto_src:%.proto=%.pb.go)

all: pre protos build

.PHONY: pre
pre:
	@mkdir -p $(binaries)
	@chmod +x ./replicas.sh

.PHONY: protos
protos: $(proto_go)
	@echo "+ compiling gorums proto files"

%.pb.go %_gorums.pb.go: %.proto
	@protoc -I=$(gorum_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<

build:
	@echo "+ compiling replica binary"
	@go build $(BUILD_FLAGS) -o $(replica) ./cmd/replica/

.PHONY: clean
clean:
	rm -rf $(binaries)
