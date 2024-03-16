package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func HttpGetWithRetry(url string, method string, retry uint) ([]byte, error) {
	retries := uint(0)
	maxRetries := retry
	backoff := 1 * time.Second

	for retries < maxRetries {
		fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			retries++
			fmt.Printf("Error making request: %v\n", err)
			fmt.Print("Retrying. . .")
			time.Sleep(backoff)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			retries++
			fmt.Printf("Error reading response body: %v\n", err)
			fmt.Print("Retrying. . .")
			time.Sleep(backoff)
			continue
		}

		return body, nil
	}

	fmt.Println("Failed to reach endpoint after", maxRetries, "retries")
	return nil, nil
}
