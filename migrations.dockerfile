FROM golang:1.23 as build

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# Copy everything
COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapp ./cmd/migrate

FROM scratch as run

WORKDIR /app
# Copy the application executable from the build image
COPY --from=build /app/goapp .
COPY --from=build /app/migrations ./migrations

ENTRYPOINT ["/app/goapp"]
