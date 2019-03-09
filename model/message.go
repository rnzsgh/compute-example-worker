package model

type JobMessage struct {
	ObjectKey string `json:"object"`
	Type      string `json:"type"`
	RequestId string `json:"requestId"`
}
