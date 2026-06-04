FROM golang:1.22-alpine

WORKDIR /app

RUN adduser -D -u 10001 appuser

COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o server .

USER appuser

ENV PING_LISTEN_PORT=8080
EXPOSE 8080

CMD ["./server"]
