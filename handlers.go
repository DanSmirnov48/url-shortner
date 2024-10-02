package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var urlStore = make(map[string]string) // map to store short key -> original URL mapping

// handleHome handles the GET and POST requests for the main form
func handleHome(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		urlInput := r.FormValue("url")

		// Validate the URL format
		if !IsValidUrl(urlInput) {
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
		if !IsUrlReachable(urlInput) {
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

		// Generate a short key and store the URL
		shortKey, _ := GenerateToken()
		urlStore[shortKey] = urlInput

		// Shortened URL
		shortUrl := "http://localhost:8080/" + shortKey

		// Render the template with original and short URL
		err := tmpl.ExecuteTemplate(w, "index.html", PageData{
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
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[1:]

	// Look up the original URL using the short key
	originalUrl, exists := urlStore[shortKey]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalUrl, http.StatusFound)
}
