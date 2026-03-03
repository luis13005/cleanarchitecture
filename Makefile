gen:
	protoc --proto_path=internal\infra\grpc\protofiles internal/infra/grpc/protofiles/*.proto --go_out=internal/infra/grpc/pb --go-grpc_out=internal/infra/grpc/pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

evans:
	docker run --rm -it -v "C:/Users/LUISFP/go/src/luisfp/pos/cleanarchitecture:/mount:ro" ghcr.io/ktr0731/evans:latest --path ./internal/infra/grpc/protofiles/ --proto order.proto --host host.docker.internal --port 50051 repl

graph:
	go run github.com/99designs/gqlgen generate