# Steam History server

FROM ubuntu:precise
MAINTAINER Roman Tsukanov <roman@tsukanov.me>

# Making sure the package repository is up to date
RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-get update

RUN apt-get install -y -q build-essential mercurial git-core wget memcached

# Setting up Go
RUN wget https://go.googlecode.com/files/go1.2.1.linux-amd64.tar.gz -q -P /tmp
RUN tar -C /usr/local -xzf /tmp/go1.2.1.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
RUN mkdir /var/go
ENV GOPATH /var/go
ENV PATH $PATH:$GOPATH/bin

# Getting sources and building
RUN go get -v -u github.com/tsukanov/steamhistory/...