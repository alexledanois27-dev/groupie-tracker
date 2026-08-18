# 🎵 Groupie Tracker

Groupie Tracker is a web application developed in Go that allows users to explore information about music artists and bands using the public **Groupie Trackers API**.

This project was developed as part of the Zone01 curriculum to practice:

- Building a web server in Go
- Working with HTML templates
- HTTP routing
- Consuming a REST API
- JSON parsing
- Structuring a web application

---

# ✨ Features

- Display all available artists
- Search artists by band name
- Search artists by member name
- Dedicated details page for each artist
- Display:
  - Creation year
  - First album release date
  - Band members
  - Concert locations
  - Concert dates
- HTTP error handling (404, 405, 500...)
- Graceful server shutdown

---

# 📂 Project Structure

```text
.
├── api.go
├── main.go
├── go.mod
├── README.md
├── static/
│   ├── style.css
│   └── ...
└── templates/
    ├── index.html
    └── details.html
```

---

# 🚀 Installation

Clone the repository:

```bash
git clone https://github.com/<organization>/groupie-tracker.git
cd groupie-tracker
```

Install the dependencies:

```bash
go mod tidy
```

Run the application:

```bash
go run .
```

The server will start on:

```
http://127.0.0.1:8080
```

---

# 🌐 Available Routes

| Route | Description |
|--------|-------------|
| `/` | Display the list of artists |
| `/artist/{id}` | Display detailed information about an artist |
| `/api/artists` | Proxy to the Artists API |
| `/api/locations` | Proxy to the Locations API |
| `/api/dates` | Proxy to the Dates API |
| `/api/relation` | Proxy to the Relations API |
| `/static/*` | Serve static assets |

---

# ⚙️ Technologies Used

- Go
- HTML5
- CSS3
- Go Templates
- `net/http`
- JSON
- Groupie Trackers API

---

# 📡 API

The application retrieves its data from the public **Groupie Trackers API**.

The API provides:

- Artists
- Band members
- Concert locations
- Concert dates
- Relationships between locations and dates

---

# 📚 Learning Objectives

This project demonstrates the use of:

- HTTP servers in Go
- HTML templating
- HTTP routing
- REST API consumption
- JSON decoding
- Error handling
- Clean project organization

---

# 👥 Authors

- Alexandre Ledanois
- Tom Marcel
- Nathalie Camargo

---

# 📄 License

This project was developed for educational purposes as part of the Zone01 curriculum.
