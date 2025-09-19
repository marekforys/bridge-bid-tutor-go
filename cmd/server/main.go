package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/yourusername/bridge-bid-tutor-go/internal/server"
)

func main() {
	s := server.New()

	mux := http.NewServeMux()
	s.RegisterRoutes(mux)

	// Static file server for simple web client
	webDir := filepath.Clean("web")
	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})

	addr := ":8080"
	fmt.Printf("Bridge Bid Tutor REST server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
