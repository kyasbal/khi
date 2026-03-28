#!/bin/bash
# Copyright 2026 Google LLC
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

# 一時的にGoogleCloudPlatform/khi を kyasbal/khi に置換するスクリプト

echo "Replacing GoogleCloudPlatform/khi with kyasbal/khi across the project..."

# macOS の sed を使用して全ファイルを置換
find . -type f \( \
  -name '*.go' -o \
  -name 'go.mod' -o \
  -name 'go.sum' -o \
  -name '*.md' -o \
  -name '*.mk' \
\) -print0 | xargs -0 sed -i '' 's|github.com/GoogleCloudPlatform/khi|github.com/kyasbal/khi|g'

echo "Done! You can now commit and push to kyasbal/khi."
