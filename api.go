package main

import (
	"context"
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

// Artist describes an artist returned by the Groupie Trackers API.
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

// Location contains all known concert locations for one artist.
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

// ConcertDates contains all known concert dates for one artist.
type ConcertDates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Relation maps each concert location to its dates for one artist.
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// LocationsResponse represents the top-level locations API payload.
type LocationsResponse struct {
	Index []Location `json:"index"`
}

// DatesResponse represents the top-level concert dates API payload.
type DatesResponse struct {
	Index []ConcertDates `json:"index"`
}

// RelationsResponse represents the top-level relations API payload.
type RelationsResponse struct {
	Index []Relation `json:"index"`
}

// GetArtists retrieves all artists from the official API.
func GetArtists(ctx context.Context, client *http.Client) ([]Artist, error) {
	return getArtistsFromURL(ctx, client, artistsURL)
}

// getArtistsFromURL performs the artist request against a configurable URL.
func getArtistsFromURL(ctx context.Context, client *http.Client, url string) ([]Artist, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artists API returned status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	if err := json.Unmarshal(data, &artists); err != nil {
		return nil, err
	}

	return artists, nil
}

// GetLocations retrieves the concert locations for all artists.
func GetLocations(ctx context.Context, client *http.Client) ([]Location, error) {
	return getLocationsFromURL(ctx, client, locationsURL)
}

// getLocationsFromURL performs the location request against a configurable URL.
func getLocationsFromURL(ctx context.Context, client *http.Client, url string) ([]Location, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("locations API returned status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response LocationsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Index, nil
}

// GetConcertDates retrieves the concert dates for all artists.
func GetConcertDates(ctx context.Context, client *http.Client) ([]ConcertDates, error) {
	return getConcertDatesFromURL(ctx, client, datesURL)
}

// getConcertDatesFromURL performs the concert date request against a configurable URL.
func getConcertDatesFromURL(ctx context.Context, client *http.Client, url string) ([]ConcertDates, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("concert dates API returned status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response DatesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Index, nil
}

// GetRelations retrieves the location-to-date relations for all artists.
func GetRelations(ctx context.Context, client *http.Client) ([]Relation, error) {
	return getRelationsFromURL(ctx, client, relationsURL)
}

// getRelationsFromURL performs the relation request against a configurable URL.
func getRelationsFromURL(ctx context.Context, client *http.Client, url string) ([]Relation, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("relations API returned status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response RelationsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Index, nil
}
