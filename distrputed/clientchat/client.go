package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
)

type Message struct {
	Username string
	Text     string
}

type ChatInput struct {
	Username string
	Text     string
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Println("Connected. Type your message (type 'exit' to quit):")

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		var history []Message
		input := ChatInput{Username: username, Text: text}
		err := client.Call("ChatService.SendMessage", input, &history)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("\n--- Chat History ---")
		for _, msg := range history {
			fmt.Printf("[%s]: %s\n", msg.Username, msg.Text)
		}
		fmt.Println("--------------------")
	}
}
