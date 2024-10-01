package main

import (
	"html/template"
	"log"
	"net/http"
)

// Handler for the root URL "/"
func handler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Pass data to the template
	data := struct {
		Title string
	}{
		Title: "There",
	}

	// Render the template with the data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Register the handler for the root URL "/"
	http.HandleFunc("/", handler)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
