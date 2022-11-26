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

set -xe

if [ -z "$1" ];then
  echo "Please supply Go runtime version as argument to release image." >> /dev/stderr
  exit 2
fi

go_version=$1
user="fnproject"
image="go"

echo "Pushing release images for Go Runtime Version ${go_version}"

./regctl image copy ${OCIR_REGION}/${OCIR_LOC}/gofdk:${go_version}-${BUILD_VERSION}-dev ${user}/${image}:${go_version}-dev
./regctl image copy ${OCIR_REGION}/${OCIR_LOC}/gofdk:${go_version}-${BUILD_VERSION} ${user}/${image}:${go_version}


