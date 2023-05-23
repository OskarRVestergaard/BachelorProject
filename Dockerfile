FROM golang:1.19.0-buster

WORKDIR /app

#COPY go.mod ./
COPY . .
RUN go mod download

#COPY *.go ./

RUN #go build -o  /go-docker-demo

EXPOSE 8080

ENTRYPOINT ["go","test","-run","TestSlow8PeerPoS","./test"]
#ENTRYPOINT ["ls"]