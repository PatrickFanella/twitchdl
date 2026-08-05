package media

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSanitizeFilename(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{input: " channel / title ", want: "channel - title"},
		{input: "...", want: "twitch"},
		{input: "one\n\ttwo", want: "one-two"},
		{input: strings.Repeat("x", 130), want: strings.Repeat("x", 120)},
	}
	for _, test := range tests {
		if got := SanitizeFilename(test.input); got != test.want {
			t.Errorf("SanitizeFilename(%q)=%q, want %q", test.input, got, test.want)
		}
	}
}

func TestEnsureOutputDirAndLivePaths(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "nested")
	got, err := EnsureOutputDir(dir)
	if err != nil {
		t.Fatalf("EnsureOutputDir: %v", err)
	}
	if _, err := os.Stat(got); err != nil {
		t.Fatalf("created directory missing: %v", err)
	}
	raw, staging, final := LivePaths(got, "test/channel", "mkv", time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC))
	if !strings.HasSuffix(raw, ".mkv.part") || !strings.HasSuffix(staging, ".mkv.finalizing") || !strings.HasSuffix(final, ".mkv") {
		t.Fatalf("unexpected paths: %q %q %q", raw, staging, final)
	}
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureOutputDir(file); err == nil {
		t.Fatal("file accepted as output directory")
	}
}
