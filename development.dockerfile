FROM golang:1.23

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# Copy everything
COPY . .

# Build the wasm
RUN GOOS=js GOARCH=wasm go build -o assets/main.wasm ./cmd/wasm


EXPOSE 8080
CMD ["go", "run", "./cmd/web"]
