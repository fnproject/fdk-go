#!/usr/bin/env bash

set -xe

goversion=${1:-"1.12.4"}
ostype=${2:-"alpine"}

pushd images/${goversion}/${ostype}/build-stage/ && docker build -t fnproject/golang:${goversion}-${ostype}-dev . && popd
pushd images/${goversion}/${ostype}/runtime/ && docker build -t fnproject/golang:${goversion}-${ostype} . && popd
