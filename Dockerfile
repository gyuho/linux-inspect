FROM ubuntu:17.10

ENV ROOT_DIR /
WORKDIR ${ROOT_DIR}
ENV HOME /root

RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections \
  && apt-get -y update \
  && apt-get -y install \
  build-essential \
  gcc \
  apt-utils \
  software-properties-common \
  curl \
  python \
  git \
  tar \
  bash \
  apt-transport-https \
  libssl-dev \
  && rm /bin/sh \
  && ln -s /bin/bash /bin/sh \
  && ls -l $(which bash) \
  && echo "root ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers \
  && apt-get -y clean \
  && rm -rf /var/lib/apt/lists/* \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y dist-upgrade \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y autoremove \
  && apt-get -y autoclean

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:${GOROOT}/bin:${PATH}
ENV GO_VERSION REPLACE_ME_GO_VERSION
ENV GO_DOWNLOAD_URL https://storage.googleapis.com/golang
RUN rm -rf ${GOROOT} \
  && curl -s ${GO_DOWNLOAD_URL}/go${GO_VERSION}.linux-amd64.tar.gz | tar -v -C /usr/local/ -xz \
  && mkdir -p ${GOPATH}/src ${GOPATH}/bin \
  && go version
