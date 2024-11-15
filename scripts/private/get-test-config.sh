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


# get-test-config.sh
# Downloads the latest set of configurations for KHI generated from khi-testing.
# A automated cluster creation/deletion system is running on the project `khi-testing`. 
# The time range and other parameters needed to query is saved on there. This script fetch the configuretion in local to test KHI easily against the data.

bq query --format prettyjson --use_legacy_sql=false --project_id kubernetes-history-inspector \
    'SELECT cases.creation_time,test_case_id,cases.title,inspection_type,query_parameter FROM kubernetes-history-inspector.testing.test_generations AS generations
LEFT JOIN kubernetes-history-inspector.testing.test_cases AS cases USING(generation_id)
WHERE generations.generation_id = (
  SELECT generation_id FROM kubernetes-history-inspector.testing.test_generations AS generations
  LEFT JOIN kubernetes-history-inspector.testing.test_cases USING(generation_id)
  ORDER BY generations.creation_time DESC
  LIMIT 1
)' > .vscode/khi-test-cases.json


echo "Test case list was updated:"
echo "Available test cases:"
echo "ID Title"
cat .vscode/khi-test-cases.json | jq ".[]|(.test_case_id)+\" \"+(.title)" -r   
