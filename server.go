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
	s.Info(fmt.Sprintf("New connection from %s", conn.RemoteAddr().String()))
	fmt.Fprintf(conn, "Welcome to the GoChat!\n\nWho are you?\n(stranger): ")

	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username:", err)
		return
	}
	username = strings.TrimSpace(username)

	// sign up/in logic
	// if user already exists, type nice to see you and ask for a password
	if s.database.IsUserExists(username) {
		// if client is online (already in clients map) - say you can't log in
		if client, ok := s.clients[username]; ok {
			fmt.Fprintf(conn, "Hmm... seems like you are already in!\n")
			s.SystemMessage(client, fmt.Sprintf("Someone trying to connect from %s", conn.RemoteAddr().String()))
			return
		}
		fmt.Fprintf(conn, "Nice to meet you, %s!\nEnter your password\n> ", username)
		password, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading password:", err)
			return
		}
		password = strings.TrimSpace(password)
		// if password is correct, welcome back
		if s.database.ValidateUser(username, password) {
			fmt.Fprintln(conn, "Everything is correct!")
		} else {
			fmt.Fprintf(conn, "Wrong password for user %s!\n", username)
			return
		}
	} else {
		fmt.Fprintf(conn, "We don't know each other, let's register, %s!\nEnter password\n> ", username)
		password, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading password:", err)
			return
		}
		// todo check for password strength
		if !s.database.Register(username, password) {
			fmt.Fprintf(conn, "Error registering user %s\n", username)
			return
		}
	}
	isAdmin := false
	if username == "admin" {isAdmin = true}  // todo make it better
	client := &ClientContext{
		Username: username,
		Conn: conn,
		Room: "main",
		Admin: isAdmin,
		LastTimeIn: time.Now().Format("2006-01-02 15:04:05"),  // dunno why
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
		case 1:  // empty command
			s.EmptyMessage(client)
		case 2:  // wrong command
			s.Info(fmt.Sprintf("User (%s) in room [%s] typed wrong command: %s\n", username, client.Room, message))
		case 3:  // unknown switch type
			s.SystemMessage(client, "Unknown switch type")
			s.Info(fmt.Sprintf("User (%s) in room [%s] typed unknown switch type: %s\n", username, client.Room, message))
		case 4: // unknown create type
			s.SystemMessage(client, "Unknown create type")
			s.Info(fmt.Sprintf("User (%s) in room [%s] typed unknown create type: %s\n", username, client.Room, message))
		case 5: // unknown list type
			s.SystemMessage(client, "Unknown list type")
			s.Info(fmt.Sprintf("User (%s) in room [%s] typed unknown list type: %s\n", username, client.Room, message))
		case 6: // unknown count type
			s.SystemMessage(client, "Unknown count type")
			s.Info(fmt.Sprintf("User (%s) in room [%s] typed unknown count type: %s\n", username, client.Room, message))
		case 10: // client wants to logout
			s.SystemMessage(client, "Goodbye, see you later.")
			s.RemoveClient(username)
			return
		// default:
			// s.SystemMessage(client, "Invalid or unknown command. Type /help <command> to see commands.")
		}
	}
}

func (s *Server) HandleMessage(client *ClientContext, message string) (int, error) {
	message = strings.TrimSpace(message)
	// if command
	if strings.HasPrefix(message, "/") {
		status, err := s.ParseCommand(client, message) 
		return status, err
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
		isPrivate := ""
		if s.database.IsPrivateRoom(roomName) {
			isPrivate = "(private) "
		}
		s.SystemMessage(client, fmt.Sprintf("%sRoom '%s' created.\n", isPrivate, roomName))
		s.Info(fmt.Sprintf("%sRoom '%s' created.\n", isPrivate, roomName))
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
	fmt.Print("] [INFO] ")
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

func (s *Server) SendHelp(client *ClientContext, command string) {
	switch command {
	case "switch":
		s.SystemMessage(client, "Usage: /switch room <room_name> [password]")
	case "create":
		s.SystemMessage(client, "Usage: /create room <room_name> <password>")
	default:
		message := "Available commands:\n"
		message += "/create | /cr room <room_name> <password>\n"
		message += "If no password specified, will create opened room: everyone can join.\n"
		message += "/switch | /sw room <room_name> [password]\n"
		message += "/list <rooms|users>\n"
		message += "/help <command>\n"
		message += "Will show description of specified command."
		message += "/quit | /q\n"
		message += "Log out from account."
		s.SystemMessage(client, message)
	}
	s.EmptyMessage(client)
}

func (s *Server) ParseCommand(client *ClientContext, message string) (int, error) {
	// assume message starts with "/"

	// preprocess message: trim, split, etc
	message = strings.TrimSpace(message)
	// split by parts
	parts := strings.Split(message, " ")
	argscount := len(parts)
	if len(parts) == 0 {
		// empty command
		return 1, nil
	}

	command := parts[0]
	switch command {
	case "/switch":
		// move to the new room example
		if argscount < 2 {
			s.SystemMessage(client, "You need to specify what do you want to switch.")
			return 2, nil
		}
		switch parts[1] {
		case "room":
			if argscount < 3 {
				s.ListSomethingToClient(client, "rooms")
				return 0, nil
			} 
			if argscount == 3 {
				s.ChangeRoom(client.Username, parts[2], "")  // non-private room
				return 0, nil
			} 
			if argscount == 4 {
				s.ChangeRoom(client.Username, parts[2], parts[3])  // private room
				return 0, nil
			}
		default:
			// unknown switch type
			return 3, nil
		}
	case "/create":
		// create room example
		if argscount < 2 {
			s.SystemMessage(client, "You need to specify what do you want to create.")
			return 2, nil
		}
		switch parts[1] {
		case "room":
			if argscount < 3 {
				s.SystemMessage(client, "You need to specify room name and (if neccessary) password")
				return 0, nil
			} 
			if argscount == 3 {
				s.CreateRoom(client, parts[2], "")  // non-private room
				return 0, nil
			} 
			if argscount == 4 {
				s.CreateRoom(client, parts[2], parts[3])  // private room
				return 0, nil
			}
		default:
			// unknown create type
			return 4, nil
		}
	case "/list":
		if argscount < 2 {
			s.SystemMessage(client, "You need to specify what do you want to list.")
			return 2, nil
		}
		switch parts[1] {
		case "rooms":
			s.ListSomethingToClient(client, "rooms")
			return 0, nil
		case "users":
			s.ListSomethingToClient(client, "users")
			return 0, nil
		default:
			// unknown list type
			return 5, nil
		}
	case "/count":
		if argscount < 2 {
			s.SystemMessage(client, "You need to specify what do you want to count.")
			return 2, nil
		}
		switch parts[1] {
		case "rooms":
			roomsCount := len(s.database.rooms)
			s.SystemMessage(client, fmt.Sprintf("There are %d rooms.", roomsCount))
			return 0, nil
		case "users":
			usersCount := len(s.database.users)
			s.SystemMessage(client, fmt.Sprintf("There are %d users.", usersCount))
			return 0, nil
		default:
			// unknown count type
			return 6, nil
		}
	case "/quit", "/q":
		return 10, nil 
	case "/help":
		if argscount < 2 {
			s.SendHelp(client, "")
			return 0, nil
		}
		s.SendHelp(client, parts[1])
		return 0, nil
	default:
		// unknown command
		return 2, nil
	}
	return 0, nil
}
