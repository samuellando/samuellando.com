FROM golang:1.23 as build

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# Copy everything
COPY . .

EXPOSE 8080
CMD ["go", "run", "./cmd/web"]
