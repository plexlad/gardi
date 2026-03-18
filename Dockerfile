# syntax=docker/dockerfile:1

FROM node:20 AS frontend-builder
WORKDIR /app/client
COPY client/package*.json ./
RUN npm install
COPY client/ ./
RUN npm run build

FROM golang:1.26.1
WORKDIR /app/server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gardi .
COPY --from=frontend-builder /app/client/dist /app/client/dist
EXPOSE 8080
CMD ["/gardi"]
