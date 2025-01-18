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
RUN go build -o goapp

FROM golang:1.23 as run

WORKDIR /app
# Copy the application executable from the build image
COPY --from=build /app/goapp .
COPY --from=build /app/templates ./templates
COPY --from=build /app/assets ./assets

EXPOSE 8080
RUN pwd
RUN ls
CMD ["/app/goapp"]
