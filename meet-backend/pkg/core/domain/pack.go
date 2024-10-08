package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PackItemTypeCodeImgJpg   = "img-jpg"
	PackItemTypeCodeImgPng   = "img-png"
	PackItemTypeCodeVideoMp4 = "video-mp4"
	PackItemTypeCodeVideoOgg = "video-ogg"
)

type Pack struct {
	Id                 *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ModelId            primitive.ObjectID  `json:"modelId" bson:"modelId"`
	PackNumber         int                 `json:"packNumber" bson:"packNumber"`
	Title              *string             `json:"title,omitempty" bson:"title,omitempty"`
	Description        *string             `json:"description,omitempty" bson:"description,omitempty"`
	DollarValue        *float64            `json:"dollarValue,omitempty" bson:"dollarValue,omitempty"`
	ReadyToPublish     bool                `json:"readyToPublish" bson:"readyToPublish"`
	Published          bool                `json:"published" bson:"published"`
	CreationDate       time.Time           `json:"creationDate" bson:"creationDate"`
	ReadyToPublishDate *time.Time          `json:"readyToPublishDate,omitempty" bson:"readyToPublishDate,omitempty"`
	PublishedDate      *time.Time          `json:"publishedDate,omitempty" bson:"publishedDate,omitempty"`
	PackItems          []PackItem          `json:"packItems" bson:"packItems"`
	Active             bool                `json:"active" bson:"active"`
}

type PackItem struct {
	TypeCode                string    `json:"typeCode" bson:"typeCode"`
	ItemNumber              int       `json:"itemNumber" bson:"itemNumber"`
	ResourceFileHash        string    `json:"resourceFileHash" bson:"resourceFile"`
	ThumbnailFileHash       string    `json:"thumbnailFileHash" bson:"thumbnailFile"`
	ThumbnailLockedFileHash string    `json:"thumbnailLockedFileHash" bson:"thumbnailLockedFile"`
	PublicItem              bool      `json:"publicItem" bson:"publicItem"`
	CreationDate            time.Time `json:"creationDate" bson:"creationDate"`
	Active                  bool      `json:"active" bson:"active"`
}
