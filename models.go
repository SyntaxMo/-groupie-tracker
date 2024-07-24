package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID        int                 `json:"id"`
	Relations map[string][]string `json:"datesLocations"`
}

type FinalArtist struct {
	ID           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    []string
	ConcertDates []string
	Relations    map[string][]string
}

func getArtists() []Artist {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Println("Error fetching artists:", err)
		return nil
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}

	var responseObject []Artist
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return nil
	}

	return responseObject
}

func getArtistByID(id string) *FinalArtist {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		log.Println("Error fetching artists by ID", err)
		return nil
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}

	var artist Artist
	err = json.Unmarshal(responseData, &artist)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
		return nil
	}

	if artist.ID == 0 {
		return nil
	}

	var FinalArtist FinalArtist
	FinalArtist.ID = artist.ID
	FinalArtist.Image = artist.Image
	FinalArtist.Name = artist.Name
	FinalArtist.Members = artist.Members
	FinalArtist.CreationDate = artist.CreationDate
	FinalArtist.FirstAlbum = artist.FirstAlbum
	FinalArtist.Locations = getLocations(id)
	FinalArtist.ConcertDates = getDates(id)
	FinalArtist.Relations = getRelations(id)

	return &FinalArtist
}

func getLocations(id string) []string {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/" + id)
	if err != nil {
		log.Println("Error fetching locations:", err)
		return nil
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}

	var locations Location
	err = json.Unmarshal(responseData, &locations)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
		return nil
	}

	for i := range locations.Locations {
		locations.Locations[i] = strings.ReplaceAll(locations.Locations[i], "_", " ")
		locations.Locations[i] = strings.ReplaceAll(locations.Locations[i], "-", ", ")
		locations.Locations[i] = strings.Title(locations.Locations[i])
	}

	return locations.Locations
}

func getDates(id string) []string {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/dates/" + id)
	if err != nil {
		log.Println("Error fetching dates:", err)
		return nil
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}

	var dates Date
	err = json.Unmarshal(responseData, &dates)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
		return nil
	}

	for i, date := range dates.Dates {
		dates.Dates[i] = strings.TrimPrefix(date, "*")
	}

	return dates.Dates
}

func getRelations(id string) map[string][]string {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + id)
	if err != nil {
		log.Println("Error fetching relations:", err)
		return nil
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}

	var relation Relation
	err = json.Unmarshal(responseData, &relation)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
		return nil
	}

	temp := make(map[string][]string)
	for key, value := range relation.Relations {
		key = strings.ReplaceAll(key, "_", " ")
		key = strings.ReplaceAll(key, "-", ", ")
		key = strings.Title(key)
		temp[key] = value
	}

	return temp
}
