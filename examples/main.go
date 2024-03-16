package main

import (
	"context"
	"fmt"
	"github.com/stdll00/grpc-web-go-client/grpcweb"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
	"os"
)

func main() {
	addr := "127.0.0.1:" + os.Getenv("GRPC_WEB_PORT")
	dial, _ := grpcweb.Dial(addr, grpcweb.WithInsecure())
	client := helloworld.NewGreeterClient(dial)
	res, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "hello,this is sample"})
	if err != nil {
		log.Fatalf("failed to send request: %v", err.Error())
	}
	fmt.Print(res.Message)
}
