package service

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
)

type StorageService interface {
	GetStorageType() string

	GetFileUploadUrl(metaData domain.FileMetaData) (string, error)

	GetFileDownloadUrl(metaData domain.FileMetaData) (string, error)

	DeleteFile(metaData domain.FileMetaData) error
}

type domainStorageService struct {
	storageRepository repository.StorageRepository
}

func NewDomainStorageService(storageRepository repository.StorageRepository) StorageService {
	return &domainStorageService{storageRepository}
}

func (port *domainStorageService) GetStorageType() string {
	return port.storageRepository.GetStorageType()
}

func (port *domainStorageService) DeleteFile(metaData domain.FileMetaData) error {
	return port.storageRepository.DeleteFile(metaData)
}

func (port *domainStorageService) GetFileDownloadUrl(metaData domain.FileMetaData) (string, error) {
	return port.storageRepository.GetFileDownloadUrl(metaData)
}

func (port *domainStorageService) GetFileUploadUrl(metaData domain.FileMetaData) (string, error) {
	return port.storageRepository.GetFileUploadUrl(metaData)
}
