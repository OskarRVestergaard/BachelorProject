FROM golang:1.19.0-buster
ARG test_name
ENV my_env_var=TestSlow8PeerPoW
WORKDIR /app
COPY . .
RUN go mod download

Entrypoint go test -timeout 3600s -run $my_env_var ./test
