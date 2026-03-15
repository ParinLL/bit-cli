package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const version = "1.0.0"

// --- Models ---

type LinkSummary struct {
	ID    int64  `json:"id"`
	Refer string `json:"refer"`
	Origin string `json:"origin"`
}

type Click struct {
	ID        int64   `json:"id"`
	UserAgent string  `json:"user_agent"`
	Country   *string `json:"country"`
	Browser   *string `json:"browser"`
	OS        *string `json:"os"`
	Referer   *string `json:"referer"`
	CreatedAt string  `json:"created_at"`
}

type Link struct {
	LinkSummary
	Clicks []Click `json:"clicks"`
}

type Pagination struct {
	HasMore bool   `json:"has_more"`
	Next    *int64 `json:"next"`
}

type APIError struct {
	Error string `json:"error"`
}

// --- Client ---

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient() *Client {
	baseURL := os.Getenv("BIT_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:4000"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	apiKey := os.Getenv("BIT_API_KEY")

	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.APIKey != "" {
		req.Header.Set("X-Api-Key", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return data, resp.StatusCode, nil
}

// --- API Methods ---

func (c *Client) Ping() error {
	data, status, err := c.do("GET", "/api/ping", nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("unexpected status %d: %s", status, string(data))
	}
	var resp struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	fmt.Println(resp.Data)
	return nil
}

func (c *Client) ListLinks(limit int, cursor string) error {
	path := fmt.Sprintf("/api/links?limit=%d", limit)
	if cursor != "" {
		path += "&cursor=" + cursor
	}

	data, status, err := c.do("GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return parseError(data, status)
	}

	var resp struct {
		Data       []LinkSummary `json:"data"`
		Pagination Pagination    `json:"pagination"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSHORT URL\tORIGINAL URL")
	fmt.Fprintln(w, "--\t---------\t------------")
	for _, l := range resp.Data {
		fmt.Fprintf(w, "%d\t%s\t%s\n", l.ID, l.Refer, l.Origin)
	}
	w.Flush()

	if resp.Pagination.HasMore && resp.Pagination.Next != nil {
		fmt.Printf("\n(more results available, use --cursor %d)\n", *resp.Pagination.Next)
	}
	return nil
}

func (c *Client) CreateLink(url string) error {
	body := map[string]string{"url": url}
	data, status, err := c.do("POST", "/api/links", body)
	if err != nil {
		return err
	}
	if status != 201 {
		return parseError(data, status)
	}

	var resp struct {
		Data Link `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	fmt.Printf("Created link #%d\n", resp.Data.ID)
	fmt.Printf("  Short: %s\n", resp.Data.Refer)
	fmt.Printf("  Origin: %s\n", resp.Data.Origin)
	return nil
}

func (c *Client) GetLink(id int64) error {
	data, status, err := c.do("GET", fmt.Sprintf("/api/links/%d", id), nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return parseError(data, status)
	}

	var resp struct {
		Data Link `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	fmt.Printf("Link #%d\n", resp.Data.ID)
	fmt.Printf("  Short:  %s\n", resp.Data.Refer)
	fmt.Printf("  Origin: %s\n", resp.Data.Origin)
	if len(resp.Data.Clicks) > 0 {
		fmt.Printf("  Recent clicks: %d\n", len(resp.Data.Clicks))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  ID\tCOUNTRY\tBROWSER\tOS\tREFERER\tTIME")
		for _, cl := range resp.Data.Clicks {
			fmt.Fprintf(w, "  %d\t%s\t%s\t%s\t%s\t%s\n",
				cl.ID, deref(cl.Country), deref(cl.Browser), deref(cl.OS), deref(cl.Referer), cl.CreatedAt)
		}
		w.Flush()
	}
	return nil
}

func (c *Client) UpdateLink(id int64, url string) error {
	body := map[string]string{"url": url}
	data, status, err := c.do("PUT", fmt.Sprintf("/api/links/%d", id), body)
	if err != nil {
		return err
	}
	if status != 200 {
		return parseError(data, status)
	}

	var resp struct {
		Data Link `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	fmt.Printf("Updated link #%d\n", resp.Data.ID)
	fmt.Printf("  Short:  %s\n", resp.Data.Refer)
	fmt.Printf("  Origin: %s\n", resp.Data.Origin)
	return nil
}

func (c *Client) DeleteLink(id int64) error {
	data, status, err := c.do("DELETE", fmt.Sprintf("/api/links/%d", id), nil)
	if err != nil {
		return err
	}
	if status != 204 {
		return parseError(data, status)
	}
	fmt.Printf("Deleted link #%d\n", id)
	return nil
}

func (c *Client) ListClicks(linkID int64, limit int, cursor string) error {
	path := fmt.Sprintf("/api/links/%d/clicks?limit=%d", linkID, limit)
	if cursor != "" {
		path += "&cursor=" + cursor
	}

	data, status, err := c.do("GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return parseError(data, status)
	}

	var resp struct {
		Data       []Click    `json:"data"`
		Pagination Pagination `json:"pagination"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tCOUNTRY\tBROWSER\tOS\tREFERER\tTIME")
	fmt.Fprintln(w, "--\t-------\t-------\t--\t-------\t----")
	for _, cl := range resp.Data {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			cl.ID, deref(cl.Country), deref(cl.Browser), deref(cl.OS), deref(cl.Referer), cl.CreatedAt)
	}
	w.Flush()

	if resp.Pagination.HasMore && resp.Pagination.Next != nil {
		fmt.Printf("\n(more results available, use --cursor %d)\n", *resp.Pagination.Next)
	}
	return nil
}

// --- Helpers ---

func deref(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}

func parseError(data []byte, status int) error {
	var apiErr APIError
	if err := json.Unmarshal(data, &apiErr); err == nil && apiErr.Error != "" {
		return fmt.Errorf("API error (%d): %s", status, apiErr.Error)
	}
	return fmt.Errorf("unexpected status %d: %s", status, string(data))
}

func mustParseInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid ID %q\n", s)
		os.Exit(1)
	}
	return n
}

// --- CLI ---

func usage() {
	fmt.Printf(`bit-cli v%s - CLI for Bit URL Shortener

Usage:
  bit <command> [arguments]

Commands:
  ping                          Health check
  list   [--limit N] [--cursor C]  List all links
  create <url>                  Create a short link
  get    <id>                   Get link details + recent clicks
  update <id> <url>             Update link URL
  delete <id>                   Delete a link
  clicks <id> [--limit N] [--cursor C]  List clicks for a link

Environment:
  BIT_API_URL   Base URL (default: http://localhost:4000)
  BIT_API_KEY   API key for authentication
`, version)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	client := NewClient()
	cmd := os.Args[1]
	args := os.Args[2:]

	var err error

	switch cmd {
	case "ping":
		err = client.Ping()

	case "list":
		limit := 100
		cursor := ""
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "--limit":
				i++
				if i < len(args) {
					limit, _ = strconv.Atoi(args[i])
				}
			case "--cursor":
				i++
				if i < len(args) {
					cursor = args[i]
				}
			}
		}
		err = client.ListLinks(limit, cursor)

	case "create":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: bit create <url>")
			os.Exit(1)
		}
		err = client.CreateLink(args[0])

	case "get":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: bit get <id>")
			os.Exit(1)
		}
		err = client.GetLink(mustParseInt64(args[0]))

	case "update":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: bit update <id> <url>")
			os.Exit(1)
		}
		err = client.UpdateLink(mustParseInt64(args[0]), args[1])

	case "delete":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: bit delete <id>")
			os.Exit(1)
		}
		err = client.DeleteLink(mustParseInt64(args[0]))

	case "clicks":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: bit clicks <id> [--limit N] [--cursor C]")
			os.Exit(1)
		}
		linkID := mustParseInt64(args[0])
		limit := 100
		cursor := ""
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--limit":
				i++
				if i < len(args) {
					limit, _ = strconv.Atoi(args[i])
				}
			case "--cursor":
				i++
				if i < len(args) {
					cursor = args[i]
				}
			}
		}
		err = client.ListClicks(linkID, limit, cursor)

	case "version", "--version", "-v":
		fmt.Printf("bit-cli v%s\n", version)

	case "help", "--help", "-h":
		usage()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
