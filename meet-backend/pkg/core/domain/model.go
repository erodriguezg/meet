package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Model struct {
	Id                            *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PersonId                      primitive.ObjectID  `json:"personId" bson:"personId"`
	NickName                      string              `json:"nickName" bson:"nickName"`
	ProfileImageFileHash          *string             `json:"profileImageFileHash,omitempty" bson:"profileImageFileHash,omitempty"`
	ProfileImageThumbnailFileHash *string             `json:"profileImageThumbnailFileHash,omitempty" bson:"profileImageThumbnailFileHash,omitempty"`
	AboutMe                       *string             `json:"aboutMe,omitempty" bson:"aboutMe,omitempty"`
	CountryCode                   *string             `json:"countryCode,omitempty" bson:"countryCode,omitempty"`
	City                          *string             `json:"city,omitempty" bson:"city,omitempty"`
	ZodiacSignCode                *string             `json:"zodiacSignCode,omitempty" bson:"ZodiacSignCode,omitempty"`
}

type FilterSearchModel struct {
	NickName       *string `json:"nickName"`
	CountryCode    *string `json:"countryCode"`
	CityName       *string `json:"cityName"`
	ZodiacSignCode *string `json:"zodiacSignCode"`
}

type SearchModelResponse struct {
	TotalCount int     `json:"totalCount"`
	Models     []Model `json:"models"`
}
