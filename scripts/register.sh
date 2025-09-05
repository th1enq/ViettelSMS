#!/bin/bash

SERVER_ID="$1"
SERVER_NAME="$2"
LOCATION="Data Center A"
OS="Linux"
INTERVAL_TIME="$3"
IPV4="$4"

RESPONSE=$(curl -X POST http://$HOST_IP:80/server/ \
  -H "Content-Type: application/json" \
  -H "X-API-Key: <API_KEY>" \
  -d "{
    \"server_id\": \"$SERVER_ID\",
    \"server_name\": \"$SERVER_NAME\",
    \"location\": \"$LOCATION\",
    \"os\": \"$OS\",
    \"ipv4\": \"$IPV4\",
    \"interval_time\": $INTERVAL_TIME
  }")

cat > .env <<EOF
SERVER_ID=$SERVER_ID
SERVER_NAME=$SERVER_NAME
DESCRIPTION=$DESCRIPTION
LOCATION=$LOCATION
OS=$OS
INTERVAL_TIME=$INTERVAL_TIME
EOF