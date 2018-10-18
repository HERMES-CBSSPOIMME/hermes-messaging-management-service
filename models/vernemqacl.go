package models

const (
	PrivateConversationTopicPath = "conversations/private/"
	GroupConversationTopicPath   = "conversations/group/"
)

// VerneMQACL : VerneMQ ACL
type VerneMQACL struct {
	Mountpoint   string `json:"mountpoint" bson:"mountpoint"`
	ClientID     string `json:"clientID" bson:"client_id"`
	Username     string `json:"username" bson:"username"`
	Passhash     string `json:"passhash" bson:"passhash"`
	PublishACL   []*ACL `json:"publish_acl" bson:"publish_acl"`
	SubscribeACL []*ACL `json:"subscribe_acl" bson:"subscribe_acl"`
}

// MQTTAuthInfos : MQTT auth informations
type MQTTAuthInfos struct {
	ClientID string `json:"clientID"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ACL : ACL entry
type ACL struct {
	Pattern string `json:"pattern" bson:"pattern"`
}

// NewVerneMQACL : Return new VerneMQACL struct pointer
func NewVerneMQACL(clientID string, username string, password string) *VerneMQACL {

	pubPrivateACL := ACL{Pattern: PrivateConversationTopicPath + clientID + "/+"}
	subPrivateACL := ACL{Pattern: PrivateConversationTopicPath + "+/" + clientID}

	subACLs := []*ACL{&subPrivateACL}
	pubACLs := []*ACL{&pubPrivateACL}

	return &VerneMQACL{
		Mountpoint:   "",
		ClientID:     clientID,
		Username:     username,
		Passhash:     password,
		SubscribeACL: subACLs,
		PublishACL:   pubACLs,
	}
}

// NewMQTTAuthInfos : Return new NewMQTTAuthInfos struct pointer
func NewMQTTAuthInfos(clientID string, token string) *MQTTAuthInfos {

	return &MQTTAuthInfos{
		ClientID: clientID,
		Username: clientID,
		Password: token,
	}
}
