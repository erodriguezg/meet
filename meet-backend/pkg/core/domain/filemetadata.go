package domain

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileMetaData struct {
	Id          *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Hash        string              `json:"hash" bson:"hash"`
	Path        string              `json:"path" bson:"path"`
	Uploaded    bool                `json:"uploaded" bson:"uploaded"`
	DownloadUrl *string             `json:"downloadUrl,omitempty" bson:"downloadUrl,omitempty"`
}

func (meta *FileMetaData) GetFileName() string {
	splitPaths := meta.splitPath()
	lenSplitPaths := len(splitPaths)
	if lenSplitPaths < 1 {
		return ""
	}
	return splitPaths[lenSplitPaths-1]
}

func (meta *FileMetaData) GetFolders() []string {
	splitPaths := meta.splitPath()
	lenSplitPaths := len(splitPaths)
	if lenSplitPaths < 2 {
		return []string{}
	}
	return splitPaths[0 : lenSplitPaths-1]
}

func (meta *FileMetaData) GetFoldersPath() string {
	var folders = meta.GetFolders()
	var output string
	first := true
	for _, folder := range folders {
		if first {
			output = folder
			first = false
		} else {
			output = output + "/" + folder
		}
	}
	return output
}

func (meta *FileMetaData) GetExtension() string {
	return ""
}

func (meta *FileMetaData) GetContentType() string {
	return ""
}

// privates

func (meta *FileMetaData) splitPath() []string {
	return strings.Split(meta.Path, "/")
}
