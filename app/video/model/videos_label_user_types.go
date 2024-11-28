package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Videos_label_user struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Lid      string             `bson:"lid,omitempty" json:"lid,omitempty"`
	Uid      int64              `bson:"uid,omitempty" json:"uid,omitempty"`
	Likes    bool               `bson:"likes" json:"likes,omitempty"`
	NoLikes  bool               `bson:"no_likes" json:"noLikes,omitempty"`
	UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
