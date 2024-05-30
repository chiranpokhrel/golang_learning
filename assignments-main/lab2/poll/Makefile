binaries := bin
participant_bin		:= $(binaries)/participant
pollmaster_bin		:= $(binaries)/pollmaster
gorums_include		:= $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
master_proto		:= proto/pollmaster.proto
participant_proto	:= proto/pollparticipant.proto

all: pre pollmaster participant

.PHONY: pre
pre:
	@mkdir -p $(binaries)

pollmaster: proto
	@echo "+ compiling pollmaster"
	@go build $(BUILD_FLAGS) -o $(pollmaster_bin) cmd/pollmaster/main.go

participant: proto
	@echo "+ compiling participant"
	@go build $(BUILD_FLAGS) -o $(participant_bin) cmd/participant/main.go

proto: grpcproto gorumsproto

grpcproto: $(master_proto)
	@echo "+ compiling grpc proto files"
	@protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. $^

gorumsproto: $(participant_proto)
	@echo "+ compiling gorums proto files"
	@protoc -I=$(gorums_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. $^

solution: pre grpcproto gorumsproto
	@echo "+ compiling pollmaster"
	@go build -tags solution -o $(pollmaster_bin) cmd/pollmaster/main.go
	@echo "+ compiling participant"
	@go build -tags solution -o $(participant_bin) cmd/participant/main.go

clean:
	rm -vf $(participant_bin) $(pollmaster_bin)
