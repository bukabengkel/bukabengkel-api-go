package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func HttpGetWithRetry(url string, maxRetries uint) ([]byte, error) {
	retries := uint(0)
	backoff := 1 * time.Second

	for retries < maxRetries {
		resp, err := http.Get(url)
		if err != nil {
			retries++
			fmt.Printf("Error making request: %v\n", err)
			fmt.Print("Retrying. . .\n")
			time.Sleep(backoff)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			retries++
			fmt.Printf("Error reading response body: %v\n", err)
			fmt.Print("Retrying. . .\n")
			time.Sleep(backoff)
			continue
		}

		return body, nil
	}

	fmt.Println("Failed to reach endpoint after", maxRetries, "retries")
	return nil, nil
}

func HttpPostWithRetry(url string, body []byte, maxRetries uint) ([]byte, error) {
	retries := uint(0)
	backoff := 1 * time.Second

	for retries < maxRetries {
		fmt.Println("Making request to", url)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			retries++
			fmt.Printf("Error making request: %v\n", err)
			fmt.Print("Retrying. . .\n")
			time.Sleep(backoff)
			continue
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			retries++
			fmt.Printf("Error reading response body: %v\n", err)
			fmt.Print("Retrying. . .\n")
			time.Sleep(backoff)
			continue
		}

		return body, nil
	}

	fmt.Println("Failed to reach endpoint after", maxRetries, "retries")
	return nil, errors.New("failed to reach endpoint after " + strconv.Itoa(int(maxRetries)) + " retries")
}
