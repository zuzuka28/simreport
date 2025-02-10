#!/bin/bash

ENV_FILE=".env"
CONFIG_DIR="./configs"

if [ ! -f "$ENV_FILE" ]; then
    echo "Error: .env file not found in the current directory"
    exit 1
fi

if [ ! -d "$CONFIG_DIR" ]; then
    echo "Error: Config directory not found at $CONFIG_DIR"
    exit 1
fi

echo "Loading environment variables from $ENV_FILE..."
set -o allexport
[ -f "$ENV_FILE" ] && source "$ENV_FILE"
set +o allexport

shopt -s nullglob
for TEMPLATE_FILE in "$CONFIG_DIR"/*.template.yaml; do
    OUTPUT_FILE="${TEMPLATE_FILE%.template.yaml}.yaml"
    envsubst < "$TEMPLATE_FILE" > "$OUTPUT_FILE"
done

echo "All configurations have been generated successfully."
