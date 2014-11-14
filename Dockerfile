# Dockerfile for development to emulate Linux environment

FROM debian:wheezy

RUN apt-get -y update && apt-get -y install curl sudo unzip git build-essential ruby

RUN mkdir -p /opt/local
RUN curl -L https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz -o /tmp/go.tar.gz
RUN tar -C /opt/local -xzf /tmp/go.tar.gz
RUN (cd /opt/local/go/src; GOOS=linux GOARCH=386 ./make.bash --no-clean)

RUN mkdir -p /go/src/github.com/mackerelio/mackerel-agent
COPY . ./go/src/github.com/mackerelio/mackerel-agent

ENV GOPATH /go
ENV GOROOT /opt/local/go
ENV PATH /opt/local/go/bin:/bin:/usr/bin
ENV GOOS linux
ENV GOARCH 386

WORKDIR /go/src/github.com/mackerelio/mackerel-agent
CMD ["make", "all"]
