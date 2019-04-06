FROM golang:1.12-alpine as builder

COPY . /
RUN apk -v --update add git gcc libc-dev && rm /var/cache/apk/* \
 && go build -o /main /main.go

FROM alpine:3.9

COPY run.sh /
COPY --from=builder /main /
RUN apk -v --update add \
        python \
        py-pip \
        groff \
        less \
        mailcap \
        ca-certificates \
        xz \
        && \
    pip install --upgrade awscli && \
    apk -v --purge del py-pip && \
    rm /var/cache/apk/* && \
    chmod a+x /run.sh

ENTRYPOINT /run.sh