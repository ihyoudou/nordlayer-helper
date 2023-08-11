generate_grpc_code:
	protoc --go_opt=paths=source_relative --go-grpc_out=protobuf/ --go-grpc_opt=paths=source_relative --go_out=protobuf/ protobuf/common.proto -I protobuf/
	protoc --go_opt=paths=source_relative --go-grpc_out=protobuf/ --go-grpc_opt=paths=source_relative --go_out=protobuf/ protobuf/status.proto -I protobuf/
	protoc --go_opt=paths=source_relative --go-grpc_out=protobuf/ --go-grpc_opt=paths=source_relative --go_out=protobuf/ protobuf/gateways.proto -I protobuf/

build:
	make generate_grpc_code
	env GO111MODULE=on go build -ldflags "-s -w"

pack:
	upx -9 -k nordlayer-helper
