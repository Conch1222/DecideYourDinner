FROM golang:1.23.1-alpine

WORKDIR /src

COPY . /src

RUN go mod tidy

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]