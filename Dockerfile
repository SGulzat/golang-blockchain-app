FROM golang:1.22.8-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main7 main7.go

CMD ["./main7"]