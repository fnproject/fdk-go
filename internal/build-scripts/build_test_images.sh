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
# This scripts triggers all test images build for each supported go version and each type of test
#
set -e +x

# Login to OCIR
echo ${OCIR_PASSWORD} | docker login --username "${OCIR_USERNAME}" --password-stdin ${OCIR_REGION}

# Build and push the test function images to OCIR for integration test framework.
# Go 1.19
(
  source internal/build-scripts/build_test_image.sh internal/tests-images/go1.19/hello-world-test 1.19
  source internal/build-scripts/build_test_image.sh internal/tests-images/go1.19/runtime-version-test 1.19
  source internal/build-scripts/build_test_image.sh internal/tests-images/go1.19/timeout-test 1.19
  source internal/build-scripts/build_test_image.sh internal/tests-images/go1.19/oci-sdk-test 1.19
)

# Go 1.18
(
  source internal/build-scripts/build_test_image.sh internal/tests-images/go1.18/runtime-version-test 1.18
)
