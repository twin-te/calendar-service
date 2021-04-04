package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"
)

func ICSHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(os.Getenv("COOKIE_NAME"))
	if err == http.ErrNoCookie {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("failed to get cookie: %+v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	ctx := WithAPICookie(r.Context(), cookie.String())

	modules, err := GetSchoolCalendar(ctx, 2021)
	if err != nil {
		log.Printf("failed to get school calendar: %+v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	courses, err := GetCourses(ctx, 2021)
	if err != nil {
		log.Printf("failed to get courses: %+v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var resp bytes.Buffer
	WriteICalendar(&resp, modules, courses)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(resp.Len()))
	_, err = resp.WriteTo(w)
	if err != nil {
		log.Printf("failed to write response: %+v", err)
	}
}
