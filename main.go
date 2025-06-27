// google-search-agent – a minimal Go-based MCP agent for Google Search.
//
// Copyright 2025
// SPDX-License-Identifier: Apache-2.0
// Author: Ruslan Magana Vsevolodovna
//
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings" // Import the strings package
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	appName    = "google-search-agent"
	appVersion = "0.3.0" // Version bumped to reflect final fixes
)

var (
	logger    = log.New(os.Stderr, "", log.LstdFlags)
	startTime = time.Now()
)

// SearchResult represents a single Google search result item.
type SearchResult struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

// googleSearchHandler is the MCP tool handler for the "google_search" tool.
func googleSearchHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError("query parameter is required"), nil
	}

	// Read credentials from ENV and trim whitespace/control characters.
	// This is the fix for the "invalid control character" error.
	apiKey := strings.TrimSpace(os.Getenv("GOOGLE_API_KEY"))
	cseID := strings.TrimSpace(os.Getenv("GOOGLE_CSE_ID"))
	if apiKey == "" || cseID == "" {
		return mcp.NewToolResultError("environment variables GOOGLE_API_KEY and GOOGLE_CSE_ID must be set"), nil
	}

	logger.Printf("performing google search for query: %q", query)

	// Build the request URL for the Google Custom Search API.
	endpoint := fmt.Sprintf(
		"https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s",
		apiKey, cseID, url.QueryEscape(query),
	)

	// Make the HTTP request to the Google API.
	// Using a client with a timeout is better practice than http.Get.
	client := &http.Client{Timeout: 10 * time.Second}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logger.Printf("ERROR: failed to create google api request: %v", err)
		return nil, err
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		logger.Printf("ERROR: failed to call google api: %v", err)
		return nil, fmt.Errorf("failed to call Google Search API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Printf("ERROR: google api returned non-200 status: %s", resp.Status)
		return mcp.NewToolResultError(fmt.Sprintf("Google Search API returned an error: %s", resp.Status)), nil
	}

	// Decode the JSON response from the API.
	var apiResponse struct {
		Items []struct {
			Title string `json:"title"`
			Link  string `json:"link"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		logger.Printf("ERROR: failed to decode google api response: %v", err)
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Map the API results to our SearchResult type.
	results := make([]SearchResult, len(apiResponse.Items))
	for i, item := range apiResponse.Items {
		results[i] = SearchResult{Title: item.Title, Link: item.Link}
	}

	// Marshal the final results to a JSON string for the tool output.
	data, err := json.Marshal(results)
	if err != nil {
		logger.Printf("ERROR: failed to marshal final results: %v", err)
		return nil, fmt.Errorf("failed to marshal results: %w", err)
	}
	return mcp.NewToolResultText(string(data)), nil
}

func main() {
	// --- Flag Parsing ---
	// Implements flags mentioned in README, using reference server's pattern.
	listenHost := flag.String("listen", "0.0.0.0", "Listen interface for the HTTP server")
	port := flag.Int("port", 8080, "TCP port for the HTTP server")
	showHelp := flag.Bool("help", false, "Show help message")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s %s - Google Search MCP Agent\n\n", appName, appVersion)
		fmt.Fprintln(flag.CommandLine.Output(), "Usage:")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "\nEnvironment Variables:")
		fmt.Fprintln(flag.CommandLine.Output(), "  GOOGLE_API_KEY: Your Google API key for Custom Search API")
		fmt.Fprintln(flag.CommandLine.Output(), "  GOOGLE_CSE_ID:  Your Custom Search Engine ID (the 'cx' value)")
	}

	flag.Parse()
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// --- MCP Server Setup ---
	srv := server.NewMCPServer(
		appName,
		appVersion,
		server.WithToolCapabilities(false),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// --- Tool Registration ---
	// The tool name is "google_search" to match conventions and test scripts.
	tool := mcp.NewTool(
		"google_search",
		mcp.WithDescription("Performs a Google Custom Search query."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search terms to query."),
		),
	)
	srv.AddTool(tool, googleSearchHandler)

	// --- HTTP Server and Routing ---
	mux := http.NewServeMux()

	// Register health and version endpoints.
	registerHealthAndVersion(mux)

	// Mount the MCP HTTP handler at the root, wrapped in logging middleware.
	mcpHandler := server.NewStreamableHTTPServer(srv)
	mux.Handle("/", loggingHTTPMiddleware(mcpHandler))

	// --- Start Server ---
	addr := fmt.Sprintf("%s:%d", *listenHost, *port)
	logger.Printf("starting %s v%s on %s", appName, appVersion, addr)
	logger.Println("MCP endpoint available at / (POST with JSON-RPC)")
	logger.Println("Test with: curl -X POST -d '{\"jsonrpc\":\"2.0\",\"method\":\"tools/list\",\"id\":1}' http://"+addr)

	if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("server error: %v", err)
	}
}

// registerHealthAndVersion adds the /health and /version endpoints.
func registerHealthAndVersion(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Matching the documented output from the project's README.
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"name":"%s","version":"%s"}`, appName, appVersion)
	})
}

// loggingHTTPMiddleware provides request logging, inspired by the reference server.
func loggingHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sessionID := r.Header.Get("Mcp-Session-Id")

		// Use a response writer wrapper to capture status code.
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		duration := time.Since(start)

		if sessionID != "" {
			logger.Printf("%s %s %s [session: ...%s] %d %v", r.Method, r.URL.Path, r.RemoteAddr, shortID(sessionID), sw.status, duration)
		} else {
			logger.Printf("%s %s %s %d %v", r.Method, r.URL.Path, r.RemoteAddr, sw.status, duration)
		}
	})
}

// shortID returns the last 6 characters of a session ID for cleaner logging.
func shortID(id string) string {
	if len(id) > 6 {
		return id[len(id)-6:]
	}
	return id
}

// statusWriter wraps http.ResponseWriter to capture the status code for logging.
type statusWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func (sw *statusWriter) WriteHeader(code int) {
	if !sw.written {
		sw.status = code
		sw.written = true
		sw.ResponseWriter.WriteHeader(code)
	}
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	if !sw.written {
		sw.WriteHeader(http.StatusOK)
	}
	return sw.ResponseWriter.Write(b)
}