package main

import (
	"io"
	"log"
	"net"
)

func server() {
	// start listener
	log.Println("Starting TCP Server")
	listener, err := net.Listen("tcp", "localhost:8853")
	if err != nil {
		log.Println("error creating listener")
	}

	for {
		log.Println("Accepting connection")
		clientConnection, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		log.Println("Spawning new goroutine to handle connection")
		go handleConnection(clientConnection)
	}
}

func handleConnection(clientConnection net.Conn) {
	//create connection to external dns server
	log.Println("Creating connection to DNS Resolver")
	resolverConnection, err := net.Dial("tcp", "1.1.1.1:53")
	if err != nil {
		log.Fatal(err)
	}
	defer resolverConnection.Close()
	log.Println("Passing request to cloudflare")
	go io.Copy(resolverConnection, clientConnection)
	log.Println("Writing response to our client")
	_, err = io.Copy(clientConnection, resolverConnection)
	if err != nil {
		log.Println(err)
	}
	log.Println("Response sent successfully")
}

func main() {
	server()
}
