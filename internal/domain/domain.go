// Package domain defines UI-independent twitchdl requests, events, and results.
package domain

import (
	"errors"
	"fmt"
	"time"
)

// ErrInterrupted identifies a recording stopped by user cancellation.
var ErrInterrupted = errors.New("recording interrupted")

// Container is a final media container.
type Container string

const (
	// ContainerMKV is the resilient default for live recordings.
	ContainerMKV Container = "mkv"
	// ContainerMP4 is the compatibility-oriented alternative.
	ContainerMP4 Container = "mp4"
)

// EventKind identifies an application lifecycle event.
type EventKind string

const (
	EventResolving  EventKind = "resolving"
	EventWaiting    EventKind = "waiting"
	EventRecording  EventKind = "recording"
	EventFinalizing EventKind = "finalizing"
	EventValidating EventKind = "validating"
	EventComplete   EventKind = "complete"
	EventWarning    EventKind = "warning"
)

// Event is a UI-neutral application status update.
type Event struct {
	Kind    EventKind
	Message string
}

// Reporter receives application lifecycle events.
type Reporter interface {
	Report(Event)
}

// NopReporter discards application events.
type NopReporter struct{}

// Report implements Reporter.
func (NopReporter) Report(Event) {}

// DownloadRequest describes a clip or VOD download.
type DownloadRequest struct {
	URL       string
	Start     *time.Duration
	End       *time.Duration
	OutputDir string
	Force     bool
}

// LiveRequest describes a live recording.
type LiveRequest struct {
	Target      string
	OutputDir   string
	Quality     string
	Container   Container
	Duration    time.Duration
	Wait        bool
	WaitTimeout time.Duration
	FromStart   bool
	KeepRaw     bool
	Force       bool
}

// Result reports a completed or recoverable media operation.
type Result struct {
	Path        string
	RawPath     string
	Interrupted bool
}

// ValidateContainer validates a final container name.
func ValidateContainer(container Container) error {
	switch container {
	case ContainerMKV, ContainerMP4:
		return nil
	default:
		return fmt.Errorf("unsupported container %q: use mkv or mp4", container)
	}
}
