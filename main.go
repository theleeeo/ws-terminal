package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

func loadTemplate(filenames ...string) (*template.Template, error) {
	tmpl, err := template.New("webpage").ParseFiles(filenames...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func termHandler(tmpl *template.Template, baseUrl string, pathPrefix string) http.HandlerFunc {
	if tmpl == nil {
		log.Fatal("terminal template is nil")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		creds := r.URL.Query().Get("credentials")
		if creds == "" {
			http.Redirect(w, r, fmt.Sprintf("%s/creds?id=%s", pathPrefix, id), http.StatusSeeOther)
			return
		}

		data := struct {
			WebSocketURL string
		}{
			WebSocketURL: fmt.Sprintf("%s/%s?credentials=%s", baseUrl, id, creds),
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Failed to execute template: %v", err)
		}
	}
}

func credHandler(tmpl *template.Template, pathPrefix string) http.HandlerFunc {
	if tmpl == nil {
		log.Fatal("credential template is nil")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		data := struct {
			PathPrefix string
		}{
			PathPrefix: pathPrefix,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Failed to execute template: %v", err)
		}
	}
}

func main() {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		log.Fatal("BASE_URL is required")
	}

	u, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatalf("Failed to parse BASE_URL: %v", err)
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		log.Fatal("BASE_URL must be a websocket URL (scheme ws or wss)")
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	pathPrefix := os.Getenv("PATH_PREFIX")

	tmpl, err := loadTemplate("public/terminal.html", "public/credentials.html")
	if err != nil {
		log.Fatalf("Failed to load template: %v", err)
	}

	http.HandleFunc(fmt.Sprintf("%s/exec/{id}", pathPrefix), termHandler(tmpl.Lookup("terminal.html"), baseUrl, pathPrefix))
	http.HandleFunc(fmt.Sprintf("%s/creds", pathPrefix), credHandler(tmpl.Lookup("credentials.html"), pathPrefix))
	log.Println("Server listening at", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
