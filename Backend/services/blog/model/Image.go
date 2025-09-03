package model

import (
	"encoding/base64"
	"encoding/json"
)

type Image struct {
	Data     []byte `json:"data" bson:"data, omitempty"`
	MimeType string `json:"mimeType" bson:"mimeType, omitempty"`
	Filename string `json:"filename" bson:"filename, omitempty"`
}

// Custom JSON marshaling to convert byte array to base64 string
func (img Image) MarshalJSON() ([]byte, error) {
	type Alias Image
	return json.Marshal(&struct {
		Data string `json:"data,omitempty"`
		*Alias
	}{
		Data:  base64.StdEncoding.EncodeToString(img.Data),
		Alias: (*Alias)(&img),
	})
}

// Custom JSON unmarshaling to convert base64 string to byte array
func (img *Image) UnmarshalJSON(data []byte) error {
	type Alias Image
	aux := &struct {
		Data string `json:"data,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(img),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	if aux.Data != "" {
		decoded, err := base64.StdEncoding.DecodeString(aux.Data)
		if err != nil {
			return err
		}
		img.Data = decoded
	}
	
	return nil
}
