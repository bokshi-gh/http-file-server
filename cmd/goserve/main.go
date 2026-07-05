package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/bokshi-gh/http-file-server/internal/handlers"
)

func main() {
	rootDir := flag.String("root", ".", "Root directory to serve")
	host := flag.String("host", "0.0.0.0", "Host to bind the server")
	port := flag.String("port", "8080", "Port to run the server on")
	verbose := flag.Bool("v", false, "Enable verbose logging")
	flag.Parse()

	if _, err := os.Stat(*rootDir); os.IsNotExist(err) {
		log.Fatalf("Root directory does not exist: %s", *rootDir)
	}

	http.HandleFunc("/", handlers.ClientHandler(*rootDir, *verbose))
	addr := *host + ":" + *port
	log.Printf("GoServe is listening on: http://%s", addr)
	log.Printf("Serving root: %s", *rootDir)
	log.Fatal(http.ListenAndServe(addr, nil))
}
