#!/usr/bin/env bash

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

#
# This scripts builds a specific test image for given go version
#
set -xe

if [ -z "$1" ]; then
  echo "Please supply function directory to build test function image" >>/dev/stderr
  exit 2
fi

if [ -z "$2" ]; then
  echo "Please supply go version as argument to build image." >>/dev/stderr
  exit 2
fi

fn_dir=$1
go_version=$2
pkg_version=${BUILD_VERSION}

(
  #Add the fdk related source code
  source internal/build-scripts/copy_source_pkg.sh $fn_dir

  # Build test function image for integration test.
  pushd ${fn_dir}

  name="$(awk '/^name:/ { print $2 }' func.yaml)"
  echo "name:$name"

  version="$(awk '/^runtime:/ { print $2 }' func.yaml)"
  echo "version:$version"

  image_identifier="${version}${go_version}-${BUILD_VERSION}"
  echo "image_identifier:$image_identifier"

  #OCIR path
  ocir_image="${OCIR_LOC}/${name}:${image_identifier}"

  docker buildx build --push --platform linux/amd64,linux/arm64 -t "${OCIR_REGION}/${ocir_image}" -f Build_file --build-arg GO_VERSION=${go_version} --build-arg OCIR_REGION=${OCIR_REGION} --build-arg OCIR_LOC=${OCIR_LOC} --build-arg BUILD_VERSION=${BUILD_VERSION} .
  popd

  #Delete the fdk related source code
  source internal/build-scripts/cleanup_source_pkg.sh $fn_dir
)
