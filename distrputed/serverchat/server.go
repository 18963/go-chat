package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Message struct {
	Username string
	Text     string
}

type ChatService struct {
	mu       sync.Mutex
	messages []Message
}

type ChatInput struct {
	Username string
	Text     string
}

func (cs *ChatService) SendMessage(input ChatInput, reply *[]Message) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Add new message
	cs.messages = append(cs.messages, Message(input))

	// Print full history to server terminal
	fmt.Println("\n--- Chat History ---")
	for _, msg := range cs.messages {
		fmt.Printf("[%s]: %s\n", msg.Username, msg.Text)
	}
	fmt.Println("--------------------")

	// Return history to client
	*reply = append([]Message{}, cs.messages...)
	return nil
}

func (cs *ChatService) GetHistory(dummy string, reply *[]Message) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	*reply = append([]Message{}, cs.messages...)
	return nil
}

func main() {
	cs := new(ChatService)
	rpc.Register(cs)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Server error:", err)
		return
	}
	fmt.Println("Chat server started on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
