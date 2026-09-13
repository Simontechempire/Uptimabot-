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
		responseTime := time.Since(start)

		status := "DOWN"

		if err == nil {
			resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				status = "UP"
			}
		}

		updateWebsiteStatus(
			target,
			status,
			responseTime.Milliseconds(),
		)

		if status == "UP" {
			fmt.Printf(
				"🟢 %s | UP | %dms\n",
				target,
				responseTime.Milliseconds(),
			)
		} else {
			fmt.Printf(
				"🔴 %s | DOWN\n",
				target,
			)
		}

		time.Sleep(1 * time.Minute)
	}
}
