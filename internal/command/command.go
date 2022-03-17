package command

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	MaxPort = 65535
)

// CLIInput represents intput of the types expected on the command line.
type CLIInput struct {
	UpstreamUrls    string
	CheckInterval   int
	PrometheusPort  int
	UpstreamTimeout int
}

// ApplicationArgs represents input to the application that has been validated and parsed into types the application understands.
type ApplicationArgs struct {
	UpstreamUrls    []string
	CheckInterval   time.Duration
	PrometheusPort  int
	UpstreamTimeout time.Duration
}

// ParseInput parses the values passed on the CLI and modifies them to be formatted properly for use by the application
// If any of the arguments are invalid, it returns with an error.
func ParseInput(input CLIInput) (ApplicationArgs, error) {
	checkInterval, err := ParseInterval(input.CheckInterval)
	if err != nil {
		return ApplicationArgs{}, fmt.Errorf("invalid check-interval duration %d %w", input.CheckInterval, err)
	}

	upstreamUrls, err := ParseUrlsToCheck(input.UpstreamUrls)
	if err != nil {
		return ApplicationArgs{}, fmt.Errorf("invalid upstream-urls %s %w", input.UpstreamUrls, err)
	}

	timeout, err := ParseTimeout(input.UpstreamTimeout)
	if err != nil {
		return ApplicationArgs{}, fmt.Errorf("invalid timeout %d %w", input.UpstreamTimeout, err)
	}

	err = ParsePrometheusPort(input.PrometheusPort)
	if err != nil {
		return ApplicationArgs{}, fmt.Errorf("invalid port %d %w", input.PrometheusPort, err)
	}

	return ApplicationArgs{
		UpstreamUrls:    upstreamUrls,
		CheckInterval:   checkInterval,
		UpstreamTimeout: timeout,
		PrometheusPort:  input.PrometheusPort,
	}, nil
}

// ParsePrometheusPort parses the prometheus port and ensures its a valid value (0 < port < 65,535).
func ParsePrometheusPort(port int) error {
	if port < 1 {
		return fmt.Errorf("port must be > 0")
	}

	if port > MaxPort {
		return fmt.Errorf("port must be < %d", MaxPort)
	}

	return nil
}

// ParseTimeout parses the timeout values, ensures its valid (> 0).
func ParseTimeout(timeout int) (time.Duration, error) {
	if timeout < 1 {
		return 0, fmt.Errorf("timeout must be >= 1")
	}

	return time.Duration(timeout) * time.Second, nil
}

// ParseInterval parses the interval values, ensures its valid (> 0) and returns that as an interval in seconds.
func ParseInterval(interval int) (time.Duration, error) {
	if interval < 1 {
		return 0, fmt.Errorf("check interval must be >= 1")
	}

	return time.Duration(interval) * time.Second, nil
}

// ParseURLsToCheck parses a list of string to ensure each string is a valid url, it errors if any of the strings are not valid urls.
func ParseUrlsToCheck(urlsInput string) ([]string, error) {
	if urlsInput == "" {
		return []string{}, fmt.Errorf("no urls to parse")
	}

	unparsedURLs := strings.Split(urlsInput, ",")

	var urlsToCheck []string
	for _, unparsedURL := range unparsedURLs {
		_, err := url.Parse(unparsedURL)
		if err != nil {
			return []string{}, err
		}

		urlsToCheck = append(urlsToCheck, unparsedURL)
	}

	return urlsToCheck, nil
}
