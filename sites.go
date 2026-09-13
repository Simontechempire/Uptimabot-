package main

import (
	"sync"
)

type Website struct {
	URL    string
	Status string
}

var (
	websites = []Website{}
	mu       sync.Mutex
)

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
}

func removeWebsite(target string) {
	mu.Lock()
	defer mu.Unlock()

	for i, site := range websites {
		if site.URL == target {
			websites = append(websites[:i], websites[i+1:]...)
			return
		}
	}
}

func getWebsites() []Website {
	mu.Lock()
	defer mu.Unlock()

	copyList := make([]Website, len(websites))
	copy(copyList, websites)

	return copyList
}
