FROM golang:alpine

LABEL authors="kstmc"

WORKDIR /ozon

EXPOSE 8080

COPY . .

RUN go build -o main ./cmd/server.go

CMD ["./main"]