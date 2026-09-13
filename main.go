package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
)

type PageData struct {
	Websites []Website
	Message  string
}

var pageTemplate = template.Must(template.New("dashboard").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Simon Tech Monitor</title>

	<style>
		* {
			box-sizing: border-box;
		}

		body {
			margin: 0;
			font-family: Arial, sans-serif;
			background: #0f172a;
			color: white;
		}

		header {
			background: #111827;
			padding: 35px 20px;
			text-align: center;
		}

		header h1 {
			margin: 0 0 10px;
			font-size: 32px;
		}

		header p {
			color: #94a3b8;
		}

		.container {
			max-width: 900px;
			margin: 30px auto;
			padding: 20px;
		}

		.panel {
			background: #1e293b;
			padding: 25px;
			border-radius: 16px;
			margin-bottom: 20px;
		}

		.form {
			display: flex;
			gap: 10px;
		}

		input {
			flex: 1;
			padding: 14px;
			border: 0;
			border-radius: 8px;
			font-size: 16px;
			outline: none;
		}

		button {
			padding: 14px 18px;
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

		.website {
			background: #0f172a;
			padding: 18px;
			border-radius: 12px;
			margin-top: 12px;
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 15px;
		}

		.url {
			word-break: break-all;
		}

		.status {
			font-weight: bold;
		}

		.up {
			color: #22c55e;
		}

		.down {
			color: #ef4444;
		}

		.checking {
			color: #facc15;
		}

		.remove {
			background: #dc2626;
		}

		.remove:hover {
			background: #b91c1c;
		}

		.empty {
			color: #94a3b8;
			text-align: center;
			padding: 20px;
		}

		footer {
			text-align: center;
			padding: 30px;
			color: #64748b;
		}

		@media (max-width: 650px) {
			.form {
				flex-direction: column;
			}

			.website {
				flex-direction: column;
				align-items: stretch;
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

	<div class="panel">
		<h2>➕ Add Website</h2>

		<form class="form" method="POST" action="/add">
			<input
				type="url"
				name="url"
				placeholder="https://example.com"
				required
			>

			<button type="submit">
				Add Website
			</button>
		</form>
	</div>

	<div class="panel">
		<h2>🌐 Monitored Websites</h2>

		{{if .Websites}}

			{{range .Websites}}
			<div class="website">

				<div class="url">
					<strong>{{.URL}}</strong>
				</div>

				<div class="status
					{{if eq .Status "UP"}}up
					{{else if eq .Status "DOWN"}}down
					{{else}}checking
					{{end}}">
					{{if eq .Status "UP"}}
						🟢 UP
					{{else if eq .Status "DOWN"}}
						🔴 DOWN
					{{else}}
						🟡 CHECKING
					{{end}}
				</div>

				<form method="POST" action="/remove">
					<input type="hidden" name="url" value="{{.URL}}">
					<button class="remove" type="submit">
						🗑️ Remove
					</button>
				</form>

			</div>
			{{end}}

		{{else}}
			<div class="empty">
				No websites are being monitored yet.
			</div>
		{{end}}

	</div>

</div>

<footer>
	<p>© 2026 Simon Tech Monitor</p>
	<p>Powered by Go 🚀</p>
</footer>

</body>
</html>
`))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Websites: getWebsites(),
	}

	err := pageTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "Unable to load dashboard", http.StatusInternalServerError)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	target := r.FormValue("url")

	parsed, err := url.ParseRequestURI(target)

	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		http.Error(w, "Invalid website URL", http.StatusBadRequest)
		return
	}

	addWebsite(target)

	go monitorWebsite(target)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	target := r.FormValue("url")

	removeWebsite(target)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {

	if err := initDatabase(); err != nil {
		fmt.Println("Database error:", err)
		return
	}

	if err := loadWebsites(); err != nil {
		fmt.Println("Failed to load websites:", err)
		return
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/remove", removeHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Simon Tech Monitor running on port %s\n", port)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
