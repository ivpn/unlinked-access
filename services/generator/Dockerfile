# stage 1: building application binary file
FROM golang:1.25 AS builder

RUN mkdir /generator
ADD services/generator /generator
ADD proto /proto
WORKDIR /generator

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o generator main.go

# stage 2: copy only the application binary file and necessary files to the alpine container
FROM alpine:latest AS production

COPY --from=builder /generator .

# run the service on container startup
CMD ["./generator", "serve"]
