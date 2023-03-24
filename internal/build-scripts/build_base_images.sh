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
  
  #Create the builder instance
  docker buildx rm builderInstance || true
  docker buildx create --name builderInstance --driver-opt=image=iad.ocir.io/oraclefunctionsdevelopm/moby/buildkit:buildx-stable-1 --platform linux/amd64,linux/arm64
  docker buildx use builderInstance

  #Teamcity uses a very old version of buildx which creates a bad request body. Pushing the images to OCIR gives a 400 bad request error. Hence, use this 
  #script to upgrade the buildx version.
  ./internal/build-scripts/update-buildx.sh

  # Build base fdk build and runtime images
  ./internal/build-scripts/build_base_image.sh 1.19
  ./internal/build-scripts/build_base_image.sh 1.18
)
