package checkurl

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zrobisho/url-checker/internal/metrics"
)

type URLChecker struct {
	Client *http.Client
}

// CheckURL checks a given URL to see determine if it returned a 200 status code, if any other status code is returned
// or if the call call to the service failed, an error is returned.
func (u URLChecker) CheckURL(url string) error {
	start := time.Now()

	resp, err := u.Client.Get(url)
	if err != nil {
		return err
	}

	metrics.ExternalURLResponseTime(url, start)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("url %s returned a non 200 status", url)
	}

	return nil
}
