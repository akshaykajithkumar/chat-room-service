package main

import (
	"chat-app/server"
	"log"
	"net/http"
)

func main() {
	server := server.NewServer()
	log.Println("Chat server started on :5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
