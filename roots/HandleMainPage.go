package roots

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"groupieTracker/tools"
)

var (
	tmpl, tmpl2 *template.Template
	err         error
)

// handler of MainPage (index.html)
func HandleMainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl2, err = template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method != "GET" {
		Err := tools.Errors{
			Message: "Invalid request method.",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		tmpl2.Execute(w, Err)
		return
	}

	if r.URL.Path != "/" {
		Err := tools.Errors{
			Message: "Page not Found",
		}
		w.WriteHeader(http.StatusNotFound)
		tmpl2.Execute(w, Err)
		return
	}

	var artists []tools.Artist
	query := r.URL.Query().Get("search")

	err := fetchDataArtist(tools.Url+"artists", w, &artists)
	if err != nil {
		return
	}
	if query != "" {
		artists = searchArtists(artists, query, w)
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fetchDataArtist(url string, w http.ResponseWriter, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		Err := tools.Errors{
			Message: "Failed to fetch locations data",
		}
		w.WriteHeader(http.StatusBadGateway)
		tmpl2.Execute(w, Err)
		return nil
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil

}

func searchArtists(data []tools.Artist, query string, w http.ResponseWriter) []tools.Artist {
	query = strings.ToLower(query)
	var results []tools.Artist
	var NewResult []tools.Artist
	var relationsData tools.Relations

	for _, artist := range data {
		idStr := strconv.Itoa(artist.ID)
		creationYearStr := strconv.Itoa(artist.CreationDate)
		err = tools.FetchData(tools.Url+"relation/"+idStr, &relationsData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if artist.FirstAlbum == query {
			results = append(results, artist)
		}
		if strings.Contains(strings.ToLower(artist.Name), query) {
			results = append(results, artist)
		}
		// Search within members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				results = append(results, artist)
				break
			}
		}

		// Search within ID and creation date
		if idStr == query {
			results = append(results, artist)
		}
		if creationYearStr == query {
			results = append(results, artist)
		}

		// Search within relations (locations and dates)
		for location, dates := range relationsData.RelatedArtists {
			if strings.Contains(strings.ToLower(location), query) {
				results = append(results, artist)
				break
			}
			for _, date := range dates {
				if strings.Contains(date, query) {
					results = append(results, artist)
					break
				}
			}
		}
		relationsData = tools.Relations{}
	}

	for i := 0; i < len(results); i++ {
		if i == 0 {
			NewResult = append(NewResult, results[i])
		} else if reflect.DeepEqual(results[i], NewResult[len(NewResult)-1]) {
			continue
		} else {
			NewResult = append(NewResult, results[i])
		}
	}

	return NewResult
}
