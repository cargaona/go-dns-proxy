package dnsproxy

import (
	"fmt"
	"log"
	"net"
)

type service struct {
	PrimaryResolver   string
	SecondaryResolver string
	LocalPort         int
	Protocol          string
	Cache             Cache
}

type Service interface {
	Serve() error
}

type Cache interface{}

func NewProxy(primaryResolver, secondaryResolver, protocol string, localPort int, cache Cache) Service {
	return &service{
		PrimaryResolver:   primaryResolver,
		SecondaryResolver: secondaryResolver,
		LocalPort:         localPort,
		Protocol:          protocol,
		Cache:             cache,
	}
}

func (svc *service) Serve() error {
	return svc.startServer()
}

func (svc *service) startServer() error {
	listener, err := net.Listen(svc.Protocol, fmt.Sprintf(":%d", svc.LocalPort))
	//	defer listener.Close()
	if err != nil {
		return err
	}

	for {
		log.Println("Accepting connection")
		clientConnection, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		log.Println("Spawning new goroutine to handle connection")
		go svc.handleConnection(clientConnection)
	}
}

func (svc *service) handleConnection(clientConnection net.Conn) {
	//create connection to external dns server
	log.Println("Creating connection to DNS Resolver", svc.PrimaryResolver)
	resolverConnection, err := net.Dial(svc.Protocol, svc.PrimaryResolver)
	if err != nil {
		log.Fatal(err)
	}
	defer resolverConnection.Close()
	log.Println("Passing request to primaryResolver")

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
