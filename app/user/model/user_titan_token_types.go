package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTitanToken struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Uid      int64              `bson:"uid,omitempty" json:"uid,omitempty"`           //用户id
	Token    string             `bson:"token,omitempty" json:"token,omitempty"`       //用户TitanToken
	Exp      int64              `bson:"exp,omitempty" json:"exp,omitempty"`           // 过期时间
	UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"` //过期时间
	CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
