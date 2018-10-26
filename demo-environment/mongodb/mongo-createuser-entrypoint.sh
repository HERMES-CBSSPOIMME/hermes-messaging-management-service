#!/usr/bin/env bash

# Constants

MONGODB_PORT=27017
MONGODB_ROOT_USER=root
MONGODB_ROOT_PASS=example

HERMES_USER=hermes-demo-user
HERMES_USER_PASSWORD=example


HERMES_DB_NAME=hermesDemoDB

echo 'Creating application user and db'

mongo ${HERMES_DB_NAME} \
        --host localhost \
        --port ${MONGODB_PORT} \
        -u ${MONGODB_ROOT_USER}  \
        -p ${MONGODB_ROOT_PASS} \
        --authenticationDatabase admin \
        --eval "db.createUser({user: '${HERMES_USER}', pwd: '${HERMES_USER_PASSWORD}', roles:[{role:'readWrite', db: '${HERMES_DB_NAME}'}]});"
