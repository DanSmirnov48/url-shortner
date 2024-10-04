package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

// handleHome handles the GET and POST requests for the main form
func handleHome(w http.ResponseWriter, r *http.Request, tmpl *template.Template, db *sql.DB) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		urlInput := r.FormValue("url")

		// Validate the URL format
		if !isValidUrl(urlInput) {
			err := tmpl.ExecuteTemplate(w, "index.html", PageData{
				OriginalUrl: urlInput,
				ErrorMsg:    "Invalid URL. Please enter a valid URL starting with http or https.",
			})
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				fmt.Println("Template execution error:", err)
			}
			return
		}

		// Check if the URL is reachable
		if !isUrlReachable(urlInput) {
			err := tmpl.ExecuteTemplate(w, "index.html", PageData{
				OriginalUrl: urlInput,
				ErrorMsg:    "The URL does not exist or cannot be reached. Please try a different URL.",
			})
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				fmt.Println("Template execution error:", err)
			}
			return
		}

		// Generate a short key and store the URL in the database
		shortKey, _ := generateToken()
		err := insertURLMapping(db, shortKey, urlInput)
		if err != nil {
			http.Error(w, "Failed to store URL", http.StatusInternalServerError)
			fmt.Println("Database insertion error:", err)
			return
		}

		// Shortened URL
		shortUrl := "http://localhost:8080/" + shortKey

		// Render the template with original and short URL
		err = tmpl.ExecuteTemplate(w, "index.html", PageData{
			OriginalUrl: urlInput,
			ShortUrl:    shortUrl,
			ErrorMsg:    "",
		})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			fmt.Println("Template execution error:", err)
		}
	} else if r.Method == http.MethodGet {
		// Render the form for URL submission
		err := tmpl.ExecuteTemplate(w, "index.html", PageData{})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			fmt.Println("Template execution error:", err)
		}
	}
}

// handleRedirect handles the short URL redirection
func handleRedirect(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	shortKey := r.URL.Path[1:]

	// Look up the original URL using the short key from the database
	originalUrl, err := getOriginalURL(db, shortKey)
	if err != nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalUrl, http.StatusFound)
}
