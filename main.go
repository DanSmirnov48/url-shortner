package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

type PageData struct {
	OriginalUrl string
	ShortUrl    string
	ErrorMsg    string
}

var urlStore = make(map[string]string)

func main() {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			urlInput := r.FormValue("url")

			if !isValidUrl(urlInput) {
				err := tmpl.ExecuteTemplate(w, "index.html", PageData{
					OriginalUrl: urlInput,
					ErrorMsg:    "Invalid URL. Please enter a valid URL starting with http or https.",
				})
				if err != nil {
					http.Error(w, "Error rendering template", http.StatusInternalServerError)
				}
				return
			}

			if !isUrlReachable(urlInput) {
				err := tmpl.ExecuteTemplate(w, "index.html", PageData{
					OriginalUrl: urlInput,
					ErrorMsg:    "The URL does not exist or cannot be reached. Please try a different URL.",
				})
				if err != nil {
					http.Error(w, "Error rendering template", http.StatusInternalServerError)
				}
				return
			}

			shortKey, _ := generateToken()
			urlStore[shortKey] = urlInput

			shortUrl := "http://localhost:8080/" + shortKey

			err := tmpl.ExecuteTemplate(w, "index.html", PageData{
				OriginalUrl: urlInput,
				ShortUrl:    shortUrl,
				ErrorMsg:    "",
			})
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
		} else if r.Method == http.MethodGet {
			err := tmpl.ExecuteTemplate(w, "index.html", PageData{})
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
		}
	})

	router.HandleFunc("/{shortKey}", func(w http.ResponseWriter, r *http.Request) {
		shortKey := r.URL.Path[1:]

		originalUrl, exists := urlStore[shortKey]
		if !exists {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, originalUrl, http.StatusFound)
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

func isValidUrl(toTest string) bool {
	u, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

func isUrlReachable(testUrl string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(testUrl)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true
	}
	return false
}

func generateToken() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)

	return token, nil
}
