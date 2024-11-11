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


TEST_CASE_ID=$1

TEST_CASE=$(cat .vscode/khi-test-cases.json | jq ".[]|select(.test_case_id ==\"${TEST_CASE_ID}\")")
echo "test case: $TEST_CASE"
echo "Requesting test scenario for $(echo $TEST_CASE | jq '.title') ($(echo $TEST_CASE | jq '.test_case_id' -r))"
INSPECTION_ID=$(curl -X POST localhost:8080/api/v2/inspection/types/gcp-gke | jq ".inspectionId" -r)

# enables all for this moment
SET_FEATURE_REQUEST=$(curl localhost:8080/api/v2/inspection/tasks/${INSPECTION_ID}/features | jq "{features:[.features[].id]}" -r)
curl -X PUT -H "Content-Type: application/json" -d "${SET_FEATURE_REQUEST}" localhost:8080/api/v2/inspection/tasks/${INSPECTION_ID}/features

PARAMS=$(echo $TEST_CASE|jq '.query_parameter' -r)
echo Requesting with following parameter
echo $PARAMS | jq
curl -X POST -H "Content-Type: application/json" -d "$(echo $PARAMS | jq -r)" localhost:8080/api/v2/inspection/tasks/${INSPECTION_ID}/run