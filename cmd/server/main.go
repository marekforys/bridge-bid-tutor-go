package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yourusername/bridge-bid-tutor-go/internal/server"
)

func main() {
	s := server.New()

	mux := http.NewServeMux()
	s.RegisterRoutes(mux)

	addr := ":8080"
	fmt.Printf("Bridge Bid Tutor REST server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
