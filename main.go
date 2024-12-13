package main

import (
	"log"
	"net"
	"flag"
	"fmt"
)

func main() {
	defaultIP := "127.0.0.1"
	defaultPort := "7878"

	// Определение флагов
	ip := flag.String("ip", defaultIP, "IP address to bind the server to")
	port := flag.String("port", defaultPort, "Port to bind the server to")

	// Разбор флагов
	flag.Parse()

	// Формируем адрес из IP и порта
	address := fmt.Sprintf("%s:%s", *ip, *port)
	server := NewServer()
	go server.Run()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server started on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go server.HandleConnection(conn)
	}
}