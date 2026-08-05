package media

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var unsafeFilename = regexp.MustCompile(`[[:cntrl:]/\\]+`)
var repeatedSpace = regexp.MustCompile(`\s+`)

// SanitizeFilename makes untrusted remote metadata safe as one path segment.
func SanitizeFilename(value string) string {
	value = unsafeFilename.ReplaceAllString(value, "-")
	value = repeatedSpace.ReplaceAllString(strings.TrimSpace(value), " ")
	value = strings.Trim(value, ".- ")
	if value == "" {
		return "twitch"
	}
	runes := []rune(value)
	if len(runes) > 120 {
		value = string(runes[:120])
	}
	return value
}

// EnsureOutputDir creates and validates an output directory.
func EnsureOutputDir(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve output directory: %w", err)
	}
	if err := os.MkdirAll(absolute, 0o755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return "", fmt.Errorf("inspect output directory: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("output path is not a directory: %s", absolute)
	}
	return absolute, nil
}

// LivePaths returns raw, staging, and final recording paths.
func LivePaths(outputDir, channel, extension string, now time.Time) (string, string, string) {
	name := fmt.Sprintf("%s - Live [%s]", SanitizeFilename(channel), now.Format("2006-01-02T15-04-05-0700"))
	final := filepath.Join(outputDir, name+"."+extension)
	return final + ".part", final + ".finalizing", final
}
