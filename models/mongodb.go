package models

import (
	// Native Go Libs
	context "context"

	// Project Libs
	utils "hermes-messaging-service/utils"

	// 3rd Party Libs
	mongo "github.com/mongodb/mongo-go-driver/mongo"
	bson "gopkg.in/mgo.v2/bson"
)

const (

	// HermesDatabaseName : Database name of the Hermes project as defined in MongoDB
	HermesDatabaseName = "hermesDB"

	// ConversationsCollection : MongoDB Collection containing conversations backups
	ConversationsCollection = "conversations"

	// VerneMQACLCollection : MongoDB Collection containing VerneMQ ACLs
	VerneMQACLCollection = "vmq_auth_acl"
)

// DatastoreInterface : DB Communication interface
type DatastoreInterface interface {
	AddUserACL(*VerneMQACL) error
}

// MongoDB : MongoDB communication interface
type MongoDB struct {
	Client                  *mongo.Client
	HermesDB                *mongo.Database
	ConversationsCollection *mongo.Collection
	VerneMQACLCollection    *mongo.Collection
}

// NewMongoDB : Return a new MongoDB abstraction struct
func NewMongoDB(connectionURL string) *MongoDB {

	// Get connection to DB
	client, err := mongo.NewClient(connectionURL)

	if err != nil {
		utils.PanicOnError(err, "Failed to connect to MongoDB")
	}

	err = client.Connect(context.TODO())

	if err != nil {
		utils.PanicOnError(err, "Failed to connect to context")
	}

	// Get database reference
	hermesDB := client.Database(HermesDatabaseName)

	// Get collections references
	conversationsCollection := hermesDB.Collection(ConversationsCollection)
	vmqACLCollection := hermesDB.Collection(VerneMQACLCollection)

	// Return new MongoDB abstraction struct
	return &MongoDB{
		Client:                  client,
		HermesDB:                hermesDB,
		ConversationsCollection: conversationsCollection,
		VerneMQACLCollection:    vmqACLCollection,
	}
}

// AddUserACL : Add VerneMQ ACL for user in database
// Should be trigerred when a user connect for the first time
func (mongoDB *MongoDB) AddUserACL(verneMQACL *VerneMQACL) error {

	// Marshal struct into bson object
	doc, err := bson.Marshal(*verneMQACL)

	if err != nil {
		return err
	}

	// Insert ACL into VerneMQ ACL Collection
	_, err = mongoDB.VerneMQACLCollection.InsertOne(context.TODO(), doc)

	if err != nil {
		return err
	}

	return nil
}
