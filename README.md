# Order APp

In memory orders app written in Golang and utilizing Go's concurrency features for concurrent processing. It uses:

1. Go channels for async communication and processing.
2. Sync.map as a safe in-memory database to prevent race conditions that typically occurs with standard maps on read/write attempts.

# Running the application

1. Clone the repository.
2. Run `go get ./...`
3. Run `go run server.go`

The server listnes for requests on port 3000

# Running unit tests

Run `go test ./...`