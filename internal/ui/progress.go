// Package ui provides a simple interactive terminal progress display for
// long-running operations. It shows a spinner, elapsed time, and step label
// without requiring external dependencies.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Step represents a single labelled step in the progress display.
type Step struct {
	Label   string
	Percent int // 0–100, used for a progress bar
}

// Progress manages a live terminal progress display.
type Progress struct {
	mu      sync.Mutex
	w       io.Writer
	current string
	percent int
	start   time.Time
	ticker  *time.Ticker
	done    chan struct{}
	stopped bool
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// New creates a new Progress that writes to stdout.
func New() *Progress {
	return &Progress{w: os.Stdout}
}

// Start begins displaying the progress indicator with the given label.
func (p *Progress) Start(label string) {
	p.mu.Lock()
	p.current = label
	p.percent = 0
	p.start = time.Now()
	p.done = make(chan struct{})
	p.stopped = false
	p.mu.Unlock()

	p.ticker = time.NewTicker(100 * time.Millisecond)
	go p.loop()
}

// Update changes the label and percentage shown.
func (p *Progress) Update(label string, percent int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = label
	p.percent = percent
}

// Stop ends the spinner and prints a final line.
func (p *Progress) Stop(finalMsg string, success bool) {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}
	p.stopped = true
	p.mu.Unlock()

	p.ticker.Stop()
	close(p.done)

	// Clear the spinner line
	fmt.Fprintf(p.w, "\r%s\r", strings.Repeat(" ", 80))

	icon := "✅"
	if !success {
		icon = "❌"
	}
	elapsed := time.Since(p.start).Round(time.Millisecond)
	fmt.Fprintf(p.w, " %s %s (%s)\n", icon, finalMsg, elapsed)
}

// loop refreshes the terminal line every tick.
func (p *Progress) loop() {
	frame := 0
	for {
		select {
		case <-p.done:
			return
		case <-p.ticker.C:
			p.mu.Lock()
			label := p.current
			pct := p.percent
			elapsed := time.Since(p.start).Round(time.Second)
			p.mu.Unlock()

			spin := spinnerFrames[frame%len(spinnerFrames)]
			frame++

			bar := progressBar(pct, 20)
			// \r returns to line start so we overwrite in place
			fmt.Fprintf(p.w, "\r %s [%s] %3d%%  %s  (%s)    ",
				spin, bar, pct, label, elapsed)
		}
	}
}

// progressBar returns an ASCII progress bar of `width` chars for pct 0–100.
func progressBar(pct, width int) string {
	filled := pct * width / 100
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// ── Simple one-liner helpers ──────────────────────────────────────────────────

// Infof prints an informational line (not part of the spinner).
func Infof(format string, a ...any) {
	fmt.Printf(" ℹ️  "+format+"\n", a...)
}

// Warnf prints a warning line.
func Warnf(format string, a ...any) {
	fmt.Printf(" ⚠️  "+format+"\n", a...)
}

// Errorf prints an error line.
func Errorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, " ❌ "+format+"\n", a...)
}

// Successf prints a success line.
func Successf(format string, a ...any) {
	fmt.Printf(" ✅ "+format+"\n", a...)
}
