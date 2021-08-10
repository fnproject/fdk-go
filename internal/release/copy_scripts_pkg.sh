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

set -ex

# list checkout directory contents
ls -al

# Create sanity_testing/internal folder at the checkout directory
mkdir -p sanity_testing
rm -rf sanity_testing/*
cd sanity_testing

mkdir -p internal
rm -rf internal/*
cd internal

mkdir -p build-scripts
rm -rf build-scripts/*

mkdir -p release
rm -rf release/*

cd ../../

# copy the needed files
cp .bitbucket/internal/Dockerfile_unit_test_run sanity_testing/internal/Dockerfile_unit_test_run
cp .bitbucket/internal/build-scripts/execute_unit_tests.sh sanity_testing/internal/build-scripts/execute_unit_tests.sh
cp .bitbucket/internal/release/release_sanity_test.sh sanity_testing/internal/release/release_sanity_test.sh

ls -alRt .