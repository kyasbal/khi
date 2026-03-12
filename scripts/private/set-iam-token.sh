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


set -e
set -u

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <project-id>"
    echo "This script prompts for an IAM token and sends it to the local KHI server."
    exit 1
fi

PROJECT_ID=$1
PORT=${KHI_PORT:-8080}
API_PATH="/api/v3/iam_token"

echo "Please paste your IAM Token for project '${PROJECT_ID}'."
echo "(Press Ctrl+D on a new line when finished):"

# Read multi-line input until EOF
RAW_TOKEN=$(cat)

# Strip out backslashes and newlines (e.g. \ \n)
CLEAN_TOKEN=$(echo "$RAW_TOKEN" | tr -d '\\\n' | tr -d '\r' | xargs)

if [ -z "$CLEAN_TOKEN" ]; then
    echo
    echo "Error: Token cannot be empty."
    exit 1
fi

echo
echo "Setting IAM token for project: ${PROJECT_ID} on port ${PORT}..."

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://localhost:${PORT}${API_PATH}" \
    -H "Content-Type: application/json" \
    -d "{
        \"projectID\": \"${PROJECT_ID}\",
        \"iamToken\": \"${CLEAN_TOKEN}\"
    }")

if [ "$HTTP_CODE" -eq 200 ]; then
    echo "Success! (HTTP 200)"
else
    echo "Failed to set token. (HTTP ${HTTP_CODE})"
    exit 1
fi
