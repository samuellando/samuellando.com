FROM golang:1.23

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# Copy everything
COPY . .

# Generate sqlc
COPY --from=sqlc/sqlc:1.28.0 /workspace/sqlc /usr/bin/sqlc
RUN sqlc generate
# Build the wasm
RUN GOOS=js GOARCH=wasm go build -o static/main.wasm ./wasm

EXPOSE 8080
CMD ["go", "run", "./cmd/web"]
