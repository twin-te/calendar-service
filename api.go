package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const APIBaseURL = "https://app.twinte.net/api/v3"

type apiCookieKey struct{}

func WithAPICookie(ctx context.Context, cookie string) context.Context {
	return context.WithValue(ctx, apiCookieKey{}, cookie)
}

func GetAPI(ctx context.Context, endpoint string, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", APIBaseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating api request: %w", err)
	}
	cookie, ok := ctx.Value(apiCookieKey{}).(string)
	if !ok {
		cookie = os.Getenv("TWINTE_COOKIE")
	}
	req.Header.Set("Cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("requesting api: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api: %s (%q)", resp.Status, string(b))
	}
	return json.NewDecoder(resp.Body).Decode(data)
}
