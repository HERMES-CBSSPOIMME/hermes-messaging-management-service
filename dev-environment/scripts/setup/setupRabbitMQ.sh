# Constants
ROOT_USER=root
ROOT_PASS=example

MESSAGING_SERVICE_USER=messaging-service-user
MESSAGING_SERVICE_PASSWORD=example

BACKUP_VHOST_NAME=backup
MESSAGING_VHOST_NAME=messaging

CONVERSATIONS_BACKUP_EXCHANGE_NAME=conversations.backup.dx
CONVERSATIONS_PRIVATE_EXCHANGE_NAME=conversations.private.tx

BACKUP_QUEUE_NAME=messages-to-backup-q

#======================================================================================================
# RabbitMQ Setup
#======================================================================================================

# Add VHOSTs
# =====================================================================================================

# backup VHOST is where messages sent by users are received in order to be stored in DB
docker exec hermes_rabbitmq rabbitmqctl add_vhost $BACKUP_VHOST_NAME

# messaging VHOST is where messages sent by users are routed
docker exec hermes_rabbitmq rabbitmqctl add_vhost $MESSAGING_VHOST_NAME


# Add users
# =====================================================================================================

# messaging-service-user is the user used by the messaging microservice to perform operations on RabbitMQ
docker exec hermes_rabbitmq rabbitmqctl add_user $MESSAGING_SERVICE_USER $MESSAGING_SERVICE_PASSWORD


# Set users permissions
# =====================================================================================================

# Set permissions for root on all VHOSTs

# Set FULL permissions on ALL resources on backup VHOST for user root
docker exec hermes_rabbitmq rabbitmqctl set_permissions -p $BACKUP_VHOST_NAME $ROOT_USER ".*" ".*" ".*"

# Set FULL permissions on ALL resources on messaging VHOST for user root
docker exec hermes_rabbitmq rabbitmqctl set_permissions -p $MESSAGING_VHOST_NAME $ROOT_USER ".*" ".*" ".*"

# messaging-service-user 

# Should be able to :
# - Read content of messages-to-backup-queue (READ rights on messages-to-backup + READ rights on backup VHOST)
# - Create group conversations exchanges on VHOST messaging (CONFIGURE rights on exchange + CONFIGURE rights on  messaging VHOST )

# +--------------------------------------------+
# | Virtual Hosts | Configure  | Write | Read  |
# |:-------------:|:----------:|:-----:|:-----:|
# | (1) backup    |      N     |   N   |   Y   |
# | (2) messaging |      Y     |   N   |   N   |
# +--------------------------------------------+

# +--------------------------------------------------------------------+
# |         Exchanges                     | Configure  | Write | Read  |
# |:-------------------------------------:|:----------:|:-----:|:-----:|
# | (3) conversations.backup.dx           |      N     |   N   |   N   |
# | (4) conversations.private.tx          |      N     |   N   |   N   |
# | (5) conversations.group.{groupID}.fx  |      Y     |   N   |   N   |
# +--------------------------------------------------------------------+

# +-----------------------------------------------------------+
# |           Queues             | Configure  | Write | Read  |
# |:----------------------------:|:----------:|:-----:|:-----:|
# | (6) messages-to-backup-q     |      N     |   N   |   Y   |
# | (7) messages-user-{userID}-q |      N     |   N   |   N   |
# +-----------------------------------------------------------+

# Set READ permissions on messages-to-backup-q (6) on backup VHOST (1) for user messaging-service-user
docker exec hermes_rabbitmq rabbitmqctl set_permissions -p $BACKUP_VHOST_NAME $MESSAGING_SERVICE_USER "^$" "^$" "^$BACKUP_QUEUE_NAME$"

# Set CONFIGURE permissions (5) on conversations.group.*.fx on backup VHOST (2) for user messaging-service-user
docker exec hermes_rabbitmq rabbitmqctl set_permissions -p $MESSAGING_VHOST_NAME $MESSAGING_SERVICE_USER "conversations.group.*.fx" "^$" "^$"


# Declare fixed exchanges
# =====================================================================================================

# Declare the direct exchange << conversations.backup.dx >> responsible to route incoming messages to backup microservice
docker exec hermes_rabbitmq rabbitmqadmin declare exchange --vhost=$BACKUP_VHOST_NAME name=$CONVERSATIONS_BACKUP_EXCHANGE_NAME type=direct durable=true -u $ROOT_USER -p $ROOT_PASS

# Declare the topic exchange << conversations.private.tx >> responsible to route incoming messages to private conversations
docker exec hermes_rabbitmq rabbitmqadmin declare exchange --vhost=$MESSAGING_VHOST_NAME name=$CONVERSATIONS_PRIVATE_EXCHANGE_NAME type=direct durable=true -u $ROOT_USER -p $ROOT_PASS

# Declare fixed queues
# =====================================================================================================

# Backup queue declaration and binding to << conversations.backup.dx >> 
# This queue is responsible to route contain incoming messages to be stored in database by the backup microservice
docker exec hermes_rabbitmq rabbitmqadmin declare queue --vhost=$BACKUP_VHOST_NAME name=$BACKUP_QUEUE_NAME auto_delete=false durable=true -u $ROOT_USER -p $ROOT_PASS
docker exec hermes_rabbitmq rabbitmqadmin declare binding --vhost=$BACKUP_VHOST_NAME source=$CONVERSATIONS_BACKUP_EXCHANGE_NAME destination=$BACKUP_QUEUE_NAME -u $ROOT_USER -p $ROOT_PASS


# Echo install informations 
# =====================================================================================================

# List exchanges
printf "\n"
echo "RabbitMQ Install Informations"
echo "====================================================================================================="

printf "\n"
echo "Overview"
docker exec hermes_rabbitmq rabbitmqadmin show overview -u $ROOT_USER -p $ROOT_PASS

printf "\n"
echo "Users"
docker exec hermes_rabbitmq rabbitmqadmin list users -u $ROOT_USER -p $ROOT_PASS

printf "\n"
echo "Permissions"
docker exec hermes_rabbitmq rabbitmqadmin list permissions user vhost configure read write -u $ROOT_USER -p $ROOT_PASS

printf "\n"
echo "Exchanges"
docker exec hermes_rabbitmq rabbitmqadmin list exchanges name vhost node type durable auto_delete internal -u $ROOT_USER -p $ROOT_PASS

printf "\n"
echo "Queues"
docker exec hermes_rabbitmq rabbitmqadmin list queues name vhost node messages message_stats.publish_details.rate -u $ROOT_USER -p $ROOT_PASS

