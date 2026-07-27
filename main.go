package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	serverAddress  = "127.0.0.1:8080"
	shutdownTimeout = 5 * time.Second
)

func newServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, "Le serveur fonctionne.")
	})

	return &http.Server{
		Addr:              serverAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
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
	server := newServer()
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
