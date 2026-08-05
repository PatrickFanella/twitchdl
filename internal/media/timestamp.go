// Package media implements media input and output normalization.
package media

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var timestampPattern = regexp.MustCompile(`^\d+(?::\d{1,2}){0,2}(?:\.\d{1,3})?$`)

// ParseTimestamp parses seconds, MM:SS, or HH:MM:SS with optional milliseconds.
func ParseTimestamp(value string) (time.Duration, error) {
	if !timestampPattern.MatchString(value) {
		return 0, fmt.Errorf("invalid timestamp %q: use seconds, MM:SS, or HH:MM:SS", value)
	}

	parts := strings.Split(value, ":")
	secondsPart := parts[len(parts)-1]
	seconds, err := strconv.ParseFloat(secondsPart, 64)
	if err != nil {
		return 0, fmt.Errorf("parse timestamp seconds: %w", err)
	}
	if len(parts) > 1 && seconds >= 60 {
		return 0, fmt.Errorf("invalid timestamp %q: seconds must be below 60", value)
	}

	minutes := 0
	hours := 0
	if len(parts) >= 2 {
		minutes, err = strconv.Atoi(parts[len(parts)-2])
		if err != nil {
			return 0, fmt.Errorf("parse timestamp minutes: %w", err)
		}
	}
	if len(parts) == 3 {
		if minutes >= 60 {
			return 0, fmt.Errorf("invalid timestamp %q: minutes must be below 60", value)
		}
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("parse timestamp hours: %w", err)
		}
	}

	total := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds*float64(time.Second))
	return total, nil
}

// FormatTimestamp formats a duration for yt-dlp section ranges.
func FormatTimestamp(value time.Duration) string {
	if value < 0 {
		value = 0
	}
	totalMilliseconds := value.Milliseconds()
	hours := totalMilliseconds / 3_600_000
	minutes := (totalMilliseconds / 60_000) % 60
	seconds := (totalMilliseconds / 1_000) % 60
	milliseconds := totalMilliseconds % 1_000
	if milliseconds == 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

// ValidateRange validates an optional media range.
func ValidateRange(start, end *time.Duration) error {
	if end != nil && start == nil {
		return fmt.Errorf("end timestamp requires a start timestamp")
	}
	if start != nil && *start < 0 {
		return fmt.Errorf("start timestamp cannot be negative")
	}
	if end != nil && *end <= *start {
		return fmt.Errorf("end timestamp must be later than start timestamp")
	}
	return nil
}
