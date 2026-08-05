package media

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var channelPattern = regexp.MustCompile(`^[A-Za-z0-9_]{1,25}$`)

// ValidateClipURL validates a public Twitch clip URL.
func ValidateClipURL(raw string) error {
	host, path, err := twitchURL(raw)
	if err != nil {
		return err
	}
	if host == "clips.twitch.tv" && strings.Trim(path, "/") != "" {
		return nil
	}
	if (host == "twitch.tv" || host == "www.twitch.tv") && strings.Contains(path, "/clip/") {
		return nil
	}
	return fmt.Errorf("URL is not a Twitch clip URL")
}

// ValidateVODURL validates a Twitch VOD URL.
func ValidateVODURL(raw string) error {
	host, path, err := twitchURL(raw)
	if err != nil {
		return err
	}
	if (host == "twitch.tv" || host == "www.twitch.tv") && strings.HasPrefix(path, "/videos/") {
		return nil
	}
	return fmt.Errorf("URL is not a Twitch VOD URL")
}

// NormalizeLiveTarget converts a channel name or Twitch channel URL to a URL.
func NormalizeLiveTarget(raw string) (string, string, error) {
	if channelPattern.MatchString(raw) {
		return "https://www.twitch.tv/" + raw, strings.ToLower(raw), nil
	}
	host, path, err := twitchURL(raw)
	if err != nil {
		return "", "", fmt.Errorf("invalid Twitch channel: %w", err)
	}
	if host != "twitch.tv" && host != "www.twitch.tv" {
		return "", "", fmt.Errorf("live URL must use twitch.tv")
	}
	channel := strings.Trim(path, "/")
	if strings.Contains(channel, "/") || !channelPattern.MatchString(channel) {
		return "", "", fmt.Errorf("URL is not a Twitch channel URL")
	}
	return "https://www.twitch.tv/" + channel, strings.ToLower(channel), nil
}

func twitchURL(raw string) (string, string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", "", fmt.Errorf("parse URL: %w", err)
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return "", "", fmt.Errorf("URL must start with http:// or https://")
	}
	return strings.ToLower(parsed.Hostname()), parsed.EscapedPath(), nil
}
