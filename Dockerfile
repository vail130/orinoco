FROM golang:1.8

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        libyaml-dev \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /go/src/github.com/vail130/orinoco

WORKDIR /go/src/github.com/vail130/orinoco

ENV PATH /go/src/github.com/vail130/orinoco/bin:$PATH

ARG GIT_COMMIT
ENV GIT_COMMIT ${GIT_COMMIT}

COPY . /go/src/github.com/vail130/orinoco/

# `touch Makefile` avoids clock skew warnings
RUN touch Makefile && make configure build
