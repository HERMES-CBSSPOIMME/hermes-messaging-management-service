package models

const (
	PrivateConversationTopicPath = "conversations/private/"
	GroupConversationTopicPath   = "conversations/group/"
)

// VerneMQACL : VerneMQ ACL
type VerneMQACL struct {
	Mountpoint   string   `json:"mountpoint" bson:"mountpoint"`
	ClientID     string   `json:"clientID" bson:"client_id"`
	Username     string   `json:"username" bson:"username"`
	Passhash     string   `json:"passhash" bson:"passhash"`
	PublishACL   []string `json:"publish_acl" bson:"publish_acl"`
	SubscribeACL []string `json:"subscribe_acl" bson:"subscribe_acl"`
}

// MQTTAuthInfos : MQTT auth informations
type MQTTAuthInfos struct {
	ClientID string `json:"clientID"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewVerneMQACL : Return new VerneMQACL struct pointer
func NewVerneMQACL(clientID string, username string, password string) *VerneMQACL {
	return &VerneMQACL{
		Mountpoint:   "",
		ClientID:     clientID,
		Username:     username,
		Passhash:     password,
		SubscribeACL: []string{PrivateConversationTopicPath + clientID},
		PublishACL:   []string{PrivateConversationTopicPath + "+"},
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
