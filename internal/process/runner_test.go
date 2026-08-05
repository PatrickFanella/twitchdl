package process

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestOSRunner(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "success", args: []string{"-c", "printf hello"}, want: "hello"},
		{name: "failure", args: []string{"-c", "exit 7"}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var output bytes.Buffer
			err := (OSRunner{}).Run(context.Background(), Spec{Name: "sh", Args: test.args, Stdout: &output, Stderr: &output})
			if test.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("Run: %v", err)
			}
			if !strings.Contains(output.String(), test.want) {
				t.Fatalf("output=%q, want substring %q", output.String(), test.want)
			}
		})
	}
}

func TestOSRunnerCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err := (OSRunner{}).Run(ctx, Spec{Name: "sh", Args: []string{"-c", "while :; do :; done"}, GracePeriod: 50 * time.Millisecond})
	if err == nil || !strings.Contains(err.Error(), "cancel") {
		t.Fatalf("Run cancellation error=%v", err)
	}
}

func TestOSRunnerMissingCommand(t *testing.T) {
	t.Parallel()
	err := (OSRunner{}).Run(context.Background(), Spec{Name: "definitely-not-a-real-twitchdl-command"})
	if err == nil || !strings.Contains(err.Error(), "find") {
		t.Fatalf("missing command error=%v", err)
	}
}
