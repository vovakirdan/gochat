package main

import (
	"net"
)

type ClientContext struct {
	Username string
	Room string
	Conn net.Conn
}