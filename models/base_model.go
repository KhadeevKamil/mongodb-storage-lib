package models

import (
	timecop "github.com/bluele/go-timecop"
	"github.com/rekamarket/mongodb-storage-lib/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel
type BaseModel struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at" bson:"updated_at"`
}

// GetID() returns ID of the document
func (bs *BaseModel) GetID() primitive.ObjectID {
	return bs.ID
}

// GetHexID() returns ID of the document
func (bs *BaseModel) GetHexId() string {
	return bs.ID.Hex()
}

// SetHexID
func (bs *BaseModel) SetHexID(hexID string) error {
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return helpers.ErrInvalidObjectID
	}

	bs.ID = oid

	return nil
}

// SetupTimestamps() sets CreatedAt and UpdatedAt fields for the model
// The method does not store any data to database
// you should use the method before InsertMany(), UpdateMany() requests from you storage
func (m *BaseModel) SetupTimestamps() {
	if m.CreatedAt == 0 {
		m.CreatedAt = timecop.Now().Unix()
	}

	m.UpdatedAt = timecop.Now().Unix()
}
