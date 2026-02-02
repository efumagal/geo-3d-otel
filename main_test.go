package main
package main

import (
    "encoding/json"
    "net/http/httptest"
    "testing"
    "time"
)

func TestRandFloatBounds(t *testing.T) {
    min, max := -5.0, 10.0
    for i := 0; i < 100; i++ {
        v := randFloat(min, max)
        if v < min || v > max {
            t.Fatalf("randFloat returned out of bounds value: %v not in [%v,%v]", v, min, max)
        }
    }
}

func TestDistanceHandler(t *testing.T) {
    app := setupApp()

    req := httptest.NewRequest("GET", "/distance", nil)
    // Fiber's app.Test takes a timeout duration
    resp, err := app.Test(req, 5*time.Second)
    if err != nil {
        t.Fatalf("failed to run request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        t.Fatalf("unexpected status code: %d", resp.StatusCode)
    }

    var body struct{
        Distance float64 `json:"distance"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
        t.Fatalf("failed to decode body: %v", err)
    }

    if body.Distance <= 0 {
        t.Fatalf("expected positive distance, got %v", body.Distance)
    }
}
