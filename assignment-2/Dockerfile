FROM golang:latest

ENV GO111MODULE=on

WORKDIR /app/assigment-2

COPY . .

RUN go mod tidy

RUN go build -o main .

CMD ["./main"]