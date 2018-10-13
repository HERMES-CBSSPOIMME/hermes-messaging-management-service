package models

import (
	context "context"
	"fmt"
	utils "hermes-messaging-service/utils"

	mongo "github.com/mongodb/mongo-go-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	HermesDatabaseName  = "hermesDB"
	VerneMQDatabaseName = "vmqDB"

	ConversationsCollection = "conversations"
	VerneMQACLCollection    = "vmq_auth_acl"
)

type DatastoreInterface interface {
	AddUserACL(clientID string, username string, password string) error
}

type MongoDB struct {
	Client                  *mongo.Client
	HermesDB                *mongo.Database
	ConversationsCollection *mongo.Collection
	VerneMQACLCollection    *mongo.Collection
}

// NewMongoDB : Return a new MongoDB communication interface
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

	hermesDB := client.Database(HermesDatabaseName)

	conversationsCollection := hermesDB.Collection(ConversationsCollection)
	vmqACLCollection := hermesDB.Collection(VerneMQACLCollection)

	return &MongoDB{
		Client:                  client,
		HermesDB:                hermesDB,
		ConversationsCollection: conversationsCollection,
		VerneMQACLCollection:    vmqACLCollection,
	}
}

func (mongoDB *MongoDB) AddUserACL(clientID string, username string, password string) error {

	acl := VerneMQACL{
		Mountpoint: "",
		ClientID:   clientID,
		Username:   username,
		Passhash:   password,
	}
	doc, err := bson.Marshal(acl)

	if err != nil {
		return err
	}

	res, err := mongoDB.VerneMQACLCollection.InsertOne(context.TODO(), doc)

	if err != nil {
		return err
	}

	fmt.Println(res.InsertedID)

	return nil
}
