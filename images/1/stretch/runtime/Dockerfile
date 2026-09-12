FROM debian:stretch-slim

RUN apt-get update && apt-get upgrade -qy && apt-get clean
RUN addgroup --system --gid 1000 --system fn && adduser --system --uid 1000 --ingroup fn fn
