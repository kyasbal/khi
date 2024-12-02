#!/bin/bash
# Copyright 2024 Google LLC
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


# This is temporal adhoc script to extract the internal KHI implementation for OSSing.
# Eventually the OSSed version should be the base and internal code should be amended in internal repository. But for this time, extract the OSS part from internal repository.
# The dependency graph will be flipped once we done OSSing.

OSSIGNORE=$(sed '/^[ \t]*#/d' $1) # Removing comment lines

echo "$OSSIGNORE" | while read -r pattern;
do
    rm -rf $pattern -r
done