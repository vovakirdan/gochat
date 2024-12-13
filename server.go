package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
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
		
		status, err := s.HandleMessage(client, message)
		if err != nil {
			fmt.Println("Error handling message:", err)
			return
		}
		
		switch status {
		case 1:
			s.Info(fmt.Sprintf("User (%s) in room [%s] typed wrong command: %s\n", username, client.Room, message))
		case 2:
			fmt.Printf("User %s in room %s typed private message: %s\n", username, client.Room, message)
		}
	}
}

func (s *Server) HandleMessage(client *ClientContext, message string) (int, error) {
	message = strings.TrimSpace(message)
	// if command
	if strings.HasPrefix(message, "/") {
		if s.ParseCommand(client, message) {
			// if everything is ok
			return 0, nil
		} else {
			s.SystemMessage(client, "Invalid or unknown command. Type /help <command> to see commands.")
			return 1, nil
		}
	} 
	// if private message
	if strings.HasPrefix(message, "@") {
		// do private message
		s.EmptyMessage(client)
		s.PrivateMessage(client, message)
		return 0, nil
	}
	// if empty message
	if message == "" {
		s.EmptyMessage(client)
		return 0, nil
	}
	// or else broadcast to room
	s.EmptyMessage(client)
	s.BroadcastMessage(client, message)
	return 0, nil
}

func (s *Server) BroadcastMessage(sender *ClientContext, message string) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for _, client := range s.clients {
		if client.Username != sender.Username && client.Room == sender.Room {
			s.SendMessage(client, sender, message)
		}
	}
}

func (s *Server) PrivateMessage(sender *ClientContext, message string) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	parts := strings.SplitN(message, " ", 2)
	if len(parts) < 2 {
		s.SystemMessage(sender, "Invalid private message format. Use @<username> <message>.")
		return
	}
	if !strings.HasPrefix(parts[0], "@") {
		s.SystemMessage(sender, "Invalid private message format. Use @<username> <message>.")
		return
	}
	recieverUsername := parts[0][1:]

	reciever, exists := s.clients[recieverUsername]
	if !exists {
		s.SystemMessage(sender, fmt.Sprintf("User %s doesn't exist.", recieverUsername))
		return
	}

	s.SendMessage(reciever, sender, fmt.Sprintf("(private) %s", parts[1]))
}

func (s *Server) CreateRoom(client *ClientContext, roomName, password string) bool {
	if s.database.AddRoom(roomName, password) {
		s.SystemMessage(client, fmt.Sprintf("Room '%s' created.\n", roomName))
		return true
	}
	s.SystemMessage(client, "Room with this name already exists.")
	return false
}

func (s *Server) ChangeRoom(username, newRoom, givenPassword string) bool {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	client, exists := s.clients[username]
	if !exists {  // strange place
		return false
	}
	room_password, exists := s.database.rooms[newRoom]
	if !exists {
		s.SystemMessage(client, "Room with this name doesn't exist.")
		return false
	}
	if room_password != "" && room_password != givenPassword {
		s.SystemMessage(client, "Wrong password.")
		return false
	}

	client.Room = newRoom
	s.SystemMessage(client, fmt.Sprintf("You have been moved to the '%s' room.\n", newRoom))
	return true
}

func (s *Server) RemoveClient(username string) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	delete(s.clients, username)
}

// ------------
// utils

func (s *Server) SystemMessage(client *ClientContext, message string) {
	fmt.Fprintf(client.Conn, "\n(system): %s\n", message)
	s.EmptyMessage(client)
}

func (s *Server) SendMessage(client, sender *ClientContext, message string) {
	fmt.Fprintf(client.Conn, "\n[%s] %s: %s\n(%s): ", sender.Room, sender.Username, message, client.Username)
}

func (s *Server) EmptyMessage(client *ClientContext) {
	fmt.Fprintf(client.Conn, "(%s): ", client.Username)
}

func (s *Server) Info(message string) {
	fmt.Print("[")
	fmt.Print(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Print("] [INFO]")
	fmt.Println(message)
}

func (s *Server) ListSomethingToClient(client *ClientContext, what string) {
	switch what {
	case "rooms":
		// iterate over database.rooms and create a message for client
		s.SystemMessage(client, "Available rooms:")
		i := 0
		for roomName := range s.database.rooms {
			if s.database.IsPrivateRoom(roomName) {
				fmt.Fprintf(client.Conn, "%d) %s private\n", i, roomName)
			} else {
				fmt.Fprintf(client.Conn, "%d) %s\n", i, roomName)
			}
			i++
		}
	case "users":
		// iterate over database.users and create a message for client
		s.SystemMessage(client, "Available users:")
		i := 0
		for username := range s.database.users {
			fmt.Fprintf(client.Conn, "%d) %s\n", i, username)
			i++
		}
	default:
		s.SystemMessage(client, "Unknown list type")
	}
	s.EmptyMessage(client)
}

func (s *Server) ParseCommand(client *ClientContext, message string) bool {
	// assume message starts with "/"
	// unlock client
	// s.clientsMu.Lock()
	// defer s.clientsMu.Unlock()

	// preprocess message: trim, split, etc
	message = strings.TrimSpace(message)
	// split by parts
	parts := strings.Split(message, " ")
	argscount := len(parts)
	if len(parts) == 0 {
		// empty command
		return false
	}

	command := parts[0]
	switch command {
	case "/switch":
		// move to the new room example
		if argscount < 2 {
			s.SystemMessage(client, "You need to specify what do you want to switch.")
			return true
		}
		switch parts[1] {
		case "room":
			if argscount < 3 {
				s.ListSomethingToClient(client, "rooms")
				return true
			} 
			if argscount == 3 {
				s.ChangeRoom(client.Username, parts[2], "")  // non-private room
				return true
			} 
			if argscount == 4 {
				s.ChangeRoom(client.Username, parts[2], parts[3])  // private room
			}
		default:
			// unknown switch type
			return false
		}
		return true
	case "/create":
		// create room example
		if argscount < 2 {
			s.SystemMessage(client, "You need to specify what do you want to create.")
			return true
		}
		switch parts[1] {
		case "room":
			if argscount < 3 {
				s.SystemMessage(client, "You need to specify room name and (if neccessary) password")
				return true
			} 
			if argscount == 3 {
				s.CreateRoom(client, parts[2], "")  // non-private room
				return true
			} 
			if argscount == 4 {
				s.CreateRoom(client, parts[2], parts[3])  // private room
				return true
			}
		default:
			// unknown create type
			return false
		}
		return true
	default:
		// unknown command
		return false
	}
}
