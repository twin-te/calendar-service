package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/timetable.ics", ICSHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "healthy", http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(":5000", mux))
}
