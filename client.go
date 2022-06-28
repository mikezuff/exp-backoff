package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
}

func DoRequest(ctx context.Context, url string) (string, error) {
	var (
		delay     time.Duration
		max       = time.Second * 5
		initRetry = 50 * time.Millisecond
		tries     = 1
		maxTries  = 10
	)

	for {
		resp, err := httpReq(ctx, url)
		if err == nil {
			return resp, nil
		}

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("getting %s: %w", url, err)
		default:
		}
		if tries >= maxTries {
			return "", fmt.Errorf("getting %s failed after %d tries with error: %w", url, tries, err)
		}
		if delay == 0 {
			delay = initRetry
		} else {
			delay *= 2
		}
		if delay > max {
			delay = max
		}
		log.Printf("request for %s failed: %s, retrying in %s", url, err, delay)
		time.Sleep(delay)

		tries++
	}
}

func httpReq(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return string(b), err
}
