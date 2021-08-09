#!/bin/bash
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

set -exuo pipefail

# ensure working dir is clean
git status
if [[ -z $(git status -s) ]]
then
  echo "tree is clean"
else
  echo "tree is dirty, please commit changes before running this"
  exit 1
fi

git pull

version_file="version.go"

if [ -z $(grep -m1 -Eo "[0-9]+\.[0-9]+\.[0-9]+" $version_file) ]; then
  echo "did not find semantic version in $version_file"
  exit 1
fi

current_version=$(grep -m1 -Eo "[0-9]+\.[0-9]+\.[0-9]+" ${version_file})

echo "Current version: $current_version"

# Calculate new/release version
version_parts=(${current_version//./ })
new_minor=$((${version_parts[2]}+1))
new_version="${version_parts[0]}.${version_parts[1]}.$new_minor"

if [[ ${new_version} =~  ^[0-9]+\.[0-9]+\.[0-9]+$ ]] ; then
   echo "New Release version: $new_version"
else
   echo "Invalid Release version: $new_version"
   exit 1
fi

SEDOPTION="-i "
if [[ "$OSTYPE" == "darwin"* ]]; then
  SEDOPTION="-i .bak"
fi

sed $SEDOPTION -e s/$current_version/$new_version/g version.go

if [[ "$OSTYPE" == "darwin"* ]]; then
  rm -rf version.go.bak
fi

tag="v$new_version"
git add -u
git commit -m "FDK Go: $tag version release"
git tag -f -a $tag -m "version $tag"
git push
git push origin $tag