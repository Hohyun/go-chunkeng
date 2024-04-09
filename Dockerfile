# syntax=docker/dockerfile:1

FROM golang:1.21.9-bullseye

WORKDIR /app
COPY  . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-chunkeng ./cmd/main.go

EXPOSE 3011
CMD ["/go-chunkeng"]
