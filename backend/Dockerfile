FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./

COPY *.go ./

RUN CGO_ENABLE=0 GOOS=linux go build -o /streamer

CMD ["/streamer"]
