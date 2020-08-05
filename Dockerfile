# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.13.9-alpine3.11 AS builder

# Copy the local package files to the container's workspace.
WORKDIR /go/src/book
COPY . .

RUN apk add --no-cache curl \
    bash \
    git

COPY go.mod go.sum ./
RUN go mod download

#note: /go/src/book/bin/book is bad place, because of binded docker-compose volumes not have compiled binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/book /go/src/book/cmd/book

FROM alpine:3.11
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/book /go/src/book
COPY --from=builder /bin/book /bin/book

CMD ["/bin/book"]

# Document that the service listens on port 8080.
EXPOSE 8080
