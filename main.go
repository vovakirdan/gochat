package main

import (
	"log"
	"net"
)

func main() {
	server := NewServer()
	go server.Run()

	listener, err := net.Listen("tcp", ":7878")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Println("Server started on :7878")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go server.HandleConnection(conn)
	}
}