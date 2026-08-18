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

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type ConcertDates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type LocationsResponse struct {
	Index []Location `json:"index"`
}

type DatesResponse struct {
	Index []ConcertDates `json:"index"`
}

type RelationsResponse struct {
	Index []Relation `json:"index"`
}

func GetArtists(ctx context.Context, client *http.Client) ([]Artist, error) {
	return getArtistsFromURL(ctx, client, artistsURL)
}

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

func GetLocations(ctx context.Context, client *http.Client) ([]Location, error) {
	return getLocationsFromURL(ctx, client, locationsURL)
}

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

func GetConcertDates(ctx context.Context, client *http.Client) ([]ConcertDates, error) {
	return getConcertDatesFromURL(ctx, client, datesURL)
}

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

func GetRelations(ctx context.Context, client *http.Client) ([]Relation, error) {
	return getRelationsFromURL(ctx, client, relationsURL)
}

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
