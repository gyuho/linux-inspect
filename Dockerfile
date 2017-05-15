FROM ubuntu:17.04
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections

##########################
# invalidates cache for every commit
RUN apt-get -y update
RUN apt-get -y install curl
RUN curl -o /version-git.json https://api.github.com/repos/gyuho/linux-inspect/git/refs/heads/master
##########################

##########################
# initialize Ubuntu
RUN apt-get -y install \
  build-essential apt-utils gcc bash tar \
  apt-transport-https libssl-dev git

RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get -y dist-upgrade
RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get -y autoremove
RUN apt-get -y autoclean

RUN cat /etc/lsb-release > /version-container.txt
RUN printf "\n" >> /version-container.txt
RUN uname -a >> /version-container.txt
RUN printf "\n" >> /version-container.txt
##########################

##########################
ENV HOME_DIR /
WORKDIR ${HOME_DIR}
ENV GIT_PATH github.com/gyuho/linux-inspect

# install Go
RUN rm -rf /root/go1.4
RUN mkdir -p /root/go1.4
RUN curl -s https://storage.googleapis.com/golang/go1.4.linux-amd64.tar.gz | tar -v -C /root/go1.4 -xz --strip-components=1
ENV GOROOT /goroot-tip
RUN rm -rf ${GOROOT}
ENV PATH ${GOROOT}/bin:/usr/local/go/bin:$PATH
RUN git clone https://go.googlesource.com/go ${GOROOT}
WORKDIR ${GOROOT}/src
RUN git reset --hard HEAD
RUN ./make.bash
RUN ${GOROOT}/bin/go version
RUN cp ${GOROOT}/bin/go /bin/go
RUN go version
WORKDIR ${HOME_DIR}

# install linux-inspect
ENV GOPATH /gopath-tip
ENV PATH ${GOPATH}/bin:/usr/local/go/bin:$PATH
RUN mkdir -p ${GOPATH}/src/github.com/gyuho
RUN git clone https://github.com/gyuho/linux-inspect --branch master ${GOPATH}/src/${GIT_PATH}
WORKDIR ${GOPATH}/src/${GIT_PATH}
RUN git reset --hard HEAD
RUN go get -d -v ./cmd/linux-inspect
RUN go get -d -v github.com/coreos/etcd/pkg/netutil
RUN go install -v ./cmd/linux-inspect
RUN linux-inspect -h

RUN ./test
##########################

##########################
ENV HOME_DIR /
WORKDIR ${HOME_DIR}
ENV GIT_PATH github.com/gyuho/linux-inspect

# install Go
ENV GOROOT /usr/local/go
ENV GOPATH /gopath-1.8.1
ENV PATH ${GOPATH}/bin:/usr/local/go/bin:$PATH
ENV GO_VERSION 1.8.1
ENV GO_DOWNLOAD_URL https://storage.googleapis.com/golang
RUN rm -rf ${GOROOT}
RUN curl -s ${GO_DOWNLOAD_URL}/go${GO_VERSION}.linux-amd64.tar.gz | tar -v -C /usr/local/ -xz
RUN go version
WORKDIR ${HOME_DIR}

# install linux-inspect
ENV GOPATH /gopath-1.8.1
ENV PATH ${GOPATH}/bin:/usr/local/go/bin:$PATH
RUN mkdir -p ${GOPATH}/src/github.com/gyuho
RUN git clone https://github.com/gyuho/linux-inspect --branch master ${GOPATH}/src/${GIT_PATH}
WORKDIR ${GOPATH}/src/${GIT_PATH}
RUN git reset --hard HEAD
RUN go get -d -v ./cmd/linux-inspect
RUN go get -d -v github.com/coreos/etcd/pkg/netutil
RUN go install -v ./cmd/linux-inspect
RUN linux-inspect -h

RUN ./test
##########################

##########################
WORKDIR ${GOPATH}/src/${GIT_PATH}
RUN pwd

RUN cat /version-git.json
RUN cat /version-container.txt
##########################
