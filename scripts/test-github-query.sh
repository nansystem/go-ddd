#!/bin/bash

TOKEN="$GITHUB_TOKEN"

# Define the GraphQL query
query='
{
  "query": "query { viewer { login } }"
}
'

# Make the API call
curl -s -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -X POST -d "$query" https://api.github.com/graphql
