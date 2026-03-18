# syntax=docker/dockerfile:1

# Stage 1: Build the Svelte frontend
FROM node:20 as frontend-builder
WORKDIR /app/client
COPY client/package*.json ./
RUN npm install
COPY client/ ./
RUN npm run build

# Stage 2: Build the Go backend and serve the Svelte frontend
FROM golang:1.25.4
WORKDIR /app

# Copy go.mod and go.sum from the root of the project
COPY go.mod go.sum ./

# Copy server files
COPY server/ ./server/

# Download Go modules
WORKDIR /app/server
RUN go mod download

# Compile the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /gardi .

# Copy frontend build output from the first stage
COPY --from=frontend-builder /app/client/dist /app/client/dist

EXPOSE 8080

# Run the binary
CMD ["/gardi"]
