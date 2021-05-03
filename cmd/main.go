package main

import (
	"log"

	"github.com/cargaona/go-dns-proxy/pkg/domain/dnsproxy"
)

func main() {
	//cache := cache.NewCache()
	proxy := dnsproxy.NewProxy("1.1.1.1:53", "", "tcp", 8181, nil)

	log.Fatal(proxy.Serve())
}
