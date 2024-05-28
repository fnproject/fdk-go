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
# This scripts triggers a build for all specific fdk go base image for given go version
#
set -xe

(
  #Login to OCIR
  echo "${OCIR_PASSWORD}" | docker login --username "${OCIR_USERNAME}" --password-stdin ${OCIR_REGION}

  # Build base fdk build and runtime images
  ./internal/build-scripts/build_base_image.sh 1.20
  ./internal/build-scripts/build_base_image.sh 1.19
)
