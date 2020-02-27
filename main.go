package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Game struct {
	Name        string `json:"name"`
	ReleaseYear int    `json:"release_year"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func GamesHandler(w http.ResponseWriter, r *http.Request) {
	games := []Game{
		Game{
			Name:        "Sonic the Hedgehog 1",
			ReleaseYear: 1991,
		},
		Game{
			Name:        "Sonic the Hedgehog 2",
			ReleaseYear: 1992,
		},
		Game{
			Name:        "Sonic the Hedgehog 3",
			ReleaseYear: 1994,
		},
		Game{
			Name:        "Sonic & Kunckles",
			ReleaseYear: 1994,
		},
	}

	b, err := json.Marshal(games)
	if err != nil {
		panic(err)
	}

	w.Write(b)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/javascript")

	errResp := ErrorResponse{
		Error: "Resource could not be found.",
	}

	b, _ := json.Marshal(errResp)
	w.WriteHeader(http.StatusNotFound)
	w.Write(b)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/javascript")

	errResp := ErrorResponse{
		Error: "Method not allowed.",
	}

	b, _ := json.Marshal(errResp)
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(b)
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/javascript")
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	r.Use(ContentTypeMiddleware)

	r.HandleFunc("/games", GamesHandler).Methods("GET")

	r.NotFoundHandler = r.NewRoute().
		BuildOnly().
		HandlerFunc(NotFoundHandler).
		GetHandler()

	r.MethodNotAllowedHandler = r.NewRoute().
		BuildOnly().
		HandlerFunc(MethodNotAllowedHandler).
		GetHandler()

	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    "0.0.0.0:80",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
