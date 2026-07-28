package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	artistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	locationsURL = "https://groupietrackers.herokuapp.com/api/locations"
	datesURL     = "https://groupietrackers.herokuapp.com/api/dates"
	relationsURL = "https://groupietrackers.herokuapp.com/api/relation"
)

// =========================
// Structures de l'API
// =========================

// Représente un artiste.
type Artist struct {
	ID              int      `json:"id"`
	Image           string   `json:"image"`
	Name            string   `json:"name"`
	Members         []string `json:"members"`
	CreationDate    int      `json:"creationDate"`
	FirstAlbum      string   `json:"firstAlbum"`
	LocationsURL    string   `json:"locations"`
	ConcertDatesURL string   `json:"concertDates"`
	RelationsURL    string   `json:"relations"`
}

// Représente les lieux de concert d'un artiste.
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

// Représente les dates de concert d'un artiste.
type ConcertDates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Représente l'association lieu -> dates.
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// =========================
// Structures des réponses API
// =========================
//
// Les endpoints locations, dates et relation
// renvoient un objet contenant une clé "index".
// Ces structures servent à décoder ces réponses.
//

type LocationsResponse struct {
	Index []Location `json:"index"`
}

type DatesResponse struct {
	Index []ConcertDates `json:"index"`
}

type RelationsResponse struct {
	Index []Relation `json:"index"`
}

// =========================
// Récupération des artistes
// =========================

func GetArtists() ([]Artist, error) {

	// Requête HTTP vers l'API.
	resp, err := http.Get(artistsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // ferme le corps de la réponse et libère les ressources une fois la fonction terminée.

	// Vérification du code de statut HTTP.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}

	// Lecture du corps de la réponse + sauvegarde en bytes dans data.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Conversion du JSON en slice de Artist.
	var artists []Artist

	err = json.Unmarshal(data, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// =========================
// Récupération des locations
// =========================

func GetLocations() ([]Location, error) {

	resp, err := http.Get(locationsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Décodage de la réponse JSON.
	var response LocationsResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response.Index, nil
}

// =========================
// Récupération des dates
// =========================

func GetConcertDates() ([]ConcertDates, error) {

	resp, err := http.Get(datesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response DatesResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response.Index, nil
}

// =========================
// Récupération des relations
// =========================

func GetRelations() ([]Relation, error) {

	resp, err := http.Get(relationsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response RelationsResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response.Index, nil
}
