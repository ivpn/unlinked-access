# stage 1: building application binary file
FROM golang:1.25 AS builder

RUN mkdir /verifier
ADD services/verifier /verifier
WORKDIR /verifier

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o verifier main.go

# stage 2: copy only the application binary file and necessary files to the alpine container
FROM alpine:latest AS production

COPY --from=builder /verifier .

# run the service on container startup
CMD ["./verifier", "serve"]
