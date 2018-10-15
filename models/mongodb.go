package models

import (

	// Native Go Libs
	context "context"

	// Project Libs
	utils "hermes-messaging-service/utils"

	// 3rd Party Libs
	mongoBSON "github.com/mongodb/mongo-go-driver/bson"
	mongo "github.com/mongodb/mongo-go-driver/mongo"
	bson "gopkg.in/mgo.v2/bson"
)

const (

	// HermesDatabaseName : Database name of the Hermes project as defined in MongoDB
	HermesDatabaseName = "hermesDB"

	// PrivateConversationsCollection : MongoDB Collection containing private conversations backups
	PrivateConversationsCollection = "privateConversations"

	// VerneMQACLCollection : MongoDB Collection containing VerneMQ ACLs
	VerneMQACLCollection = "vmq_auth_acl"

	// GroupConversationCollection : MongoDB Collection containing group private conversations backups
	GroupConversationCollection = "groupConversations"
)

// MongoDBInterface : MongoDB Communication interface
type MongoDBInterface interface {
	AddGroupConversation(groupConversation *GroupConversation) error
	AddProfileACL(*VerneMQACL) error
	AuthorizePublishing(string, string) error
	UpdateProfilesWithGroupACL(*GroupConversation) error
}

// MongoDB : MongoDB communication interface
type MongoDB struct {
	Client                         *mongo.Client
	HermesDB                       *mongo.Database
	PrivateConversationsCollection *mongo.Collection
	VerneMQACLCollection           *mongo.Collection
	GroupConversationCollection    *mongo.Collection
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
	privateConversationsCollection := hermesDB.Collection(PrivateConversationsCollection)
	vmqACLCollection := hermesDB.Collection(VerneMQACLCollection)
	groupConversationCollection := hermesDB.Collection(GroupConversationCollection)

	// Return new MongoDB abstraction struct
	return &MongoDB{
		Client:                         client,
		HermesDB:                       hermesDB,
		PrivateConversationsCollection: privateConversationsCollection,
		VerneMQACLCollection:           vmqACLCollection,
		GroupConversationCollection:    groupConversationCollection,
	}
}

// AddGroupConversation : Add group conversation entry in database
func (mongoDB *MongoDB) AddGroupConversation(groupConversation *GroupConversation) error {

	// Marshal struct into bson object
	doc, err := bson.Marshal(*groupConversation)

	if err != nil {
		return err
	}

	// Insert group conversation into DB
	_, err = mongoDB.GroupConversationCollection.InsertOne(nil, doc)

	if err != nil {
		return err
	}

	return nil
}

// AddProfileACL : Add VerneMQ ACL for user in database
// Should be trigerred when a user connect for the first time
func (mongoDB *MongoDB) AddProfileACL(verneMQACL *VerneMQACL) error {

	// Marshal struct into bson object
	doc, err := bson.Marshal(*verneMQACL)

	if err != nil {
		return err
	}

	// Insert ACL into VerneMQ ACL Collection
	_, err = mongoDB.VerneMQACLCollection.InsertOne(nil, doc)

	if err != nil {
		return err
	}

	return nil
}

// AuthorizePublishing : Authorize publishing on MQTT topic for userID
func (mongoDB *MongoDB) AuthorizePublishing(userID string, topic string) error {

	_, err := mongoDB.VerneMQACLCollection.UpdateOne(
		nil,
		mongoBSON.NewDocument(
			mongoBSON.EC.String("client_id", userID),
		),
		mongoBSON.NewDocument(
			mongoBSON.EC.SubDocumentFromElements("$push",
				mongoBSON.EC.String("publish_acl", topic),
			),
		),
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateProfilesWithGroupACL : Update VerneMQ Acls in database to grant publish and read access to all members of the group
func (mongoDB *MongoDB) UpdateProfilesWithGroupACL(groupConversation *GroupConversation) error {

	for _, userID := range groupConversation.Members {

		_, err := mongoDB.VerneMQACLCollection.UpdateOne(
			nil,
			mongoBSON.NewDocument(
				mongoBSON.EC.String("client_id", userID),
			),
			mongoBSON.NewDocument(
				mongoBSON.EC.SubDocumentFromElements("$push",
					mongoBSON.EC.String("publish_acl", GroupConversationTopicPath+groupConversation.GroupConversationID),
				),
				mongoBSON.EC.SubDocumentFromElements("$push",
					mongoBSON.EC.String("subscribe_acl", GroupConversationTopicPath+groupConversation.GroupConversationID),
				),
			),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
