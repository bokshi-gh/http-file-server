package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ClientHandler(rootDir string, verbose bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path
		if requestedPath == "" || requestedPath == "/" {
			requestedPath = "/"
		}
		fullPath := filepath.Join(rootDir, filepath.Clean(requestedPath))

		// Prevent directory traversal
		absRoot, _ := filepath.Abs(rootDir)
		absPath, _ := filepath.Abs(fullPath)
		if len(absPath) < len(absRoot) || absPath[:len(absRoot)] != absRoot {
			http.Error(w, "Invalid path", http.StatusForbidden)
			return
		}

		info, err := os.Stat(fullPath)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Redirect to path with trailing slash for directories
		if info.IsDir() && !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
			return
		}

		if info.IsDir() {
			// Try to serve index.html if it exists
			indexPath := filepath.Join(fullPath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				serveFile(indexPath, verbose, w)
				return
			}

			// Otherwise, generate directory listing with links
			entries, err := os.ReadDir(fullPath)
			if err != nil {
				http.Error(w, "Failed to read directory", http.StatusInternalServerError)
				return
			}

			style := fmt.Sprintf(`
					     <style>
					     	* { margin: 0px; padding: 0px; box-sizing: border-box; }
							body { padding: 1rem; display: flex; flex-direction: column; gap: 1rem; font-size: 18px; }
							div { display: flex; flex-direction: column; gap: 0.5rem; }
							div p { display: flex; gap: 0.5rem; align-items: center; }
					     </style>
			         `)

			head := fmt.Sprintf(`
					     <head>
						     <meta charset=\"utf-8\"><title>Index of %s</title>
						     <meta name="viewport" content="width=device-width, initial-scale=1.0">
						     <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/7.0.1/css/all.min.css" integrity="sha512-2SwdPD6INVrV/lHTZbO2nodKhrnDdJK9/kg2XD1r9uGqPo1cUbujc+IYdlYdEErWNu69gVcYgdxlmVmzTWnetw==" crossorigin="anonymous" referrerpolicy="no-referrer" />
			 			     %s
				         </head>
			        `, requestedPath, style)

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "<html>%s<body><h1>Index of %s</h1><hr><div>", head, requestedPath)
			for _, e := range entries {
				name := e.Name()
				if e.IsDir() {
					fmt.Fprintf(w, `<p><i class="fa-solid fa-folder-open"></i><a href="%s/">%s</a></p>`, name, name)
					continue
				}
				fmt.Fprintf(w, `<p><i class="fa-solid fa-file"></i><a href="%s">%s</a></p>`, name, name)
			}
			fmt.Fprint(w, "</div></body></html>")

			if verbose {
				log.Printf("Listed directory: %s (%d entries)", fullPath, len(entries))
			}
			return
		}

		serveFile(fullPath, verbose, w)
	}
}

func serveFile(fullPath string, verbose bool, w http.ResponseWriter) {
	start := time.Now()
	
	mime := "application/octet-stream"
	ext := strings.ToLower(filepath.Ext(fullPath))
	switch ext {
	case ".txt":
		mime = "text/plain"
	case ".html":
		mime = "text/html"
	case ".css":
		mime = "text/css"
	case ".js":
		mime = "application/javascript"
	case ".png":
		mime = "image/png"
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".gif":
		mime = "image/gif"
	case ".svg":
		mime = "image/svg+xml"
	case ".json":
		mime = "application/json"
	case ".pdf":
		mime = "application/pdf"
	case ".mp4":
		mime = "video/mp4"
	case ".mp3":
		mime = "audio/mpeg"
	case ".zip":
		mime = "application/zip"
	}

	w.Header().Set("Content-Type", mime)

	f, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusNotFound)
		return
	}
	defer f.Close()

	info, _ := f.Stat()
	fileSize := info.Size()

	_, _ = io.Copy(w, f)

	elapsed := time.Since(start)
	if verbose {
		log.Printf(
			"Served file: %s (%d bytes, %s)",
			fullPath,
			fileSize,
			elapsed,
		)
	}
}
