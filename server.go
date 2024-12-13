package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Server struct {
	clients map[string]*ClientContext
	clientsMu sync.Mutex
	database *Database
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]*ClientContext),
		database: NewDatabase(),
	}
}

func (s *Server) Run() {
	// here will be background logic 
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	fmt.Fprintf(conn, "Welcome to the GoChat!\n\nWho are you?\n(stranger): ")

	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username:", err)
		return
	}
	username = strings.TrimSpace(username)

	client := &ClientContext{
		Username: username,
		Conn: conn,
		Room: "main",
	}

	s.clientsMu.Lock()
	s.clients[username] = client
	s.clientsMu.Unlock()

	fmt.Fprintf(conn, "Welcome, %s! You are now in the 'main' room.\n(%s): ", username, username)

	for {
		message, err := reader.ReadString('\n')
		// if error
		if err != nil {
			fmt.Printf("User %s disconnected\n", username)
			s.RemoveClient(username)
			return
		}
		
		// if command
		message = strings.TrimSpace(message)
		if strings.HasPrefix(message, "/room") {
			// move to the new room example
			// todo refactor
			newRoom := strings.TrimPrefix(message, "/room ")
			s.ChangeRoom(username, newRoom)
			continue
		}

		// if empty message
		if message == "" {
			s.EmptyMessage(*client)
			continue
		}

		// or else broadcast to room
		s.BroadcastMessage(*client, message)
	}
}

func (s *Server) BroadcastMessage(sender ClientContext, message string) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for _, client := range s.clients {
		if client.Username != sender.Username && client.Room == sender.Room {
			// fmt.Fprintf(client.Conn, "[%s]: %s\n", sender, message)
			s.EmptyMessage(sender)
			s.SendMessage(*client, sender, message)
		}
	}
}

func (s *Server) ChangeRoom(username, newRoom string) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	client, exists := s.clients[username]
	if !exists {  // strange place
		return
	}

	client.Room = newRoom
	fmt.Fprintf(client.Conn, "You have been moved to the '%s' room.\n", newRoom)
	s.EmptyMessage(*client)
}

func (s *Server) RemoveClient(username string) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	delete(s.clients, username)
}

// ------------
// utils

func (s *Server) SystemMessage(conn net.Conn, message string) {
	fmt.Fprintf(conn, "(system): %s\n", message)
}

func (s *Server) SendMessage(client, sender ClientContext, message string) {
	fmt.Fprintf(client.Conn, "\n[%s] %s: %s\n(%s): ", sender.Room, sender.Username, message, client.Username)
}

func (s *Server) EmptyMessage(client ClientContext) {
	fmt.Fprintf(client.Conn, "(%s): ", client.Username)
}
