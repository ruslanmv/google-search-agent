package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	appName    = "google-search-agent"
	appVersion = "0.1.0"
)

// SearchResult represents a single Google result.
type SearchResult struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

// googleSearchHandler is the MCP tool handler for "google_search".
func googleSearchHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract the "query" argument
	query, err := req.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError("query parameter is required"), nil
	}

	// Read credentials from ENV
	apiKey := os.Getenv("GOOGLE_API_KEY")
	cseID := os.Getenv("GOOGLE_CSE_ID")
	if apiKey == "" || cseID == "" {
		return mcp.NewToolResultError("GOOGLE_API_KEY and GOOGLE_CSE_ID must be set"), nil
	}

	// Build the request URL
	escaped := url.QueryEscape(query)
	endpoint := fmt.Sprintf(
		"https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s",
		apiKey, cseID, escaped,
	)

	// Call the Google API
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response
	var out struct {
		Items []struct {
			Title string `json:"title"`
			Link  string `json:"link"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	// Map to our result type
	results := make([]SearchResult, len(out.Items))
	for i, item := range out.Items {
		results[i] = SearchResult{Title: item.Title, Link: item.Link}
	}

	// Marshal results to JSON and return as text
	data, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(data)), nil
}

func main() {
	// 1) Create the MCP server
	srv := server.NewMCPServer(
		appName,
		appVersion,
		server.WithToolCapabilities(false),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// 2) Register the google_search tool
	tool := mcp.NewTool(
		"google_search",
		mcp.WithDescription("Perform a Google Custom Search query"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search terms to query"),
		),
	)
	srv.AddTool(tool, googleSearchHandler)

	// 3) Expose health & version over plain HTTP
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	}))
	mux.Handle("/version", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"name":"%s","version":"%s"}`, appName, appVersion)
	}))

	// 4) Mount the MCP HTTP handler at root
	mux.Handle("/", server.NewStreamableHTTPServer(srv))

	// 5) Start listening
	addr := ":8080"
	fmt.Printf("Starting %s on %s\n", appName, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
