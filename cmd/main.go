package main

import (
	"log"

	mygrpc "github.com/shch989/my-grpc-go-server/internal/adapter/grpc"
	app "github.com/shch989/my-grpc-go-server/internal/application"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(log.Writer())

	hs := &app.HelloService{}
	bs := &app.BankService{}

	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)

	grpcAdapter.Run()
}
