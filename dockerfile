FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build main.go

RUN go mod download

EXPOSE 8080

CMD ["./main"]
