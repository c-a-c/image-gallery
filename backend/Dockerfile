# Dockerfile

FROM golang:1.24

WORKDIR /build

COPY go.mod . go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o app main.go

EXPOSE 8080

CMD ["go", "run", "main.go"]
