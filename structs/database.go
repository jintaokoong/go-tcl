package structs

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entry struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Channel         string             `bson:"channel"`
	UserID          string             `bson:"userId"`
	DisplayName     string             `bson:"displayName,omitempty"`
	Message         string             `bson:"message"`
	Roles           []string           `bson:"roles"`
	CreatedDatetime time.Time          `bson:"createdDatetime"`
}
