FROM fnproject/go:dev as build-stage
ADD . /go/src/func/
ENV GO111MODULE=on 
RUN cd /go/src/func/ && go mod vendor && go build -o func
FROM fnproject/go
WORKDIR /function
COPY --from=build-stage /go/src/func/func /function/
ENTRYPOINT ["./func"]
