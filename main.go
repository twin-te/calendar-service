package main

import (
	"context"
	"encoding/json"
	"os"
)

func main() {
	sc, err := GetSchoolCalendar(context.Background(), 2021)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(os.Stdout).Encode(sc)
}
