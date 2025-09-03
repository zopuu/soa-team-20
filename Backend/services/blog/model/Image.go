package model

type Image struct {
	Data     []byte `json:"data" bson:"data, omitempty"`
	MimeType string `json:"mimeType" bson:"mimeType, omitempty"`
	Filename string `json:"filename" bson:"filename, omitempty"`
}
