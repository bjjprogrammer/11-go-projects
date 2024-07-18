package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string   `json:"id"`
	Isbn     string   `json:"isbn"`
	Title    string   `json:"title"`
	Director Director `json:"director"`
}

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)

}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)

}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "448743", Title: "Movie One", Director: Director{FirstName: "John", LastName: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "448744", Title: "Movie Two", Director: Director{FirstName: "Steve", LastName: "Smith"}})
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id:[0-9]+}", getMovie).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id:[0-9]+}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id:[0-9]+}", deleteMovie).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
