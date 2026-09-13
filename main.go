package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type CheckResult struct {
	URL        string
	Status     string
	StatusCode int
	Response   time.Duration
	Error      string
}

func checkWebsite(target string) CheckResult {
	start := time.Now()

	parsed, err := url.ParseRequestURI(target)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return CheckResult{
			URL:    target,
			Status: "INVALID",
			Error:  "Please enter a valid URL.",
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(target)
	responseTime := time.Since(start)

	if err != nil {
		return CheckResult{
			URL:      target,
			Status:   "DOWN",
			Response: responseTime,
			Error:    err.Error(),
		}
	}

	defer resp.Body.Close()

	status := "UP"

	if resp.StatusCode >= 400 {
		status = "DOWN"
	}

	return CheckResult{
		URL:        target,
		Status:     status,
		StatusCode: resp.StatusCode,
		Response:   responseTime,
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("url")

	resultHTML := ""

	if target != "" {
		result := checkWebsite(target)

		if result.Status == "UP" {
			resultHTML = fmt.Sprintf(`
				<div class="result up">
					<h2>🟢 WEBSITE ONLINE</h2>
					<p><strong>Website:</strong> %s</p>
					<p><strong>Status:</strong> %d</p>
					<p><strong>Response:</strong> %v</p>
				</div>
			`, result.URL, result.StatusCode, result.Response)
		} else {
			resultHTML = fmt.Sprintf(`
				<div class="result down">
					<h2>🔴 WEBSITE DOWN</h2>
					<p><strong>Website:</strong> %s</p>
					<p>%s</p>
				</div>
			`, result.URL, result.Error)
		}
	}

	fmt.Fprintln(w, `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Simon Tech Monitor</title>

	<style>
		body {
			margin: 0;
			font-family: Arial, sans-serif;
			background: #0f172a;
			color: white;
		}

		header {
			text-align: center;
			padding: 40px 20px;
			background: #111827;
		}

		header h1 {
			font-size: 32px;
		}

		header p {
			color: #94a3b8;
		}

		.container {
			max-width: 800px;
			margin: 40px auto;
			padding: 20px;
		}

		.monitor {
			background: #1e293b;
			padding: 30px;
			border-radius: 16px;
			text-align: center;
		}

		input {
			width: 65%;
			padding: 15px;
			border: 0;
			border-radius: 8px;
			font-size: 16px;
		}

		button {
			padding: 15px 20px;
			margin-left: 8px;
			border: 0;
			border-radius: 8px;
			background: #2563eb;
			color: white;
			font-weight: bold;
			cursor: pointer;
		}

		button:hover {
			background: #1d4ed8;
		}

		.result {
			margin-top: 25px;
			padding: 20px;
			border-radius: 12px;
			text-align: left;
		}

		.up {
			background: #064e3b;
			border: 1px solid #10b981;
		}

		.down {
			background: #450a0a;
			border: 1px solid #ef4444;
		}

		footer {
			text-align: center;
			padding: 30px;
			color: #64748b;
		}

		@media (max-width: 600px) {
			input {
				width: 100%;
				margin-bottom: 10px;
			}

			button {
				width: 100%;
				margin: 0;
			}
		}
	</style>
</head>

<body>

<header>
	<h1>🚀 Simon Tech Monitor</h1>
	<p>Website Monitoring Dashboard</p>
</header>

<div class="container">
	<div class="monitor">

		<h2>🔍 Check a Website</h2>

		<form method="GET">
			<input
				type="url"
				name="url"
				placeholder="https://example.com"
				required
			>

			<button type="submit">
				Check Website
			</button>
		</form>
`)

	fmt.Fprintln(w, resultHTML)

	fmt.Fprintln(w, `
	</div>
</div>

<footer>
	<p>© 2026 Simon Tech Monitor</p>
	<p>Powered by Go 🚀</p>
</footer>

</body>
</html>
`)
}

func main() {

	http.HandleFunc("/", homeHandler)

	// Start automatic monitoring in the background.
	go monitorWebsite("https://example.com")

	fmt.Println("🚀 Simon Tech Monitor running on port 8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
