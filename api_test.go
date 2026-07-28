package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetArtist vérifie que la fonction getArtistsFromURL peut récupérer et décoder correctement les données d'un artiste depuis un serveur HTTP simulé.
func TestGetArtist(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// Simuler une réponse JSON pour l'API des artistes
			w.WriteHeader(http.StatusOK)

			// Écrire la réponse JSON dans le corps de la réponse HTTP
			_, err := w.Write([]byte(`[
				{
					"id": 1,
					"image": "https://example.com/queen.jpg",
					"name": "Queen",
					"members": ["Freddie Mercury", "Brian May"],
					"creationDate": 1970,
					"firstAlbum": "14-12-1973",
					"locations": "https://example.com/locations/1",
					"concertDates": "https://example.com/dates/1",
					"relations": "https://example.com/relation/1"
				}
			]`))

			if err != nil {
				t.Fatalf("could not write fake API response: %v", err)
			}
		},
	))
	defer server.Close()

	artists, err := getArtistsFromURL(server.URL)
	if err != nil {
		t.Fatalf("getArtistsFromURL failed: %v", err)
	}

	if len(artists) != 1 {
		t.Fatalf("expected 1 artist, got %d", len(artists))
	}

	artist := artists[0]

	if artist.ID != 1 {
		t.Errorf("expected artist ID 1, got %d", artist.ID)
	}

	if artist.Name != "Queen" {
		t.Errorf("expected artist name Queen, got %q", artist.Name)
	}

	if artist.CreationDate != 1970 {
		t.Errorf("expected creation date 1970, got %d", artist.CreationDate)
	}

	if len(artist.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(artist.Members))
	}

	if artist.Members[0] != "Freddie Mercury" {
		t.Errorf("expected first member Freddie Mercury, got %q", artist.Members[0])
	}
}

// TestGetArtistsHTTPError vérifie que la fonction getArtistsFromURL gère correctement les erreurs HTTP en renvoyant une erreur lorsque le serveur retourne un code de statut 500.
func TestGetArtistsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer server.Close()

	artists, err := getArtistsFromURL(server.URL)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if artists != nil {
		t.Errorf("expected nil artists, got %v", artists)
	}
}

// TestGetArtistsInvalidJSON vérifie que la fonction getArtistsFromURL gère correctement les erreurs de décodage JSON en renvoyant une erreur lorsque le serveur retourne un JSON invalide.
func TestGetArtistsInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`[
				{
					"id": 1,
					"name": "Queen"
			]`))
			// accolade manquante exprèsément pour provoquer une erreur de décodage JSON
			if err != nil {
				t.Fatalf("could not write fake API response: %v", err)
			}
		},
	))
	defer server.Close()

	artists, err := getArtistsFromURL(server.URL)

	if err == nil {
		t.Fatal("expected JSON decoding error, got nil")
	}

	if artists != nil {
		t.Errorf("expected nil artists, got %v", artists)
	}
}
