// main_test.go
// Copyright 2025
// SPDX-License-Identifier: Apache-2.0
// Authors: Mihai Criveti
package main

import (
    "context"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "testing"
)

// --- Test helpers for mocking HTTP ---

// roundTripperFunc lets us replace http.DefaultTransport
type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
    return f(req)
}

// --- Tests ---

func TestGoogleSearchEnvCheck(t *testing.T) {
    // Ensure error when API key or CSE ID not set
    os.Unsetenv("GOOGLE_API_KEY")
    os.Unsetenv("GOOGLE_CSE_ID")

    _, err := googleSearch(context.Background(), json.RawMessage(`{"query":"test"}`))
    if err == nil || !strings.Contains(err.Error(), "GOOGLE_API_KEY and GOOGLE_CSE_ID must be set") {
        t.Fatalf("expected missing env error, got %v", err)
    }
}

func TestGoogleSearchMocked(t *testing.T) {
    // Provide dummy env vars
    os.Setenv("GOOGLE_API_KEY", "dummy-key")
    os.Setenv("GOOGLE_CSE_ID", "dummy-cx")

    // Swap in a fake transport
    originalTransport := http.DefaultTransport
    http.DefaultTransport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
        // Verify that the request URL includes our dummy key and cx
        if !strings.Contains(req.URL.String(), "key=dummy-key") ||
           !strings.Contains(req.URL.String(), "cx=dummy-cx") ||
           !strings.Contains(req.URL.RawQuery, "q=test") {
            t.Errorf("unexpected URL: %s", req.URL.String())
        }

        // Return a fake JSON response
        body := `{"items":[{"title":"Foo","link":"http://example.com"}]}`
        return &http.Response{
            StatusCode: 200,
            Body:       ioutil.NopCloser(strings.NewReader(body)),
            Header:     http.Header{"Content-Type": []string{"application/json"}},
        }, nil
    })
    defer func() { http.DefaultTransport = originalTransport }()

    // Call the function
    res, err := googleSearch(context.Background(), json.RawMessage(`{"query":"test"}`))
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    // Assert on returned type and content
    results, ok := res.([]SearchResult)
    if !ok {
        t.Fatalf("expected []SearchResult, got %T", res)
    }
    if len(results) != 1 {
        t.Fatalf("expected 1 result, got %d", len(results))
    }
    if results[0].Title != "Foo" || results[0].Link != "http://example.com" {
        t.Errorf("unexpected result: %+v", results[0])
    }
}
