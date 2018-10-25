#!/usr/bin/env bash

# TODO: Change credentials handling method

# Constants
PROJECT=hermes-demo

printf "\n"
echo "====================================================================================================="
echo "Removing previous containers ..."
echo "====================================================================================================="

# Stop and remove previous container
docker rm -f "${PROJECT}_mongodb"
docker rm -f "${PROJECT}_mongoexpress"

printf "\n"
echo "====================================================================================================="
echo "Building & starting containers ..."
echo "====================================================================================================="
# Build & start services
docker-compose build && docker-compose -p $PROJECT up -d
