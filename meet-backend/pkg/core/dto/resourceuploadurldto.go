package dto

type ResourceUploadUrlDto struct {
	UploadUrl   string `json:"uploadUrl"`
	FileHash    string `json:"fileHash"`
	IsThumbnail bool   `json:"isThumbnail"`
	IsBlurred   bool   `json:"isBlurred"`
}
