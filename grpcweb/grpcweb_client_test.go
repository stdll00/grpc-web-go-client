package grpcweb

import (
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"testing"
)

func TestWithGeneratedClient(t *testing.T) {
	client, err := DialContext(":50051")
	if err != nil {
		panic("err")
	}
	helloworld.NewGreeterClient(client) // just check interface
}
