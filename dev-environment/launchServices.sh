#!/bin/bash

# TODO: Change credentials handling method 

# Constants
PROJECT=hermes

printf "\n"
echo "====================================================================================================="
echo "Removing previous containers ..."
echo "====================================================================================================="

# Stop and remove previous container
docker rm -f "${PROJECT}_rabbitmq"
docker rm -f "${PROJECT}_mongodb"
docker rm -f "${PROJECT}_mongoexpress"
docker rm -f "${PROJECT}_redis-cache"
docker rm -f "${PROJECT}_redis-realtime"

printf "\n"
echo "====================================================================================================="
echo "Creating containers ..."
echo "====================================================================================================="
# Start services
docker-compose -p $PROJECT up -d

# Wait for containers initialization
printf "\n"
echo "Waiting on containers to initialize ..."

sleep 10

printf "\n"
echo "====================================================================================================="
echo "RabbitMQ Setup"
echo "====================================================================================================="

printf "\n"
bash ./scripts/setup/setupRabbitMQ.sh