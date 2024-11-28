package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetFile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Uid       int64              `bson:"uid,omitempty" json:"uid,omitempty"`             // 用户ID
	Cid       string             `bson:"cid,omitempty" json:"cid,omitempty"`             // titan id
	AssetId   int64              `bson:"assetId,omitempty" json:"assetId,omitempty"`     // 资源ID
	AssetName string             `bson:"assetName,omitempty" json:"assetName,omitempty"` // 文件名
	AssetSize int64              `bson:"assetSize,omitempty" json:"assetSize,omitempty"` // 文件大小
	AssetType int64              `bson:"assetType,omitempty" json:"assetType,omitempty"` // 1文件夹2文件3视频4图片
	IsTag     int64              `bson:"isTag,omitempty" json:"isTag,omitempty"`         // 是否标记 1是2否
	Pid       int64              `bson:"pid,omitempty" json:"pid,omitempty"`             // 资源所属父资源ID
	Source    int64              `bson:"source,omitempty" json:"source,omitempty"`       // 来源 1本地上次2云上传3Telegram
	Status    int64              `bson:"status,omitempty" json:"status,omitempty"`       // 状态 1禁用2进行中3完成4失败
	Link      string             `bson:"link,omitempty" json:"link,omitempty"`           // 链接
	Path      string             `bson:"path,omitempty" json:"path,omitempty"`           // 文件路径

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
