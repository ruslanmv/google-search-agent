// main_test.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// mockRoundTripper lets us fake HTTP responses for testing the Google API call.
type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestGoogleSearchHandler_EnvCheck(t *testing.T) {
	// Ensure the handler returns an error when API key or CSE ID are not set.
	// t.Setenv is used for setting env vars in tests (Go 1.17+).
	t.Setenv("GOOGLE_API_KEY", "")
	t.Setenv("GOOGLE_CSE_ID", "")

	// Create a dummy MCP request.
	req := mcp.CallToolRequest{
		Arguments: json.RawMessage(`{"query":"test"}`),
	}

	// Call the handler.
	result, err := googleSearchHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler returned an unexpected error: %v", err)
	}

	// Check for the specific error message inside the tool result.
	if result.Error == nil || !strings.Contains(result.Error.Message, "must be set") {
		t.Fatalf("expected missing env error in tool result, got: %+v", result)
	}
}

func TestGoogleSearchHandler_MockedSuccess(t *testing.T) {
	// Provide dummy env vars for this test.
	t.Setenv("GOOGLE_API_KEY", "dummy-key")
	t.Setenv("GOOGLE_CSE_ID", "dummy-cx")

	// Set up the mocked HTTP client.
	originalClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify the request URL contains our dummy credentials and query.
				if !strings.Contains(req.URL.String(), "key=dummy-key") {
					t.Errorf("request URL missing 'key=dummy-key'")
				}
				if !strings.Contains(req.URL.String(), "cx=dummy-cx") {
					t.Errorf("request URL missing 'cx=dummy-cx'")
				}
				if !strings.Contains(req.URL.RawQuery, "q=test+query") {
					t.Errorf("request URL missing 'q=test+query', got: %s", req.URL.RawQuery)
				}

				// Return a fake JSON response from the Google API.
				body := `{"items":[{"title":"Test Title","link":"https://example.com"}]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
		},
	}
	// Restore the original client after the test.
	defer func() { http.DefaultClient = originalClient }()

	// Create the MCP request.
	req := mcp.CallToolRequest{
		Name:      "Google Search",
		Arguments: json.RawMessage(`{"query":"test query"}`),
	}

	// Call the handler.
	result, err := googleSearchHandler(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error from handler, got %v", err)
	}
	if result.Error != nil {
		t.Fatalf("expected no error in tool result, got message: %s", result.Error.Message)
	}

	// Assert on the returned content.
	var searchResults []SearchResult
	if err := json.Unmarshal([]byte(result.Text), &searchResults); err != nil {
		t.Fatalf("failed to unmarshal result text into []SearchResult: %v", err)
	}

	if len(searchResults) != 1 {
		t.Fatalf("expected 1 result, got %d", len(searchResults))
	}
	if searchResults[0].Title != "Test Title" || searchResults[0].Link != "https://example.com" {
		t.Errorf("unexpected result content: %+v", searchResults[0])
	}
}

func TestGoogleSearchHandler_ApiFailure(t *testing.T) {
	t.Setenv("GOOGLE_API_KEY", "dummy-key")
	t.Setenv("GOOGLE_CSE_ID", "dummy-cx")

	// Mock an API failure (e.g., 403 Forbidden).
	originalClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(strings.NewReader(`{"error":"invalid key"}`)),
				}, nil
			},
		},
	}
	defer func() { http.DefaultClient = originalClient }()

	req := mcp.CallToolRequest{Arguments: json.RawMessage(`{"query":"test"}`)}
	result, err := googleSearchHandler(context.Background(), req)
	
	if err != nil {
		t.Fatalf("handler returned an unexpected error: %v", err)
	}
	if result.Error == nil || !strings.Contains(result.Error.Message, "403 Forbidden") {
		t.Fatalf("expected API error in tool result, but got: %+v", result)
	}
}