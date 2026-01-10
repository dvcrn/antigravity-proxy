#!/bin/bash

curl 'http://localhost:9888/v1beta/models/gemini-2.5-pro:generateContent' \
  -X POST \
  -H 'Authorization: Bearer Zo3CyNtqVbzRYhZL3JWarLHfiB3YXU3Mkv34' \
  -H 'User-Agent: GeminiCLI/v23.5.0 (darwin; arm64) google-api-nodejs-client/9.15.1' \
  -H 'x-goog-api-client: gl-node/23.5.0' \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  --data-raw '{
    "contents": [
      {
        "role": "user",
        "parts": [
          {
            "text": "Say hello and tell me what 2+2 is."
          }
        ]
      }
    ]
  }' | jq
