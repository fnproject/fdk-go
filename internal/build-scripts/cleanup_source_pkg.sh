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
# This scripts clean-up all temporary folders/files got created during test image build
#
set -xe

if [ -z "$1" ]; then
  echo "Supply function directory to build test function image" >>/dev/stderr
  exit 2
fi

fn_dir=$1

if [[ -d $fn_dir/vendor ]]; then
    echo "Cleaning the vendor directory[$fn_dir/vendor]"
    rm -rf "$fn_dir/vendor"
fi


