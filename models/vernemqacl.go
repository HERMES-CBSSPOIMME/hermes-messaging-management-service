package models

type VerneMQACL struct {
	Mountpoint   string   `json:"mountpoint" bson:"mountpoint"`
	ClientID     string   `json:"client_id" bson:"client_id"`
	Username     string   `json:"username" bson:"username"`
	Passhash     string   `json:"passhash" bson:"passhash"`
	PublishACL   []string `json:"publish_acl" bson:"publish_acl"`
	SubscribeACL []string `json:"subscribe_acl" bson:"subscribe_acl"`
}
