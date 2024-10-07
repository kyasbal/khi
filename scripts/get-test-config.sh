#!/bin/bash

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
