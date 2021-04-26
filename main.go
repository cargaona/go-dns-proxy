package main

import (
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

	// UnresolvedDNS holds the actual unresolved dns message.
	var unresolvedDNS [2222]byte
	responseSize, err := clientConnection.Read(unresolvedDNS[:])
	if err != nil {
		log.Println(err)
	}

	// Write to the resolver
	responseSize, err = resolverConnection.Write(unresolvedDNS[:responseSize])
	if err != nil {
		log.Println(err)
	}

	// Get response from resolver
	var resolvedDNS [2222]byte
	responseSize, err = resolverConnection.Read(resolvedDNS[:])
	if err != nil {
		log.Println(err)
	}

	// Send response to client
	responseSize, err = clientConnection.Write(resolvedDNS[:])
	if err != nil {
		log.Println(err)
	}

}

func main() {
	server()
}
