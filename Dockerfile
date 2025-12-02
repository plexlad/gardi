# syntax=docker/dockerfile:1
FROM golang:1.25.4

WORKDIR /server

# download dependencies
COPY * ./
RUN go mod download

# compile the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /gardi

EXPOSE 8080

# run the binary
CMD ["/gardi"]
