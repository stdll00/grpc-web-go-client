package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	helloworldpb "google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
	"net"
	"os"
)

type helloServer struct {
	helloworldpb.UnimplementedGreeterServer
}

func (*helloServer) SayHello(_ context.Context, req *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	return &helloworldpb.HelloReply{Message: "" + req.Name}, nil
}

func main() {
	addr := fmt.Sprintf("127.0.0.1:%v", os.Getenv("GRPC_PORT"))
	log.Printf("listening %s\n", addr)
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	helloworldpb.RegisterGreeterServer(s, &helloServer{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
