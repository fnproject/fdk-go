#!/bin/sh
#
# Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
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

set -euo pipefail

# Pipe stdout to file, as Fn CLI init-image reads stdout to untar (outputting messages to stderr is fine).
go mod init "$FN_FUNCTION_NAME" > /tmp/go_mod_output
go get > /tmp/go_get_output

tar c go.mod func.go func.init.yaml
