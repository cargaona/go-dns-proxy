package main

import (
	"io"
	"log"
	"net"
)

func server() {
	// start server
	log.Println("Starting TCP Server.")
	server, err := net.Listen("tcp", "localhost:8853")
	if err != nil {
		log.Println("error creating server")
	}

	for {
		log.Println("Accepting connection")
		conn, err := server.Accept()
		if err != nil {
			log.Println(err)
		}
		log.Println("Spawning new goroutine to handle connection")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	//create connection to external dns server
	log.Println("Creating connection to DNS Resolver")
	client, err := net.Dial("tcp", "1.1.1.1:53")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	log.Println("Copying messages?")
	// Sending message to cloudflare
	io.Copy(client, conn)
}

func main() {
	server()
}
