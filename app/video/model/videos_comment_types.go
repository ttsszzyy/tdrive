package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Videos_comment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Vid         int64              `bson:"vid,omitempty" json:"vid,omitempty"`
	Uid         int64              `bson:"uid,omitempty" json:"uid,omitempty"`
	ComposeType string             `bson:"composeType,omitempty" json:"composeType,omitempty"`
	Content     string             `bson:"content,omitempty" json:"content,omitempty"`
	UpdateAt    time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt    time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
