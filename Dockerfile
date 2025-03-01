FROM node:18 as tailwind

WORKDIR /app

# Copy package.json and package-lock.json if they exist
COPY package.json .
COPY package-lock.json .

# Install npm dependencies
RUN npm install

# Copy the rest of the application files
COPY static static

# Run the Tailwind CLI to build the CSS
RUN npx tailwindcss -i ./static/input.css -o ./static/output.css --minify

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

# Build the app
RUN go build -o goapp ./cmd/web
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
