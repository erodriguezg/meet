package dto

type ModelRegisterDto struct {
	PersonId string `json:"personId" bson:"personId"`
	NickName string `json:"nickName" bson:"nickName"`
}
