package server

import (
	"chat-app/chatroom"
	"chat-app/client"
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	chatRoom *chatroom.ChatRoom
}

func NewServer() *Server {
	return &Server{
		chatRoom: chatroom.NewChatRoom(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/join":
		s.handleJoin(w, r)
	case "/send":
		s.handleSend(w, r)
	case "/leave":
		s.handleLeave(w, r)
	case "/messages":
		s.handleMessages(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handleJoin(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	_, exists := s.chatRoom.GetClientByID(clientID)
	if exists {
		http.Error(w, "Client already in the chat", http.StatusConflict)
		return
	}

	client := client.NewClient(clientID)
	s.chatRoom.AddClient(client)
	s.chatRoom.SendMessageToRoom(fmt.Sprintf("%s joined the chat", clientID))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s joined the chat\n", clientID)
}

func (s *Server) handleSend(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("id")
	message := r.URL.Query().Get("message")

	if clientID == "" || message == "" {
		http.Error(w, "Client ID and message required", http.StatusBadRequest)
		return
	}

	_, exists := s.chatRoom.GetClientByID(clientID)
	if !exists {
		http.Error(w, "Unknown client, please join the chat first", http.StatusUnauthorized)
		return
	}

	s.chatRoom.SendMessageToRoom(fmt.Sprintf("%s: %s", clientID, message))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message sent\n")
}

func (s *Server) handleLeave(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	_, exists := s.chatRoom.GetClientByID(clientID)
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	s.chatRoom.RemoveClient(clientID)
	s.chatRoom.SendMessageToRoom(fmt.Sprintf("%s left the chat", clientID))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s left the chat\n", clientID)
}

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	client, exists := s.chatRoom.GetClientByID(clientID)
	if !exists {
		http.Error(w, "Unknown client, cannot receive messages", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case msg := <-client.MessageCh:
			_, exists := s.chatRoom.GetClientByID(clientID)
			if !exists {
				http.Error(w, "Client has left the chat", http.StatusGone)
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()

		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				http.Error(w, "Connection timed out", http.StatusGatewayTimeout)
			}
			return
		}
	}
}
