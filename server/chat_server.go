package main

import (
	"context"
	"fmt"
	"math/rand"

	"gitlab.com/smallwood/chatapp-server/db"
	"gitlab.com/smallwood/chatapp-server/msg"
)

type ChatServiceServer struct {
	msg.UnimplementedChatServiceServer
	clients map[int]chan *msg.Msg
}

func (c *ChatServiceServer) Broadcast(msg *msg.Msg) {
	for _, client := range c.clients {
		client <- msg
	}
}

// TODO: save message with sqlite

func (c *ChatServiceServer) InitChat(context.Context, *msg.Empty) (*msg.Msgs, error) {
	messages, err := db.Handler().SelectAllMessages()
	if err != nil {
		fmt.Println("[ERROR] failed to select all messages:", err.Error())
		return nil, err
	}
	return messages, nil
}

func (c *ChatServiceServer) RecvMessage(user *msg.User, stream msg.ChatService_RecvMessageServer) error {
	userChannel := make(chan *msg.Msg)
	clientId := rand.Intn(100000) + len(c.clients)
	c.clients[clientId] = userChannel
	fmt.Printf("user %s has joined\n", user.GetName())
	m := &msg.Msg{
		Sender: &msg.User{Name: "server"},
		Msg:    "user " + user.GetName() + " has joined",
	}
	go c.Broadcast(m)

	for {
		select {
		case <-stream.Context().Done():
			fmt.Printf("streaming for %s is done!\n", user.GetName())
			delete(c.clients, clientId)
			m := &msg.Msg{
				Sender: &msg.User{Name: "server"},
				Msg:    user.GetName() + " has left the party",
			}
			go c.Broadcast(m)
			return nil
		case msg := <-userChannel:
			stream.Send(msg)
		}
	}
}

func (c *ChatServiceServer) SendMessage(ctx context.Context, m *msg.Msg) (*msg.MsgAck, error) {
	go c.Broadcast(m)
	go func() {
		err := db.Handler().InsertMessage(m)
		if err != nil {
			fmt.Println("[ERROR] failed to save message:", err.Error())
		}
	}()
	return &msg.MsgAck{Status: "SENT"}, nil
}
