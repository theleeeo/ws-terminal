package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func loadTemplate(filenames ...string) (*template.Template, error) {
	tmpl, err := template.New("webpage").ParseFiles(filenames...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func termHandler(tmpl *template.Template, baseUrl string) http.HandlerFunc {
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
			http.Redirect(w, r, fmt.Sprintf("/creds?id=%s", id), http.StatusSeeOther)
			return
		}

		data := struct {
			WebSocketURL string
		}{
			WebSocketURL: fmt.Sprintf("%s%s?credentials=%s", baseUrl, id, creds),
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Failed to execute template: %v", err)
		}
	}
}

func credHandler(tmpl *template.Template) http.HandlerFunc {
	if tmpl == nil {
		log.Fatal("credential template is nil")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		if err := tmpl.Execute(w, struct{}{}); err != nil {
			log.Printf("Failed to execute template: %v", err)
		}
	}
}

func main() {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		log.Fatal("BASE_URL is required")
		os.Exit(1)
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	tmpl, err := loadTemplate("public/terminal.html", "public/credentials.html")
	if err != nil {
		log.Fatalf("Failed to load template: %v", err)
	}

	http.HandleFunc("/exec/{id}", termHandler(tmpl.Lookup("terminal.html"), baseUrl))
	http.HandleFunc("/creds", credHandler(tmpl.Lookup("credentials.html")))
	log.Println("Server listening at", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
