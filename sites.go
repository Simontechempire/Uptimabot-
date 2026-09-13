package main

import "sync"

type Website struct {
	URL          string
	Status       string
	ResponseTime int64
}

var (
	websites = []Website{}
	mu       sync.Mutex
)

func loadWebsites() error {
	rows, err := db.Query(`
		SELECT url, status, response_time
		FROM websites
		ORDER BY id DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	mu.Lock()
	defer mu.Unlock()

	websites = []Website{}

	for rows.Next() {
		var site Website

		err := rows.Scan(
			&site.URL,
			&site.Status,
			&site.ResponseTime,
		)
		if err != nil {
			return err
		}

		websites = append(websites, site)
	}

	return rows.Err()
}

func addWebsite(target string) {
	mu.Lock()
	defer mu.Unlock()

	for _, site := range websites {
		if site.URL == target {
			return
		}
	}

	websites = append(websites, Website{
		URL:    target,
		Status: "CHECKING",
	})

	_ = saveWebsite(target)
}

func removeWebsite(target string) {
	mu.Lock()
	defer mu.Unlock()

	for i, site := range websites {
		if site.URL == target {
			websites = append(websites[:i], websites[i+1:]...)
			break
		}
	}

	_ = deleteWebsite(target)
}

func updateWebsiteStatus(
	target string,
	status string,
	responseTime int64,
) {
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

func getWebsites() []Website {
	mu.Lock()
	defer mu.Unlock()

	copyList := make([]Website, len(websites))
	copy(copyList, websites)

	return copyList
}
