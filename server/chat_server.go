package main

import (
	"context"
	"fmt"
	"math/rand"

	"gitlab.com/smallwood/chatapp-server/msg"
)

type ChatServiceServer struct {
	msg.UnimplementedChatServiceServer
	clients map[int]chan *msg.Msg
}

// TODO: grab message from sqlite

func (c *ChatServiceServer) RecvChat(user *msg.User, stream msg.ChatService_RecvChatServer) error {
	userChannel := make(chan *msg.Msg)
	clientId := rand.Intn(100000) + len(c.clients)
	c.clients[clientId] = userChannel
	fmt.Printf("user %s has joined\n", user.GetName())

	for {
		select {
		case <-stream.Context().Done():
			fmt.Printf("streaming for %s is done!\n", user.GetName())
			delete(c.clients, clientId)
			return nil
		case msg := <-userChannel:
			fmt.Printf("sending '%s' to %s\n", msg.GetMsg(), user.GetName())
			stream.Send(msg)
		}
	}
}

func (c *ChatServiceServer) SendMessage(ctx context.Context, m *msg.Msg) (*msg.MsgAck, error) {
	go func() {
		for _, client := range c.clients {
			client <- m
		}
	}()

	return &msg.MsgAck{Status: "SENT"}, nil
}
