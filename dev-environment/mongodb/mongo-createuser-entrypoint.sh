#!/usr/bin/env bash

# Constants

MONGODB_PORT=27017
MONGODB_ROOT_USER=root
MONGODB_ROOT_PASS=example

VERNEMQ_USER=vmq-user
VERNEMQ_USER_PASSWORD=example

MESSAGING_SERVICE_USER=messaging-user
MESSAGING_SERVICE_USER_PASSWORD=example

HERMES_DB_NAME=hermesDB
VERNEMQ_DB_NAME=vmqDB


echo 'Creating application user and db'

# VerneMQ system user must be able to read ACLs from vmqDB and to write conversations to hermesDB
mongo ${VERNEMQ_DB_NAME} \
        --host localhost \
        --port ${MONGODB_PORT} \
        -u ${MONGODB_ROOT_USER}  \
        -p ${MONGODB_ROOT_PASS} \
        --authenticationDatabase admin \
        --eval "db.createUser({user: '${VERNEMQ_USER}', pwd: '${VERNEMQ_USER_PASSWORD}', roles:[{role:'read', db: '${VERNEMQ_DB_NAME}'},{role:'readWrite', db: '${HERMES_DB_NAME}'}]});"

# Messaging service system user must be able to write ACLs to vmqDB and read conversations from hermesDB
mongo ${HERMES_DB_NAME} \
        --host localhost \
        --port ${MONGODB_PORT} \
        -u ${MONGODB_ROOT_USER}  \
        -p ${MONGODB_ROOT_PASS} \
        --authenticationDatabase admin \
        --eval "db.createUser({user: '${MESSAGING_SERVICE_USER}', pwd: '${MESSAGING_SERVICE_USER_PASSWORD}', roles:[{role:'read', db: '${HERMES_DB_NAME}'},{role:'readWrite', db: '${VERNEMQ_DB_NAME}'}]});"