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
# This script orchestrat the all steps needed to build and release the fdk-go images
#
set -xe

BUILD_VERSION=${BUILD_VERSION:-1.0.0-SNAPSHOT}
LOCAL=${LOCAL:-true}

export BUILD_VERSION
export LOCAL

#Unit test run section
(
  mydir=$(cd "$(dirname "$0")"; pwd)

  echo "LOCAL[${LOCAL}]"
  if [ ${LOCAL} = "true" ]; then
     cd ${mydir}/../../
  fi

  docker build -t fdk_go_test_build -f ./internal/Dockerfile_unit_test_run .

  #Run docker container and run the unit tests
  docker run  -v $(pwd):/build  fdk_go_test_build ./internal/build-scripts/execute_unit_tests.sh --rm
)

(
  # Build base fdk build and runtime
  source internal/build-scripts/build_base_images.sh
)

(
  # Build the test integration images
  source internal/build-scripts/build_test_images.sh
)

