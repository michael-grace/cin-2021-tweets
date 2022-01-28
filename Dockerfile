FROM golang:1.16

WORKDIR /usr/src/app

COPY . .

RUN go build ./cmd/tweets/

EXPOSE 3000

ENTRYPOINT ["./tweets"]