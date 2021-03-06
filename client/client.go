package main

import (
	"context"
	"fmt"

	"gitlab.com/smallwood/sw-chat/channel"
	"gitlab.com/smallwood/sw-chat/msg"
	"google.golang.org/grpc"
)

func startClient() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:8000", opts...)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client := msg.NewChatServiceClient(conn)

	initClient(ctx, client)
	go recv(ctx, client)
	go send(ctx, client)

	return conn, nil
}

func initClient(ctx context.Context, client msg.ChatServiceClient) {
	empty := &msg.Empty{}
	messages, err := client.InitChat(ctx, empty)
	if err != nil {
		errStr := fmt.Sprintf("client %s failed to init [server_error=%s]", username, err.Error())
		channel.ScreenChan <- channel.Data{Type: "msg", Object: errStr}
		return
	}

	for _, message := range messages.Msgs {
		channel.ScreenChan <- channel.Data{Type: "msg", Object: fmt.Sprintf("[%s] %s", message.GetSender(), message.GetMsg())}
	}
}

func recv(ctx context.Context, client msg.ChatServiceClient) {
	user := &msg.User{Name: username}
	stream, err := client.RecvMessage(ctx, user)
	if err != nil {
		errStr := fmt.Sprintf("client %s failed to joined [server_error=%s]", user.Name, err.Error())
		channel.ScreenChan <- channel.Data{Type: "msg", Object: errStr}
		return
	}

	for {
		select {
		case <-channel.GlobalShutdown:
			return
		default:
			in, err := stream.Recv()
			if err != nil {
				channel.ScreenChan <- channel.Data{Type: "msg", Object: "[ERROR] in recv: " + err.Error()}
				return
			}

			channel.ScreenChan <- channel.Data{Type: "msg", Object: fmt.Sprintf("[%s] %s", in.GetSender(), in.GetMsg())}
		}
	}
}

func send(ctx context.Context, client msg.ChatServiceClient) {
	for {
		select {
		case <-channel.GlobalShutdown:
			return
		case m := <-channel.MsgChan:
			ack, err := client.SendMessage(ctx, &msg.Msg{
				Sender: &msg.User{
					Name: username,
				},
				Msg: m,
			})
			if err != nil {
				channel.ScreenChan <- channel.Data{Type: "msg", Object: "[ERROR] in send: " + err.Error()}
				continue
			}

			channel.HeaderChan <- channel.Data{Type: "msg", Object: ack.Status}
		}
	}
}
