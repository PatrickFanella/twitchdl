package media

import "testing"

func TestTwitchInputs(t *testing.T) {
	t.Parallel()
	clipTests := []struct {
		value string
		ok    bool
	}{
		{value: "https://clips.twitch.tv/SomeClip", ok: true},
		{value: "https://www.twitch.tv/name/clip/SomeClip", ok: true},
		{value: "https://www.twitch.tv/videos/123", ok: false},
	}
	for _, test := range clipTests {
		err := ValidateClipURL(test.value)
		if (err == nil) != test.ok {
			t.Errorf("ValidateClipURL(%q) error=%v, want ok=%v", test.value, err, test.ok)
		}
	}

	if err := ValidateVODURL("https://www.twitch.tv/videos/123"); err != nil {
		t.Fatalf("valid VOD rejected: %v", err)
	}
	if err := ValidateVODURL("https://clips.twitch.tv/SomeClip"); err == nil {
		t.Fatal("clip accepted as VOD")
	}
}

func TestNormalizeLiveTarget(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input       string
		wantURL     string
		wantChannel string
		wantErr     bool
	}{
		{input: "TwitchDev", wantURL: "https://www.twitch.tv/TwitchDev", wantChannel: "twitchdev"},
		{input: "https://www.twitch.tv/TwitchDev", wantURL: "https://www.twitch.tv/TwitchDev", wantChannel: "twitchdev"},
		{input: "https://example.com/TwitchDev", wantErr: true},
		{input: "https://www.twitch.tv/videos/123", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			gotURL, gotChannel, err := NormalizeLiveTarget(test.input)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil || gotURL != test.wantURL || gotChannel != test.wantChannel {
				t.Fatalf("got (%q, %q, %v), want (%q, %q, nil)", gotURL, gotChannel, err, test.wantURL, test.wantChannel)
			}
		})
	}
}
