FROM golang:1.23

WORKDIR /app

# Some build dependencies
RUN apt update && apt install npm -y
COPY --from=sqlc/sqlc:1.28.0 /workspace/sqlc /usr/bin/sqlc

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# Copy everything
COPY . .

RUN npm install

# Build the wasm
RUN GOOS=js GOARCH=wasm go build -o static/main.wasm ./wasm

EXPOSE 8080
CMD sqlc generate && \ 
    npx tailwindcss -i ./static/input.css -o ./static/output.css --minify && \
    go run ./cmd/web
