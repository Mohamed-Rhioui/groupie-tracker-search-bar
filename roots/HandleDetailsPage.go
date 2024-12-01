package roots

import (
	"encoding/json"
	"net/http"
	"strings"
	"text/template"

	"groupieTracker/tools"
)

// handler of Details Page (details.html)
func HandleDetailsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/details.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl2, err := template.ParseFiles("templates/error.html")
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
	var str []string
	id := r.URL.Query().Get("ID")
	str = strings.Split(id, "/")
	if len(str) > 1 {
		Err := tools.Errors{
			Message: "Page not Found",
		}
		w.WriteHeader(http.StatusNotFound)
		tmpl2.Execute(w, Err)
		return
	}

	if !tools.ValidateID(id) {
		Err := tools.Errors{
			Message: "Invalid ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		tmpl2.Execute(w, Err)
		return

	}

	resp, err := http.Get(tools.Url + "artists/" + id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Err := tools.Errors{
			Message: "Failed to fetch artist data",
		}
		w.WriteHeader(http.StatusBadGateway)
		tmpl2.Execute(w, Err)
		return
	}

	var artist tools.Artist
	err = json.NewDecoder(resp.Body).Decode(&artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var locationsData tools.Locations
	err = tools.FetchData(tools.Url+"locations/"+id, &locationsData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var datesData tools.ConcertDates
	err = tools.FetchData(tools.Url+"dates/"+id, &datesData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var relationsData tools.Relations
	err = tools.FetchData(tools.Url+"relation/"+id, &relationsData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := struct {
		Artist       tools.Artist
		Locations    tools.Locations
		ConcertDates tools.ConcertDates
		Relations    tools.Relations
	}{
		Artist:       artist,
		Locations:    locationsData,
		ConcertDates: datesData,
		Relations:    relationsData,
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
