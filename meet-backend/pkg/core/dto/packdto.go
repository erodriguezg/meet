package dto

type PackDto struct {
	PackNumber         int     `json:"packNumber"`
	Title              *string `json:"title,omitempty"`
	CoverImageFileHash *string `json:"coverImageFileHash,omitempty"`
	IsLocked           bool    `json:"isLocked"`
}

type PackItemDto struct {
	TypeCode          string  `json:"typeCode"`
	ItemNumber        int     `json:"itemNumber"`
	ResourceFileHash  *string `json:"resourceFileHash,omitempty"`
	ThumbnailFileHash string  `json:"thumbnailFileHash"`
	IsLocked          bool    `json:"isLocked"`
}
