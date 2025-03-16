FROM golang:alpine

RUN mkdir /app

WORKDIR $GOPATH/src/github.com/mboiar/swift-restful
COPY go.mod go.sum $GOPATH/src/github.com/mboiar/swift-restful/
RUN go mod download

COPY . $GOPATH/src/github.com/mboiar/swift-restful/
RUN go build -o main main.go

EXPOSE 8000

CMD ["./main"]