package roots

import (
	"net/http"
	"strings"
	"text/template"

	"groupieTracker/tools"
)

func Handlecss(w http.ResponseWriter, r *http.Request) {
	Err := tools.Errors{
		Message: "Page Not Found",
	}
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/style/") || r.URL.Path == "/style" {
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, Err)
		return
	}
}
