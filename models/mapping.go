package models

// Mapping : Mapping between external and minternal user ID
type Mapping struct {
	OriginalUserID     string `json:"originalUserID"`
	InternalWaveUserID string `json:"internalWaveUserID"`
}
