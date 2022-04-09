package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetYear(r *http.Request) int {
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err == nil {
		return year
	}
	now := time.Now()
	year = now.Year()
	if now.Month() <= time.March {
		year -= 1
	}
	return year
}

func FilterByTags(courses []Course, tags []string) []Course {
	m := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		m[t] = struct{}{}
	}
	filtered := make([]Course, 0, len(courses))
	for _, c := range courses {
		for _, t := range c.Tags {
			if _, ok := m[t]; ok {
				filtered = append(filtered, c)
				break
			}
		}
	}
	return filtered
}

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
	year := GetYear(r)

	modules, err := GetSchoolCalendar(ctx, year)
	if err != nil {
		log.Printf("failed to get school calendar: %+v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	courses, err := GetCourses(ctx, year)
	if err != nil {
		log.Printf("failed to get courses: %+v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Has("tags[]") {
		tags := r.URL.Query()["tags[]"]
		courses = FilterByTags(courses, tags)
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
