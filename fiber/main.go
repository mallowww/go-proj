package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// server := http.NewServeMux()
	server := mux.NewRouter()
	server.HandleFunc("/rentals/{id}", Rentals).Methods(http.MethodGet)
	http.ListenAndServe(":8080", server)

}

func Rentals(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
}
