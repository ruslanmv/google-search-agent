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

// SearchRequest is the JSON payload for google_search
type SearchRequest struct {
    Query string `json:"query"`
}

// SearchResult is one result item from Google Custom Search
type SearchResult struct {
    Title string `json:"title"`
    Link  string `json:"link"`
}

// googleSearch queries the Google Custom Search API and returns a slice of results.
func googleSearch(ctx context.Context, payload json.RawMessage) (interface{}, error) {
    var req SearchRequest
    if err := json.Unmarshal(payload, &req); err != nil {
        return nil, err
    }

    apiKey := os.Getenv("GOOGLE_API_KEY")
    cseID  := os.Getenv("GOOGLE_CSE_ID")
    if apiKey == "" || cseID == "" {
        return nil, fmt.Errorf("GOOGLE_API_KEY and GOOGLE_CSE_ID must be set")
    }

    q := url.QueryEscape(req.Query)
    endpoint := fmt.Sprintf(
        "https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s",
        apiKey, cseID, q,
    )

    resp, err := http.Get(endpoint)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var out struct {
        Items []struct {
            Title string `json:"title"`
            Link  string `json:"link"`
        } `json:"items"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, err
    }

    results := make([]SearchResult, len(out.Items))
    for i, item := range out.Items {
        results[i] = SearchResult{Title: item.Title, Link: item.Link}
    }
    return results, nil
}

func main() {
    // Create an MCP server instance
    s := server.NewMCPServer(
        appName,
        appVersion,
        server.WithToolCapabilities(false),
        server.WithLogging(),
        server.WithRecovery(),
    )

    // Register health and version endpoints
    s.AddHTTPHandler("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, `{"status":"ok"}`)
    })
    s.AddHTTPHandler("/version", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, `{"name":"%s","version":"%s"}`, appName, appVersion)
    })

    // Define and register the google_search tool
    googleTool := mcp.NewTool("google_search",
        mcp.WithDescription("Perform a Google Custom Search query"),
        mcp.WithString("query",
            mcp.Required(),
            mcp.Description("Search terms to query"),
        ),
    )
    s.AddTool(googleTool, googleSearch)

    // Start the HTTP server on port 8080
    addr := ":8080"
    fmt.Printf("Starting %s on %s\n", appName, addr)
    if err := http.ListenAndServe(addr, server.WithHandler(s)); err != nil {
        fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
        os.Exit(1)
    }
}
