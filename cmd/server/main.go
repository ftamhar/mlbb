package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"example_nrpc/proto/hello"

	"github.com/darmawan01/toldata"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type server struct{}

func (s *server) Upload2(ctx context.Context, req *hello.UploadRequest) (*hello.UploadResponse, error) {
	name := uuid.NewString() + ".jpeg"
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	b := req.GetData()
	_, err = f.Write(b)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	res := hello.UploadResponse{
		Name: name,
	}

	return &res, nil
}

func (s *server) Greeting(ctx context.Context, req *hello.GreetingRequest) (resp *hello.GreetingResponse, err error) {
	resp = &hello.GreetingResponse{
		Fullname: req.Firstname + " " + req.Lastname,
	}
	return
}

func (s *server) Upload(stream hello.HelloServices_UploadToldataServer) {
	var err error
	var ur *hello.UploadRequest

	name := uuid.NewString() + ".jpeg"
	f, err := os.Create(name)
	if err != nil {
		stream.Error(err)
	}
	defer f.Close()

	for {
		ur, err = stream.Receive()
		if err == io.EOF {
			break
		}

		if err != nil {
			stream.Error(err)
		}
		f.Write(ur.Data)
	}

	err = stream.Done(&hello.UploadResponse{
		Name: name,
	})
	if err != nil {
		stream.Error(err)
	}
}

func (s *server) ToldataHealthCheck(ctx context.Context, req *toldata.Empty) (*toldata.ToldataHealthCheckInfo, error) {
	return nil, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Server run context
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	bus, err := toldata.NewBus(ctx, toldata.ServiceConfiguration{URL: nats.DefaultURL})
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	// The NATS handler from the helloworld.nrpc.proto file.
	svc := hello.NewHelloServicesToldataServer(bus, &server{})
	c2, err := svc.SubscribeHelloServices()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server is running, ^C quits.")
	<-sig
	cancel()
	<-c2
}
