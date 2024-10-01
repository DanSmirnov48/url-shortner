package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	Name string
}

func main() {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := tmpl.ExecuteTemplate(w, "index.html", PageData{
			Name: "Joe",
		})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			fmt.Println("Template execution error:", err)
		}
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting website at localhost:8080")

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("An error occured:", err)
	}
}
