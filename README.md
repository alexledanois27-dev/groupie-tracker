# Groupie Tracker

Groupie Tracker is a Go web server that displays artist data from the public
Groupie Trackers API. It listens on `127.0.0.1:8080` and uses the standard
library for HTTP handling, JSON decoding, templates, and graceful shutdown.

## Server overview

At startup, the server loads `templates/index.html`, creates a shared HTTP
client with a 10-second timeout, and registers the application routes. If the
template cannot be loaded, the application exits before opening the listening
socket.

The home page handler accepts `GET /`, retrieves the artist list through
`GetArtists`, and renders the HTML template with the decoded artist data. The
page is rendered into a buffer first so a template error can return a clean
`500 Internal Server Error` before any partial response is sent.

The `/api/*` routes act as a small proxy to the public API. They preserve query
parameters, forward the remote status code and response body, and reuse the
same timeout-enabled HTTP client.

| Route | Purpose |
| --- | --- |
| `GET /` | Render the artist list as HTML |
| `GET /static/*` | Serve files from the local `static` directory |
| `GET /api/artists` | Proxy the artists endpoint |
| `GET /api/locations` | Proxy the locations endpoint |
| `GET /api/dates` | Proxy the concert dates endpoint |
| `GET /api/relation` | Proxy the relations endpoint |

Unsupported methods return `405 Method Not Allowed`. Requests to the remote
API inherit the incoming request context, so they are cancelled if the client
disconnects.

The process listens for `SIGINT` and `SIGTERM`. When either signal arrives, it
stops accepting new connections and gives active requests up to five seconds
to complete before shutting down.

## Run locally

From the project root, run:

```sh
go run .
```

Then open `http://127.0.0.1:8080` in a browser.
