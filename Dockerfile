FROM golang:1.17.1-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o yandex-eda ./cmd/main.go

EXPOSE 8080
EXPOSE 8081
EXPOSE 5432
ENTRYPOINT ["./yandex-eda"]

