package chatroom

import (
	"chat-app/client"
	"fmt"
	"sync"
)

type ChatRoom struct {
	connectedClients map[string]*client.Client
	messageChannel   chan string
	mu               sync.Mutex
}

func NewChatRoom() *ChatRoom {
	room := &ChatRoom{
		connectedClients: make(map[string]*client.Client),
		messageChannel:   make(chan string),
	}

	go room.handleMessageBroadcasting()
	return room
}

func (room *ChatRoom) AddClient(client *client.Client) {
	room.mu.Lock()
	defer room.mu.Unlock()
	room.connectedClients[client.ID] = client
}

func (room *ChatRoom) RemoveClient(clientID string) {
	room.mu.Lock()
	defer room.mu.Unlock()
	if client, exists := room.connectedClients[clientID]; exists {
		close(client.MessageCh)
		delete(room.connectedClients, clientID)
	}
}

func (room *ChatRoom) GetClientByID(clientID string) (*client.Client, bool) {
	room.mu.Lock()
	defer room.mu.Unlock()
	client, exists := room.connectedClients[clientID]
	return client, exists
}

func (room *ChatRoom) SendMessageToRoom(message string) {
	room.messageChannel <- message
}

func (room *ChatRoom) handleMessageBroadcasting() {
	for message := range room.messageChannel {
		room.mu.Lock()
		for _, cl := range room.connectedClients {
			go func(clientInstance *client.Client) {
				select {
				case clientInstance.MessageCh <- message:
				default:
					fmt.Printf("Client %s is not receiving messages. Channel full or disconnected.\n", clientInstance.ID)
				}
			}(cl)
		}
		room.mu.Unlock()
	}
}
