package main

import (
	"fmt"
	"net/http"
	"time"
)

func monitorWebsite(target string) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		start := time.Now()

		resp, err := client.Get(target)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("🔴 %s is DOWN | Error: %v\n", target, err)
		} else {
			resp.Body.Close()

			if resp.StatusCode >= 400 {
				fmt.Printf("🔴 %s is DOWN | Status: %d | Response: %v\n",
					target, resp.StatusCode, elapsed)
			} else {
				fmt.Printf("🟢 %s is UP | Status: %d | Response: %v\n",
					target, resp.StatusCode, elapsed)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
