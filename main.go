package main

import (
	"encoding/json"
	"fmt"
	"github.com/edjmore/edbot/groupme"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/edbot", handleEdBot)
	http.HandleFunc("/edbot/avatar", handleEdBotAvatar)

	// When running on Heroku we must bind to the port specified in env.
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "5000"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

// The main EdBot endpoint; handles incoming group messages.
func handleEdBot(w http.ResponseWriter, r *http.Request) {
	if validateMethod(w, r, http.MethodGet, http.MethodPost) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprint(w, "Hello, I am EdBot.")
		case http.MethodPost:
			var m groupme.Message
			if r.Body == nil {
				http.Error(w, "Request has no body.", http.StatusBadRequest)
			} else if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				groupme.Respond(m)
			}
		}
	}
}

func handleEdBotAvatar(w http.ResponseWriter, r *http.Request) {
	if validateMethod(w, r, http.MethodGet) {
		http.ServeFile(w, r, "static/avatar.jpg")
	}
}

// Writes an appropriate error response if the incoming request uses an invalid method.
// Returns false if the request method is invalid (or OPTIONS).
func validateMethod(w http.ResponseWriter, r *http.Request, validMethods ...string) bool {
	for _, method := range validMethods {
		if method == r.Method {
			return true
		}
	}
	w.Header().Set("Allow", strings.Join(validMethods, ", "))
	if r.Method != http.MethodOptions {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	return false
}
