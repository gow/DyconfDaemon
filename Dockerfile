# sshd

FROM ubuntu:16.04
MAINTAINER Chetan Gowda <chetan.hs@gmail.com>


### install required packages ###
RUN apt-get update
RUN apt-get install -y --no-install-recommends \
  #git \
  #python \
  #curl \
  ca-certificates 
  #vim-nox \
  #vim \
  #openssh-server \
  #mosh \
  #sudo \
  #ctags

RUN mkdir -p /usr/sbin/

COPY ./DyconfDaemon /usr/sbin/
### ssh stuff
#RUN mkdir /var/run/sshd
#RUN set -i 's/Port 22/Port 2222/' /etc/ssh/sshd_config
#RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
#RUN sed -i 's/^X11Forwarding yes/X11Forwarding no/' /etc/ssh/sshd_config
#RUN sed -i 's/^LogLevel INFO/LogLevel VERBOSE/' /etc/ssh/sshd_config

### add user ###
RUN useradd chetan \
  --uid 1000 \
  --home-dir /home/chetan \
  --no-create-home \
  --groups sudo \
  --shell /bin/bash
RUN echo "chetan:chetan" | chpasswd
RUN mkdir -p /home/chetan
#WORKDIR /home/chetan

### Install Go ###
#RUN curl https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz -o /tmp/go.tar.gz
#RUN tar -C /usr/local -xzf /tmp/go.tar.gz

### Environment variables ###
#ENV GOPATH /home/chetan/go
#ENV PATH $GOPATH/bin/:$PATH:/usr/local/go/bin/

### Volumes ###
VOLUME /home/chetan

### expose port ###
EXPOSE 9009

#CMD ["/usr/sbin/sshd", "-D"]
