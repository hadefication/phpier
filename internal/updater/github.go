package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	GitHubAPIURL = "https://api.github.com"
	RepoOwner    = "phpier" // TODO: Update with actual GitHub org/user
	RepoName     = "phpier"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName     string        `json:"tag_name"`
	Name        string        `json:"name"`
	Body        string        `json:"body"`
	Draft       bool          `json:"draft"`
	Prerelease  bool          `json:"prerelease"`
	CreatedAt   time.Time     `json:"created_at"`
	PublishedAt time.Time     `json:"published_at"`
	Assets      []GitHubAsset `json:"assets"`
}

// GitHubAsset represents a release asset
type GitHubAsset struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Label              string `json:"label"`
	ContentType        string `json:"content_type"`
	Size               int64  `json:"size"`
	DownloadCount      int    `json:"download_count"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// GitHubClient handles GitHub API interactions
type GitHubClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: GitHubAPIURL,
	}
}

// GetLatestRelease fetches the latest release from GitHub
func (c *GitHubClient) GetLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", c.baseURL, RepoOwner, RepoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add User-Agent header as required by GitHub API
	req.Header.Set("User-Agent", "phpier-cli")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &release, nil
}

// GetRelease fetches a specific release by tag
func (c *GitHubClient) GetRelease(tag string) (*GitHubRelease, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", c.baseURL, RepoOwner, RepoName, tag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "phpier-cli")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("release %s not found", tag)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &release, nil
}

// FindAssetForPlatform finds the appropriate asset for the current platform
func (r *GitHubRelease) FindAssetForPlatform(os, arch string) (*GitHubAsset, error) {
	expectedName := fmt.Sprintf("phpier-%s-%s", os, arch)

	for _, asset := range r.Assets {
		if asset.Name == expectedName {
			return &asset, nil
		}

		// Also check for common variations
		if strings.Contains(asset.Name, os) && strings.Contains(asset.Name, arch) {
			return &asset, nil
		}
	}

	return nil, fmt.Errorf("no asset found for platform %s/%s", os, arch)
}

// DownloadAsset downloads a release asset with progress tracking
func (c *GitHubClient) DownloadAsset(asset *GitHubAsset, progressCallback func(downloaded, total int64)) ([]byte, error) {
	req, err := http.NewRequest("GET", asset.BrowserDownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	req.Header.Set("User-Agent", "phpier-cli")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Use progress reader if callback provided
	var reader io.Reader = resp.Body
	if progressCallback != nil {
		reader = &ProgressReader{
			Reader:   resp.Body,
			Total:    asset.Size,
			Callback: progressCallback,
		}
	}

	return io.ReadAll(reader)
}
