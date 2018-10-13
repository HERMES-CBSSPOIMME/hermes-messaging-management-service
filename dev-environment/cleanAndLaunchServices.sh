#!/usr/bin/env bash

# TODO: Change credentials handling method 

# Constants
PROJECT=hermes

printf "\n"
echo "====================================================================================================="
echo "Removing previous containers ..."
echo "====================================================================================================="

# Stop and remove previous container
docker rm -f "${PROJECT}_vernemq"
docker rm -f "${PROJECT}_mongodb"
docker rm -f "${PROJECT}_mongoexpress"
docker rm -f "${PROJECT}_redis-sessions-cache"
docker rm -f "${PROJECT}_redis-real-time"

printf "\n"
echo "====================================================================================================="
echo "Building & starting containers ..."
echo "====================================================================================================="
# Build & start services
docker-compose build && docker-compose -p $PROJECT up -d
