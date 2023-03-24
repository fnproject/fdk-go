#
# Copyright (c) 2021, 2022 Oracle and/or its affiliates. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

FROM --platform=$TARGETPLATFORM ocr-docker-remote.artifactory.oci.oraclecorp.com/os/oraclelinux:8-slim

RUN set -eux; \
	microdnf update; \
	microdnf install \
    glibc wget gzip tar; \
  microdnf clean all;

ENV PATH /usr/local/go/bin:$PATH

ENV GOLANG_VERSION 1.19.4

RUN set -eux; \
	arch="$(arch)"; \
	url=; \
	case "$arch" in \
		'x86_64') \
			url='https://dl.google.com/go/go1.19.4.linux-amd64.tar.gz'; \
            sha256='c9c08f783325c4cf840a94333159cc937f05f75d36a8b307951d5bd959cf2ab8'; \
			;; \
		'arm64') \
			url='https://dl.google.com/go/go1.19.4.linux-arm64.tar.gz'; \
            sha256='9df122d6baf6f2275270306b92af3b09d7973fb1259257e284dba33c0db14f1b'; \
			;; \
		'aarch64') \
			url='https://dl.google.com/go/go1.19.4.linux-arm64.tar.gz'; \
            sha256='9df122d6baf6f2275270306b92af3b09d7973fb1259257e284dba33c0db14f1b'; \
			;; \
		*) echo >&2 "error: unsupported architecture '$arch' (likely packaging update needed)"; exit 1 ;; \
	esac; \
	if [ -z "$url" ]; then \
		echo >&2; \
		echo >&2 "warning: current architecture ($arch) does not have a compatible Go binary release"; \
		echo >&2; \
	fi; \
	\
	wget -O go.tgz.asc "$url.asc"; \
	wget -O go.tgz "$url" --progress=dot:giga; \
	echo "$sha256 *go.tgz" | sha256sum -c -; \
	GNUPGHOME="$(mktemp -d)"; \
    export GNUPGHOME; \
	gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 'EB4C 1BFD 4F04 2F6D DDCC EC91 7721 F63B D38B 4796'; \
	gpg --batch --verify go.tgz.asc go.tgz; \
	gpgconf --kill all; \
	rm -rf "$GNUPGHOME" go.tgz.asc; \
    mkdir -p /usr/local; \
	gzip -d go.tgz; \
	tar -C /usr/local -xpf go.tar; \
	rm go.tar; \
	go version;

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH
CMD ["/bin/go", "version"]
