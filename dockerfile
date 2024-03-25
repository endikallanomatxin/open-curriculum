FROM golang:alpine

WORKDIR /app

COPY main.go .

RUN go build main.go

# Expose port 8080 to the outside world

EXPOSE 8080

# Command to run the executable
CMD ["./main"]
