package main

import (
	"context"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"sync"
	"time"
)

const (
	address = "localhost:50055"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := desc.NewChatV1Client(conn)

	chatID, err := createChat(ctx, client)
	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, "vova", 5*time.Second)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, "kate", 7*time.Second)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	wg.Wait()
}

func connectChat(ctx context.Context, client desc.ChatV1Client, chatID int64, username string, period time.Duration) error {
	stream, err := client.Connect(ctx, &desc.ConnectRequest{
		Id:       chatID,
		Username: username,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}

			log.Printf("[%v] - [from: %s]: %s\n",
				message.GetTimestamp().AsTime().Format(time.RFC3339),
				message.GetFrom(),
				message.GetText(),
			)
		}
	}()

	for {
		// Ниже пример того, как можно считывать сообщения из консоли
		// в демонстрационных целях будем засылать в чат рандомный текст раз в 5 секунд
		//scanner := bufio.NewScanner(os.Stdin)
		//var lines strings.Builder
		//
		//for {
		//	scanner.Scan()
		//	line := scanner.Text()
		//	if len(line) == 0 {
		//		break
		//	}
		//
		//	lines.WriteString(line)
		//	lines.WriteString("\n")
		//}
		//
		//err = scanner.Err()
		//if err != nil {
		//	log.Println("failed to scan message: ", err)
		//}

		time.Sleep(period)

		text := gofakeit.Word()

		_, err = client.SendMessage(ctx, &desc.SendMessageRequest{
			Message: &desc.Message{
				From:      username,
				Text:      text,
				ChatId:    chatID,
				Timestamp: timestamppb.Now(),
			},
		})
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}

func createChat(ctx context.Context, client desc.ChatV1Client) (int64, error) {
	res, err := client.Create(ctx, &desc.CreateRequest{Usernames: []string{"vova", "kate"}})
	if err != nil {
		return 0, err
	}

	return res.GetId(), nil
}
