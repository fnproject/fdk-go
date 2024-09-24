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
# This scripts add folders/filesrelated to fdk-go(that will be needed during test image build)
#

set -xe

if [ -z "$1" ]; then
  echo "Please supply function directory to build test function image" >>/dev/stderr
  exit 2
fi

fn_dir=$1

(
 #Prepare the path for copying source file related to fdk
  echo "Function directory provided[${fn_dir}]"

  echo "Trying to create vendor folder[${fn_dir}/vendor/github.com/fnproject/fdk-go/]"
  mkdir -p "${fn_dir}/vendor/github.com/fnproject/fdk-go/"

  # Clean-up if any left over
  echo "Clean-up vendor folder[${fn_dir}/vendor/github.com/fnproject/fdk-go/]"
  rm -rf ${fn_dir}/vendor/github.com/fnproject/fdk-go/*

  #Copy fdk related file one-by-one to vendor specific directory
  echo "Copying fdk-go related source files to vendor folder[${fn_dir}/vendor/github.com/fnproject/fdk-go/]"
  cp  ./fdk.go ${fn_dir}/vendor/github.com/fnproject/fdk-go/
  cp ./fdk_example_test.go ${fn_dir}/vendor/github.com/fnproject/fdk-go/
  cp ./fdk_test.go ${fn_dir}/vendor/github.com/fnproject/fdk-go/
  cp ./go.mod ${fn_dir}/vendor/github.com/fnproject/fdk-go/
  cp ./handler.go ${fn_dir}/vendor/github.com/fnproject/fdk-go/
  cp ./version.go ${fn_dir}/vendor/github.com/fnproject/fdk-go/

  echo "Listing copied files at vendor folder[${fn_dir}/vendor/github.com/fnproject/fdk-go/]"
  #Listing all copied files name
  ls -alRt  "${fn_dir}"

)