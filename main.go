package main

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	serverAddress   = "127.0.0.1:8080"
	shutdownTimeout = 5 * time.Second
	apiTimeout      = 10 * time.Second
)

// DetailPage contains the artist and concert data rendered on a detail page.
type DetailPage struct {
	Artist   Artist
	Relation Relation
}

// IndexPage contains the filtered artists and the current search term.
type IndexPage struct {
	Artists []Artist
	Search  string
}

// formatLocation converts API location keys into human-readable labels.
func formatLocation(location string) string {
	location = strings.ReplaceAll(location, "_", " ")
	location = strings.ReplaceAll(location, "-", ", ")
	return strings.Title(location)
}

// filterArtists matches the search term against group and member names.
func filterArtists(artists []Artist, search string) []Artist {
	search = strings.ToLower(strings.TrimSpace(search))
	if search == "" {
		return artists
	}

	filtered := make([]Artist, 0)
	for _, artist := range artists {
		matches := strings.Contains(strings.ToLower(artist.Name), search)
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), search) {
				matches = true
				break
			}
		}

		if matches {
			filtered = append(filtered, artist)
		}
	}

	return filtered
}

// newServer builds the HTTP server and registers the application routes.
func newServer() (*http.Server, error) {
	indexTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return nil, err
	}

	funcMap := template.FuncMap{
		"formatLocation": formatLocation,
	}
	detailsTemplate, err := template.New("details.html").Funcs(funcMap).ParseFiles("templates/details.html")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	apiClient := &http.Client{Timeout: apiTimeout}

	mux.HandleFunc("/", indexHandler(apiClient, indexTemplate))
	mux.HandleFunc("/artist/", detailHandler(apiClient, detailsTemplate))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/api/artists", apiProxyHandler(apiClient, artistsURL))
	mux.HandleFunc("/api/locations", apiProxyHandler(apiClient, locationsURL))
	mux.HandleFunc("/api/dates", apiProxyHandler(apiClient, datesURL))
	mux.HandleFunc("/api/relation", apiProxyHandler(apiClient, relationsURL))

	return &http.Server{
		Addr:              serverAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}, nil
}

// indexHandler fetches the artists and renders the home page.
func indexHandler(client *http.Client, indexTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		artists, err := GetArtists(r.Context(), client)
		if err != nil {
			log.Printf("Artists API unavailable: %v", err)
			http.Error(w, "Remote API unavailable", http.StatusBadGateway)
			return
		}

		search := strings.TrimSpace(r.URL.Query().Get("search"))
		pageData := IndexPage{
			Artists: filterArtists(artists, search),
			Search:  search,
		}

		var page bytes.Buffer
		if err := indexTemplate.ExecuteTemplate(&page, "index.html", pageData); err != nil {
			log.Printf("Failed to render the home page: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := page.WriteTo(w); err != nil {
			log.Printf("Failed to send the home page: %v", err)
		}
	}
}

// detailHandler renders the artist and concert information for one artist ID.
func detailHandler(client *http.Client, detailsTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idText := strings.TrimPrefix(r.URL.Path, "/artist/")
		id, err := strconv.Atoi(idText)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		artists, err := GetArtists(r.Context(), client)
		if err != nil {
			log.Printf("Artists API unavailable: %v", err)
			http.Error(w, "Remote API unavailable", http.StatusBadGateway)
			return
		}

		relations, err := GetRelations(r.Context(), client)
		if err != nil {
			log.Printf("Relations API unavailable: %v", err)
			http.Error(w, "Remote API unavailable", http.StatusBadGateway)
			return
		}

		if id < 1 || id > len(artists) || id > len(relations) {
			http.NotFound(w, r)
			return
		}

		pageData := DetailPage{
			Artist:   artists[id-1],
			Relation: relations[id-1],
		}

		var page bytes.Buffer
		if err := detailsTemplate.ExecuteTemplate(&page, "details.html", pageData); err != nil {
			log.Printf("Failed to render artist detail page: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := page.WriteTo(w); err != nil {
			log.Printf("Failed to send artist detail page: %v", err)
		}
	}
}

// apiProxyHandler forwards GET requests to one Groupie Trackers API endpoint.
func apiProxyHandler(client *http.Client, apiURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Preserve query parameters when forwarding the request.
		requestURL := apiURL
		if r.URL.RawQuery != "" {
			requestURL += "?" + r.URL.RawQuery
		}

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, requestURL, nil)
		if err != nil {
			log.Printf("Failed to create request to %s: %v", apiURL, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response, err := client.Do(req)
		if err != nil {
			log.Printf("API unavailable (%s): %v", apiURL, err)
			http.Error(w, "Remote API unavailable", http.StatusBadGateway)
			return
		}
		defer response.Body.Close()

		if contentType := response.Header.Get("Content-Type"); contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(response.StatusCode)

		if _, err := io.Copy(w, response.Body); err != nil {
			log.Printf("Failed to forward response from %s: %v", apiURL, err)
		}
	}
}

// startServer starts listening in a goroutine and reports the final server error.
func startServer(server *http.Server) <-chan error {
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Serveur démarré sur http://%s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	return serverErrors
}

// stopServer gives active requests a limited amount of time to finish.
func stopServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Println("Arrêt du serveur...")
	return server.Shutdown(ctx)
}

func main() {
	server, err := newServer()
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}
	serverErrors := startServer(server)

	// Handle both interactive interrupts and container termination signals.
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stopSignal)

	select {
	case <-stopSignal:
		if err := stopServer(server); err != nil {
			log.Fatalf("Failed to gracefully stop the server: %v", err)
		}
		log.Println("Serveur arrêté.")
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}
}
