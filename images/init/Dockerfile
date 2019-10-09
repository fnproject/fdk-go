FROM golang:1.12.7-alpine3.10
RUN apk add --no-cache git

COPY boilerplate /boilerplate

WORKDIR /boilerplate
CMD ["/bin/sh", "mod.sh"]
