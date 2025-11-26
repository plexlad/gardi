# syntax=docker/dockerfile:1
FROM golang:1.25.4

# set server directory
WORKDIR /server

# copy binary files and download dependencies
COPY * ./
RUN go mod download

# compile the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /gardi

EXPOSE 5499

# run the binary
CMD ["/gardi"]
