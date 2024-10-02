package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	OriginalUrl string
	ShortUrl    string
}

func main() {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			url := r.FormValue("url")

			shortUrl := "https://short.ly/abcd1234"

			err := tmpl.ExecuteTemplate(w, "index.html", PageData{
				OriginalUrl: url,
				ShortUrl:    shortUrl,
			})
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				fmt.Println("Template execution error:", err)
			}
		} else if r.Method == http.MethodGet {
			err := tmpl.ExecuteTemplate(w, "index.html", PageData{
				OriginalUrl: "",
				ShortUrl:    "",
			})
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				fmt.Println("Template execution error:", err)
			}
		}
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting website at http://localhost:8080")

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("An error occurred:", err)
	}
}
