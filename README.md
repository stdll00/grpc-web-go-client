# gRPC Web Go client
[![GoDoc](https://godoc.org/github.com/stdll00/grpc-web-go-client/grpcweb?status.svg)](https://godoc.org/github.com/stdll00/grpc-web-go-client/grpcweb)
[![GitHub Actions](https://github.com/stdll00/grpc-web-go-client/workflows/main/badge.svg)](https://github.com/stdll00/grpc-web-go-client/actions)  

Easy to use, compatibles with generated client by `protoc-gen-go-grpc`.  
Stream, grpc-web-text is not supported now.  
This is fork of https://github.com/ktr0731/grpc-web-go-client 

Usage: (examples/main.go)
```go
dial, _ := grpcweb.Dial(addr, grpcweb.WithInsecure())
client := helloworld.NewGreeterClient(dial) // compatible with normal grpc-go code!
res, _ := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "hello"})
```

### How to run
```
export ENVOY_ADMIN_PORT=8001 
export GRPC_PORT=8002 
export GRPC_WEB_PORT=8003
```


```
ENVOY_CONFIG=$(envsubst < examples/envoy.envsubst.yaml) &&  envoy --config-yaml $ENVOY_CONFIG
go run examples/server/main.go
go run examples/main.go
```
