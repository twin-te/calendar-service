package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/timetable.ics", ICSHandler)
	log.Fatal(http.ListenAndServe(":5000", mux))
}
