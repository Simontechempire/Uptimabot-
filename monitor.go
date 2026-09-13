func updateWebsiteStatus(target string, status string, responseTime int64) {
	mu.Lock()
	defer mu.Unlock()

	for i := range websites {
		if websites[i].URL == target {
			websites[i].Status = status
			websites[i].ResponseTime = responseTime
			break
		}
	}

	_ = updateDatabaseStatus(
		target,
		status,
		responseTime,
	)
}
