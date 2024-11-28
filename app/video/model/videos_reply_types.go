package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Videos_reply struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	CommentId  string `bson:"comment_id,omitempty" json:"commentId,omitempty"`
	ToUserid   int64  `bson:"to_userid,omitempty" json:"toUserid,omitempty"`
	FromUserid int64  `bson:"from_userid,omitempty" json:"fromUserid,omitempty"`
	ReplyType  string `bson:"reply_type,omitempty" json:"replyType,omitempty"`
	Content    string `bson:"content,omitempty" json:"content,omitempty"`

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
