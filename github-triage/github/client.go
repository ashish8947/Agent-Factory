package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github-triage/models"
)

type Client struct {
	token string
	owner string
	repo  string
}

func NewClient(token, owner, repo string) *Client {
	return &Client{
		token: token,
		owner: owner,
		repo:  repo,
	}
}

func (c *Client) makeRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println("[GitHub][ERROR] Failed to create request:", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 20 * time.Second}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[GitHub][ERROR] Request failed:", err)
		return nil, err
	}

	log.Printf("[GitHub] %s %s → status=%d duration=%s\n",
		method, url, resp.StatusCode, time.Since(start))

	return resp, nil
}

func (c *Client) FetchOpenIssues() ([]models.Issue, error) {
	log.Println("[GitHub] Fetching open issues")

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=open", c.owner, c.repo)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issues []models.Issue
	err = json.NewDecoder(resp.Body).Decode(&issues)
	if err != nil {
		log.Println("[GitHub][ERROR] Failed to decode issues:", err)
		return nil, err
	}

	log.Printf("[GitHub] Fetched %d issues\n", len(issues))

	return issues, nil
}

func (c *Client) AddLabel(issueNumber int, label string) error {
	log.Printf("[GitHub] Adding label '%s' to issue #%d\n", label, issueNumber)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/labels", c.owner, c.repo, issueNumber)

	payload := []string{label}
	body, _ := json.Marshal(payload)

	resp, err := c.makeRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) Comment(issueNumber int, comment string) error {
	log.Printf("[GitHub] Adding comment to issue #%d\n", issueNumber)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments", c.owner, c.repo, issueNumber)

	payload := map[string]string{
		"body": comment,
	}
	body, _ := json.Marshal(payload)

	resp, err := c.makeRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
func (c *Client) FetchRecentIssues(limit int) ([]models.Issue, error) {

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/issues?state=open&per_page=%d",
		c.owner, c.repo, limit,
	)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issues []models.Issue
	err = json.NewDecoder(resp.Body).Decode(&issues)
	if err != nil {
		return nil, err
	}

	return issues, nil
}