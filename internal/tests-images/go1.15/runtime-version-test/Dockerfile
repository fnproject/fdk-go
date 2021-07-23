ARG GO_VERSION=latest
FROM fnproject/go:${GO_VERSION}-dev as build-stage
WORKDIR /function
WORKDIR /go/src/func/
ENV GO111MODULE=on
COPY . .
RUN go build -o func -v
FROM fnproject/go:${GO_VERSION}
WORKDIR /function
COPY --from=build-stage /go/src/func/func /function/
ENTRYPOINT ["./func"]