package app

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// ToolStatus describes one external media dependency.
type ToolStatus struct {
	Name    string
	Path    string
	Version string
	Error   error
}

// Doctor checks required external media tools without reading credentials.
func Doctor(ctx context.Context) []ToolStatus {
	tools := []struct {
		name string
		args []string
	}{
		{name: "yt-dlp", args: []string{"--version"}},
		{name: "streamlink", args: []string{"--version"}},
		{name: "ffmpeg", args: []string{"-version"}},
		{name: "ffprobe", args: []string{"-version"}},
	}
	statuses := make([]ToolStatus, 0, len(tools))
	for _, tool := range tools {
		status := ToolStatus{Name: tool.name}
		path, err := exec.LookPath(tool.name)
		if err != nil {
			status.Error = fmt.Errorf("not installed: %w", err)
			statuses = append(statuses, status)
			continue
		}
		status.Path = path
		output, err := exec.CommandContext(ctx, path, tool.args...).CombinedOutput()
		if err != nil {
			status.Error = fmt.Errorf("read version: %w", err)
		} else {
			status.Version = strings.TrimSpace(strings.SplitN(string(output), "\n", 2)[0])
		}
		statuses = append(statuses, status)
	}
	return statuses
}
