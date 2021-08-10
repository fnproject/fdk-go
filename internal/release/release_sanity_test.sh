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

cd ..
cp -R sanity_testing/internal .bitbucket/internal
cd .bitbucket

ls -al

docker build -t fdk_go_test_build -f ./internal/Dockerfile_unit_test_run .

#Run docker container and run the unit tests
docker run  -v $(pwd):/build  fdk_go_test_build ./internal/build-scripts/execute_unit_tests.sh --rm

# Remove the copied over internal folder for sanity testing
rm -rf internal
rm -rf ./../sanity_testing


# this step is to ensure we don't commit any files produced by unit test or build package to github branch
git status
if [[ -z $(git status -s) ]]
then
  echo "tree is clean"
else
  echo "tree is dirty, please commit changes before running this"
  exit 1
fi