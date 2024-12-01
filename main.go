package main

import (
	"fmt"
	"net/http"

	"groupieTracker/roots"
)

func main() {
	http.HandleFunc("/style", roots.Handlecss)
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))

	http.HandleFunc("/", roots.HandleMainPage)
	http.HandleFunc("/details", roots.HandleDetailsPage)

	fmt.Println("\x1b[92m" + " -----------------------------------------------")
	fmt.Println(" |  " + "\033[1m" + "im working, port :" + "\x1b[91m" + " http://localhost:7000" + "\x1b[0m" + "\x1b[92m" + "   |")
	fmt.Println(" -----------------------------------------------" + "\x1b[0m")
	http.ListenAndServe(":7000", nil)
}
