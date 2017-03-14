FROM golang:alpine

RUN apk update \
  && apk add curl bash binutils tar git \
  && rm -rf /var/cache/apk/* \
  && /bin/bash \
  && touch ~/.bashrc \
  && curl -o- -L https://glide.sh/get | bash \
  && apk del curl tar binutils

WORKDIR "${GOPATH}/src/github.com/duckclick/wing"

ADD glide.yaml glide.yaml
ADD glide.lock glide.lock
RUN glide install

ADD . .
RUN go build -ldflags "-s -w"

CMD ["./wing"]
