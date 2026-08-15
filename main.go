package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func ollamaHealthCheck(ollamaBaseURL string) string {
	ollamaHealthCheck, err := http.Get(fmt.Sprintf("%s", ollamaBaseURL))
	if err != nil {
		log.Fatal(err)
	}
	ollamaHealthCheckBody, err := io.ReadAll(ollamaHealthCheck.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(ollamaHealthCheckBody)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("API_KEY")
	ollamaBaseURL := os.Getenv("OLLAMA_URL")
	exposeURL := os.Getenv("EXPOSE_URL")
	exposePort := os.Getenv("EXPOSE_PORT")

	var ollamaEndpoints = map[string]map[string]bool{
		"GET": {
			"/":            true,
			"/api/tags":    true,
			"/api/ps":      true,
			"/api/version": true,
		},
		"POST": {
			"/api/generate": true,
			"/api/chat":     true,
			"/api/embed":    true,
			"/api/show":     true,
			"/api/create":   true,
			"/api/copy":     true,
			"/api/pull":     true,
			"/api/push":     true,
		},
		"DELETE": {
			"/api/delete": true,
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		if !ollamaEndpoints[r.Method][r.URL.Path] {
			http.Error(w, "Endpoint not allowed", http.StatusNotFound)
			return
		} else {
			headerToken := r.Header.Get("Authorization")
			if headerToken == apiKey {
				requestURL := fmt.Sprintf("%s%s", ollamaBaseURL, r.URL.Path)
				repassReq, err := http.NewRequest(r.Method, requestURL, r.Body)
				if err != nil {
					log.Fatal(err)
				}
				repassReq.Header = r.Header.Clone()

				resp, err := http.DefaultClient.Do(repassReq)
				if err != nil {
					log.Fatal(err)
				}
				w.WriteHeader(resp.StatusCode)
				io.Copy(w, resp.Body)
				{
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	})

	err = http.ListenAndServe(fmt.Sprintf("%s:%s", exposeURL, exposePort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
