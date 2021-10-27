package main

import (
	"fmt"
	"log"
	"net"

	"gitlab.com/smallwood/chatapp-server/msg"
	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err.Error())
	}

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)

	msg.RegisterChatServiceServer(server, &ChatServiceServer{clients: make(map[int]chan *msg.Msg)})
	fmt.Println("listen on localhost:8000")
	server.Serve(listen)
}
