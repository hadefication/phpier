package updater

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// ProgressReader wraps an io.Reader to track download progress
type ProgressReader struct {
	Reader   io.Reader
	Total    int64
	Current  int64
	Callback func(downloaded, total int64)
}

// Read implements io.Reader with progress tracking
func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Current += int64(n)

	if pr.Callback != nil {
		pr.Callback(pr.Current, pr.Total)
	}

	return n, err
}

// ProgressBar represents a simple console progress bar
type ProgressBar struct {
	total       int64
	current     int64
	width       int
	lastPrinted time.Time
	startTime   time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{
		total:     total,
		width:     50,
		startTime: time.Now(),
	}
}

// Update updates the progress bar with current progress
func (pb *ProgressBar) Update(current int64) {
	pb.current = current

	// Only update display every 100ms to avoid flickering
	now := time.Now()
	if now.Sub(pb.lastPrinted) < 100*time.Millisecond && current < pb.total {
		return
	}
	pb.lastPrinted = now

	pb.render()
}

// Finish completes the progress bar
func (pb *ProgressBar) Finish() {
	pb.current = pb.total
	pb.render()
	fmt.Println() // Add newline after completion
}

// render draws the progress bar to stdout
func (pb *ProgressBar) render() {
	percent := float64(pb.current) / float64(pb.total) * 100
	filled := int(float64(pb.width) * percent / 100)

	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.width-filled)

	// Calculate download speed and ETA
	elapsed := time.Since(pb.startTime)
	speed := float64(pb.current) / elapsed.Seconds()
	remaining := time.Duration(float64(pb.total-pb.current) / speed * float64(time.Second))

	// Format sizes
	currentStr := formatBytes(pb.current)
	totalStr := formatBytes(pb.total)
	speedStr := formatBytes(int64(speed)) + "/s"

	var eta string
	if pb.current < pb.total && speed > 0 {
		eta = fmt.Sprintf(" ETA: %s", formatDuration(remaining))
	}

	// Clear line and print progress
	fmt.Printf("\r[%s] %.1f%% (%s/%s) %s%s",
		bar, percent, currentStr, totalStr, speedStr, eta)
}

// formatBytes formats byte count as human-readable string
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// formatDuration formats duration as human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	} else {
		return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
	}
}
