package media

import (
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{name: "seconds", input: "90", want: 90 * time.Second},
		{name: "minutes", input: "1:30", want: 90 * time.Second},
		{name: "hours", input: "01:02:03", want: time.Hour + 2*time.Minute + 3*time.Second},
		{name: "milliseconds", input: "1:30.500", want: 90*time.Second + 500*time.Millisecond},
		{name: "bad seconds", input: "1:70", wantErr: true},
		{name: "bad minutes", input: "1:70:00", wantErr: true},
		{name: "negative", input: "-1", wantErr: true},
		{name: "text", input: "later", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseTimestamp(test.input)
			if test.wantErr {
				if err == nil {
					t.Fatalf("ParseTimestamp(%q) succeeded, want error", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseTimestamp(%q): %v", test.input, err)
			}
			if got != test.want {
				t.Fatalf("ParseTimestamp(%q) = %s, want %s", test.input, got, test.want)
			}
		})
	}
}

func TestFormatTimestamp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input time.Duration
		want  string
	}{
		{input: 90 * time.Second, want: "00:01:30"},
		{input: time.Hour + 2*time.Minute + 3*time.Second, want: "01:02:03"},
		{input: 1500 * time.Millisecond, want: "00:00:01.500"},
	}
	for _, test := range tests {
		if got := FormatTimestamp(test.input); got != test.want {
			t.Errorf("FormatTimestamp(%s) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestValidateRange(t *testing.T) {
	t.Parallel()
	start := 2 * time.Minute
	earlier := time.Minute
	later := 3 * time.Minute
	if err := ValidateRange(&start, &later); err != nil {
		t.Fatalf("valid range rejected: %v", err)
	}
	if err := ValidateRange(&start, &earlier); err == nil {
		t.Fatal("reversed range accepted")
	}
	if err := ValidateRange(nil, &later); err == nil {
		t.Fatal("end-only range accepted")
	}
}
