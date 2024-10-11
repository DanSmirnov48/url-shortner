package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

// handleHome handles the GET request for rendering the main form
func handleHome(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == http.MethodGet {
		// Render the form for URL submission
		err := tmpl.ExecuteTemplate(w, "index.html", PageData{})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			fmt.Println("Template execution error:", err)
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleShorten(w http.ResponseWriter, r *http.Request, tmpl *template.Template, db *sql.DB) {
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

		// Check if the original URL has already been shortened
		existingShortKey, err := findURLMappingByOriginal(db, urlInput)
		if err != nil {
			http.Error(w, "Failed to check URL", http.StatusInternalServerError)
			fmt.Println("Database lookup error:", err)
			return
		}

		// If the URL is already shortened, use the existing short key
		var shortKey string
		if existingShortKey != "" {
			shortKey = existingShortKey
		} else {
			// If not, generate a new short key and store the URL in the database
			shortKey = generateShortKey()
			err := insertURLMapping(db, shortKey, urlInput)
			if err != nil {
				http.Error(w, "Failed to store URL", http.StatusInternalServerError)
				fmt.Println("Database insertion error:", err)
				return
			}
		}

		shortUrl := "http://localhost:8080/" + shortKey

		// Render the shorten.html template with the shortened URL
		err = tmpl.ExecuteTemplate(w, "shorten.html", PageData{
			OriginalUrl: urlInput,
			ShortUrl:    shortUrl,
		})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			fmt.Println("Template execution error:", err)
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
