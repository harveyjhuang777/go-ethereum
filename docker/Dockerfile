FROM golang:1.15 AS go-builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY . .
RUN go build -o srv /app/cmd/

FROM alpine:3.15
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates && \
    apk add git
WORKDIR /app
COPY --from=go-builder /app/srv /app/srv
COPY --from=go-builder /app/conf.d /app/conf.d
ENTRYPOINT ["/app/srv"]