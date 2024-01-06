package grpcweb

import (
	"context"
	"fmt"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"io"
	"net/http"
	"testing"
)

// httpServer
type httpServer struct {
	lastRequest    []byte
	lastHeader     http.Header
	response       []byte
	responseHeader http.Header
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Store the body of the request
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}
	s.lastRequest = body
	// Store the header of the request
	s.lastHeader = r.Header

	// Set the response header
	for key, value := range s.responseHeader {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	w.Write(s.response)
}

func TestWithGeneratedClient(t *testing.T) {
	addr := ":50051"
	ctx := context.Background()
	s := &httpServer{}
	go http.ListenAndServe(addr, s)

	dial, err := Dial(addr)
	if err != nil {
		panic("err")
	}
	client := helloworld.NewGreeterClient(dial)
	_, err = client.SayHello(ctx, &helloworld.HelloRequest{Name: "hello"})
	fmt.Print(s.lastRequest)
	if err != nil {
		return
	}
}
