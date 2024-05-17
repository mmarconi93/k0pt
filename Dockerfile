FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go build -o kopt cmd/main.go

ENTRYPOINT ["./kopt"]
