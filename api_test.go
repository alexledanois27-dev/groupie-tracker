package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//--------------------------
// Tests Artists
//--------------------------

// TestGetArtist vérifie que la fonction getArtistsFromURL peut récupérer et décoder correctement les données d'un artiste depuis un serveur HTTP simulé.
func TestGetArtistsSuccess(t *testing.T) {
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

//--------------------------
// Tests Locations
//--------------------------

func TestGetLocationsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// la réponse complète est un objet JSON avec un tableau "index" contenant les locations, pas une liste comme dans Artists
			_, err := w.Write([]byte(`{
				"index": [
					{
						"id": 1,
						"locations": [
							"paris-france",
							"london-uk"
						]
					}
				]
			}`))
			if err != nil {
				t.Fatalf("could not write fake API response: %v", err)
			}
		},
	))
	defer server.Close()

	locations, err := getLocationsFromURL(server.URL)
	if err != nil {
		t.Fatalf("getLocationsFromURL failed: %v", err)
	}
	if len(locations) != 1 {
		t.Fatalf("expected 1 location entry, got %d", len(locations))
	}

	location := locations[0]

	if location.ID != 1 {
		t.Errorf("expected location ID 1, got %d", location.ID)
	}

	if len(location.Locations) != 2 {
		t.Fatalf("expected 2 locations, got %d", len(location.Locations))
	}

	if location.Locations[0] != "paris-france" {
		t.Errorf(
			"expected first location paris-france, got %q",
			location.Locations[0],
		)
	}
}

func TestGetLocationsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer server.Close()

	locations, err := getLocationsFromURL(server.URL)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if locations != nil {
		t.Errorf("expected nil locations, got %v", locations)
	}
}

func TestGetLocationsInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{
				"index": [
					{
						"id": 1,
						"locations": ["paris-france"]
				]
			}`))
			if err != nil {
				t.Fatalf("could not write fake API response: %v", err)
			}
		},
	))
	defer server.Close()

	locations, err := getLocationsFromURL(server.URL)

	if err == nil {
		t.Fatal("expected JSON decoding error, got nil")
	}

	if locations != nil {
		t.Errorf("expected nil locations, got %v", locations)
	}
}

//--------------------------
// Tests Concert Dates
//--------------------------

func TestGetConcertDatesSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{
				"index": [
					{
						"id": 1,
						"dates": [
							"*23-08-2019",
							"*22-09-2019"
						]
					}
				]
			}`))
			if err != nil {
				t.Fatalf("could not write fake API response: %v", err)
			}
		},
	))
	defer server.Close()

	concertDates, err := getConcertDatesFromURL(server.URL)
	if err != nil {
		t.Fatalf("getConcertDatesFromURL failed: %v", err)
	}

	if len(concertDates) != 1 {
		t.Fatalf(
			"expected 1 concert dates entry, got %d",
			len(concertDates),
		)
	}

	dateEntry := concertDates[0]

	if dateEntry.ID != 1 {
		t.Errorf("expected concert dates ID 1, got %d", dateEntry.ID)
	}

	if len(dateEntry.Dates) != 2 {
		t.Fatalf("expected 2 dates, got %d", len(dateEntry.Dates))
	}

	if dateEntry.Dates[0] != "*23-08-2019" {
		t.Errorf(
			"expected first date *23-08-2019, got %q",
			dateEntry.Dates[0],
		)
	}
}

func TestGetConcertDatesHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer server.Close()

	concertDates, err := getConcertDatesFromURL(server.URL)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if concertDates != nil {
		t.Errorf("expected nil concert dates, got %v", concertDates)
	}
}

func TestGetConcertDatesInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{
				"index": [
					{
						"id": 1,
						"dates": ["*23-08-2019"]
				]
			}`))
			if err != nil {
				t.Fatalf("could not write fake API response: %v", err)
			}
		},
	))
	defer server.Close()

	concertDates, err := getConcertDatesFromURL(server.URL)

	if err == nil {
		t.Fatal("expected JSON decoding error, got nil")
	}

	if concertDates != nil {
		t.Errorf("expected nil concert dates, got %v", concertDates)
	}
}

//--------------------------
// Tests Relations
//--------------------------

func TestGetRelationsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{
				"index": [
					{
						"id": 1,
						"datesLocations": {
							"paris-france": [
								"23-08-2019",
								"24-08-2019"
							]
						}
					}
				]
			}`))
			if err != nil {
				t.Fatalf("could not write fake API response: %v", err)
			}
		},
	))
	defer server.Close()

	relations, err := getRelationsFromURL(server.URL)
	if err != nil {
		t.Fatalf("getRelationsFromURL failed: %v", err)
	}

	if len(relations) != 1 {
		t.Fatalf("expected 1 relation entry, got %d", len(relations))
	}

	relation := relations[0]

	if relation.ID != 1 {
		t.Errorf("expected relation ID 1, got %d", relation.ID)
	}

	dates, exists := relation.DatesLocations["paris-france"]
	if !exists {
		t.Fatal("expected paris-france relation to exist")
	}

	if len(dates) != 2 {
		t.Fatalf(
			"expected 2 dates for paris-france, got %d",
			len(dates),
		)
	}

	if dates[0] != "23-08-2019" {
		t.Errorf(
			"expected first relation date 23-08-2019, got %q",
			dates[0],
		)
	}
}

func TestGetRelationsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer server.Close()

	relations, err := getRelationsFromURL(server.URL)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if relations != nil {
		t.Errorf("expected nil relations, got %v", relations)
	}
}

func TestGetRelationsInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(`{
				"index": [
					{
						"id": 1,
						"datesLocations": {
							"paris-france": ["23-08-2019"]
						}
				]
			}`))
			if err != nil {
				t.Fatalf("could not write fake API response: %v", err)
			}
		},
	))
	defer server.Close()

	relations, err := getRelationsFromURL(server.URL)

	if err == nil {
		t.Fatal("expected JSON decoding error, got nil")
	}

	if relations != nil {
		t.Errorf("expected nil relations, got %v", relations)
	}
}
