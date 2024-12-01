package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// Define the structure that matches the JSON response
type Artist struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Image           string   `json:"image"`
	CreationDate    int      `json:"creationDate"`
	Members         []string `json:"members"`
	FirstAlbum      string   `json:"firstAlbum"`
	LocationsURL    string   `json:"locations"`
	ConcertDatesURL string   `json:"concertDates"`
	RelationsURL    string   `json:"relations"`
}

type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type ConcertDates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relations struct {
	ID             int                 `json:"id"`
	RelatedArtists map[string][]string `json:"datesLocations"`
}

type Errors struct {
	Message string
}

var (
	Url = "https://groupietrackers.herokuapp.com/api/"
)

// Define a single function to fetch data from a URL and unmarshal it into a struct
func FetchData(url string, dst interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dst)
}

func ValidateID(id string) bool {
	pattern := `^[0-9]{1,2}$`
	matched, err := regexp.MatchString(pattern, id)
	return err == nil && matched
}
