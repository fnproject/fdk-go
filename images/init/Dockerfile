FROM alpine:3.10

COPY boilerplate /boilerplate

ARG tag=latest
RUN sed -i "s/TAG/${tag}/" /boilerplate/func.init.yaml

WORKDIR /boilerplate
CMD ["tar", "c", "Gopkg.toml", "func.go", "func.init.yaml"]
