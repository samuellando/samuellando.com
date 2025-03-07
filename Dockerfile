FROM node:22 as tailwind

WORKDIR /app

# Copy package.json and package-lock.json if they exist
COPY package.json .
COPY package-lock.json .

# Install npm dependencies
RUN npm install

# Copy the the application files
COPY . .

# Run the Tailwind CLI to build the CSS
RUN npx tailwindcss -i ./static/input.css -o ./static/output.css --minify

FROM debian as sqlc
WORKDIR /app
COPY --from=sqlc/sqlc:1.28.0 /workspace/sqlc /usr/bin/sqlc
COPY queries ./queries
COPY migrations ./migrations
COPY sqlc.yml . 
RUN sqlc generate

FROM golang:1.23 as build

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# Copy everything
COPY . .
COPY --from=tailwind /app/static ./static
COPY --from=sqlc /app/data ./data

# Build the app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapp ./cmd/web
# Build wasm
RUN GOOS=js GOARCH=wasm go build -o static/main.wasm ./wasm

FROM scratch as run

WORKDIR /app
# Copy the application executable from the build image
COPY --from=build /app/goapp .
COPY --from=build /app/templates ./templates
COPY --from=build /app/static ./static
COPY --from=build /app/migrations ./migrations

EXPOSE 8080
ENTRYPOINT ["/app/goapp"]
