package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/marekforys/bridge-bid-tutor-go/internal/server"
)

// Minimal Swagger UI HTML using CDN, pointing to /docs/openapi.yaml
const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Bridge Bid Tutor API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style> body { margin:0 } </style>
  <script defer src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script defer src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.addEventListener('load', () => {
      window.ui = SwaggerUIBundle({
        url: '/docs/openapi.yaml',
        dom_id: '#swagger',
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
        layout: 'StandaloneLayout'
      });
    });
  </script>
  </head>
  <body>
    <div id="swagger"></div>
  </body>
</html>`

func main() {
	s := server.New()

	mux := http.NewServeMux()
	s.RegisterRoutes(mux)

	// Swagger UI (served at /docs) and OpenAPI spec (/docs/openapi.yaml)
	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/docs" { // keep /docs exact here
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, swaggerHTML)
	})
	// Serve the raw OpenAPI YAML
	mux.Handle("/docs/openapi.yaml", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))

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
