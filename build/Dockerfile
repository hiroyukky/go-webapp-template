FROM golang:1.19.1-alpine

RUN apk update && apk add git

RUN mkdir -p /app
WORKDIR /app
ADD . /app

RUN rm go.mod
RUN go mod init go-template && \
    go mod edit -replace go-template/cmd/app=./cmd/app && \
    go mod tidy

CMD ["go", "run", "cmd/app/main.go"]
