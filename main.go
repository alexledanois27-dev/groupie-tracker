package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	serverAddress   = "127.0.0.1:8080"
	shutdownTimeout = 5 * time.Second
	apiTimeout      = 10 * time.Second
)

var apiURLs = map[string]string{
	"/api/artists":   "https://groupietrackers.herokuapp.com/api/artists",
	"/api/locations": "https://groupietrackers.herokuapp.com/api/locations",
	"/api/dates":     "https://groupietrackers.herokuapp.com/api/dates",
	"/api/relation":  "https://groupietrackers.herokuapp.com/api/relation",
}

type Artist struct {
	ID           int    `json:"id"`
	Image        string `json:"image"`
	Name         string `json:"name"`
	CreationDate int    `json:"creationDate"`
	FirstAlbum   string `json:"firstAlbum"`
}

func newServer() (*http.Server, error) {
	indexTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	apiClient := &http.Client{Timeout: apiTimeout}

	mux.HandleFunc("/", indexHandler(apiClient, indexTemplate))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	for route, apiURL := range apiURLs {
		mux.HandleFunc(route, apiProxyHandler(apiClient, apiURL))
	}

	return &http.Server{
		Addr:              serverAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}, nil
}

func indexHandler(client *http.Client, indexTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, apiURLs["/api/artists"], nil)
		if err != nil {
			log.Printf("Impossible de créer la requête artists : %v", err)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		response, err := client.Do(req)
		if err != nil {
			log.Printf("API artists indisponible : %v", err)
			http.Error(w, "API distante indisponible", http.StatusBadGateway)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Printf("L'API artists a répondu avec le statut %s", response.Status)
			http.Error(w, "Erreur de l'API distante", http.StatusBadGateway)
			return
		}

		var artists []Artist
		if err := json.NewDecoder(response.Body).Decode(&artists); err != nil {
			log.Printf("Réponse artists invalide : %v", err)
			http.Error(w, "Réponse invalide de l'API distante", http.StatusBadGateway)
			return
		}

		var page bytes.Buffer
		if err := indexTemplate.ExecuteTemplate(&page, "index.html", artists); err != nil {
			log.Printf("Impossible de générer la page d'accueil : %v", err)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := page.WriteTo(w); err != nil {
			log.Printf("Impossible d'envoyer la page d'accueil : %v", err)
		}
	}
}

func apiProxyHandler(client *http.Client, apiURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		requestURL := apiURL
		if r.URL.RawQuery != "" {
			requestURL += "?" + r.URL.RawQuery
		}

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, requestURL, nil)
		if err != nil {
			log.Printf("Impossible de créer la requête vers %s : %v", apiURL, err)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		response, err := client.Do(req)
		if err != nil {
			log.Printf("API indisponible (%s) : %v", apiURL, err)
			http.Error(w, "API distante indisponible", http.StatusBadGateway)
			return
		}
		defer response.Body.Close()

		if contentType := response.Header.Get("Content-Type"); contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(response.StatusCode)

		if _, err := io.Copy(w, response.Body); err != nil {
			log.Printf("Impossible de transmettre la réponse de %s : %v", apiURL, err)
		}
	}
}

func startServer(server *http.Server) <-chan error {
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Serveur démarré sur http://%s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	return serverErrors
}

func stopServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Println("Arrêt du serveur...")
	return server.Shutdown(ctx)
}

func main() {
	server, err := newServer()
	if err != nil {
		log.Fatalf("Impossible de charger les templates : %v", err)
	}
	serverErrors := startServer(server)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stopSignal)

	select {
	case <-stopSignal:
		if err := stopServer(server); err != nil {
			log.Fatalf("Impossible d'arrêter proprement le serveur : %v", err)
		}
		log.Println("Serveur arrêté.")
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erreur du serveur : %v", err)
		}
	}
}
