FROM golang:1.19.0-buster
ARG test_name
ENV my_env_var=TestSlow8PeerPoW
WORKDIR /app
RUN echo "The ARG variable value is $test_name, the ENV variable value is $my_env_var"
#COPY go.mod ./
COPY . .
RUN go mod download
#ENV my_env_var="TestSlow8PeerPoW"
#COPY *.go ./

#RUN ["go","test","-timeout","3600s","-run","TestSlow4PeerPoS","./test"]
RUN echo "The ARG variable value is $test_name, the ENV variable value is $my_env_var"
#EXPOSE 8080
#ENTYPOINT["$my_env_var"] \
Entrypoint go test -timeout 3600s -run $my_env_var ./test
#Entrypoint echo "$my_env_var"
#ENTRYPOINT ["go","test","-timeout","3600s","-run","TestSlow8PeerPoW","./test"]
#ENTRYPOINT ["go","test","-timeout","3600s","-run","TestSlow8PeerPoW","./test"]

#-------------------------------

#FROM golang:1.19.0-buster
#ARG test_name
#ENV my_env_var=$test_name
#WORKDIR /app
#
##COPY go.mod ./
#COPY . .
#RUN go mod download
#
##COPY *.go ./
#
#RUN #go build -o  /go-docker-demo
#
#EXPOSE 8080
#
#ENTRYPOINT ["go","test","-timeout","3600s","-run","$my_env_var","./test"]
#---------------------------------------


#ENTRYPOINT["./test"]
#CMD["go","test","-timeout","3600s","-run","$my_env_var","./test"]
#ENTRYPOINT ["go","test","-timeout","3600s","-run","TestSlow8PeerPoW","./test"]
#ENTRYPOINT ["ls"]
#ENTRYPOINT ["go","test","-timeout","3600s","-run","TestSlow4PeerPoS","./test"]
#ENTRYPOINT ["ls"]