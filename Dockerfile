FROM golang:1.19 AS build

WORKDIR /app

COPY . .

RUN go build -o main .

CMD ["./main"]