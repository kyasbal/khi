#!/bin/bash

# generate-testlog-data.sh
# This script receives a file path that includes YAML format logs and replace project ID, project number...etc to use the log in test cases.

INPUT_JSON=$(cat $1|yq e -o json|jq .)

REQUEST_JSON='{
    "contents": [{
      "parts":[
        {"text": ""
        }
      ]
    }],
    "generationConfig": { "response_mime_type": "application/json" }
}'

# Concat the prompt text and given YAML data converted to JSON.
# This script receives the log data as YAML file, but give it to Gemini in JSON and receives JSON response from the API.
# This is because Gemini API is not supporting YAML data as output of structured data output mode. 
REQUEST_JSON=$(echo $REQUEST_JSON|jq --arg prompt "$(cat ./scripts/private/prompts/convert-to-testingdata.txt)" --arg content "`echo -E "$INPUT_JSON"`" '.contents[0].parts[0].text = ($prompt + "\n" + $content)')

OUTPUT=`curl "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=${GOOGLE_AI_KEY}" \
-H 'Content-Type: application/json' \
-d "$REQUEST_JSON" \
2> /dev/null`

RESULT_YAML=$(echo $OUTPUT | jq .candidates[0].content.parts[0].text -r | yq -p json -r)

# Show the diff.
diff  -u -U 10000000 --color -b <(echo "$INPUT_JSON" |yq -p json -r) <(echo "$RESULT_YAML")

read -p "Overwrite? [y/N]" -n 1 -r

if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo -E "$RESULT_YAML" > $1
fi