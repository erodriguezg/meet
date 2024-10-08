package repository

import "github.com/erodriguezg/meet/pkg/core/domain"

type StorageRepository interface {
	GetStorageType() string

	GetFileUploadUrl(metaData domain.FileMetaData) (string, error)

	GetFileDownloadUrl(metaData domain.FileMetaData) (string, error)

	DeleteFile(metaData domain.FileMetaData) error
}
