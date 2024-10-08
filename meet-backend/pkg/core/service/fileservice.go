package service

import (
	"fmt"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"github.com/erodriguezg/meet/pkg/util/hashutil"
)

type FileService interface {
	GetStorageType() string

	FindByHash(hash string) (*domain.FileMetaData, error)

	GetDownloadUrl(hash string) (string, error)

	CreateForUpload(path string, hashSeeds []string) (*domain.FileMetaData, string, error)

	ConfirmUploaded(hash string) error

	Delete(hash string) error
}

type domainFileService struct {
	storageService   StorageService
	fileMetaDataRepo repository.FileMetaDataRepository
}

func NewFileService(storageService StorageService, fileMetaDataRepo repository.FileMetaDataRepository) FileService {
	return &domainFileService{storageService, fileMetaDataRepo}
}

func (port *domainFileService) GetStorageType() string {
	return port.storageService.GetStorageType()
}

func (port *domainFileService) ConfirmUploaded(hash string) error {

	fileMetaData, err := port.mustGetOneByHash(hash)
	if err != nil {
		return err
	}

	fileMetaData.Uploaded = true

	_, err = port.fileMetaDataRepo.Save(*fileMetaData)
	if err != nil {
		return fmt.Errorf("error at fileMetaDataRepo.Save, error: %w", err)
	}

	return nil
}

func (port *domainFileService) CreateForUpload(path string, hashSeeds []string) (*domain.FileMetaData, string, error) {

	hashPlain := ""
	for i, seed := range hashSeeds {
		if i == 0 {
			hashPlain = seed
		} else {
			hashPlain = hashPlain + "-" + seed
		}
	}

	hashBcrypt, err := hashutil.BCryptHash(hashPlain)
	if err != nil {
		return nil, "", fmt.Errorf("errot at generating BCryptHash for seed: %s, error: %w", hashPlain, err)
	}

	b64UrlEncodedHash := hashutil.B64UrlEncoding(hashBcrypt)

	fileMetaData, err := port.fileMetaDataRepo.Save(domain.FileMetaData{
		Hash:     b64UrlEncodedHash,
		Path:     path,
		Uploaded: false,
	})
	if err != nil {
		return nil, "", fmt.Errorf("error at fileMetaDataRepo.Save. error: %w", err)
	}

	uploadUrl, err := port.storageService.GetFileUploadUrl(*fileMetaData)
	if err != nil {
		return nil, "", fmt.Errorf("errot at storageService.GetFileUploadUrl. err: %w", err)
	}

	return fileMetaData, uploadUrl, nil
}

func (port *domainFileService) Delete(hash string) error {

	fileMetaData, err := port.mustGetOneByHash(hash)
	if err != nil {
		return err
	}

	err = port.storageService.DeleteFile(*fileMetaData)
	if err != nil {
		return fmt.Errorf("error at storageService.DeleteFile. error: %w", err)
	}

	err = port.fileMetaDataRepo.Delete(*fileMetaData.Id)
	if err != nil {
		return fmt.Errorf("error at fileMetaDataRepo.Delete. error: %w", err)
	}

	return nil
}

func (port *domainFileService) FindByHash(hash string) (*domain.FileMetaData, error) {
	return port.fileMetaDataRepo.FindByHash(hash)
}

func (port *domainFileService) GetDownloadUrl(hash string) (string, error) {
	fileMetaData, err := port.mustGetOneByHash(hash)
	if err != nil {
		return "", err
	}

	if fileMetaData.DownloadUrl != nil {
		return *fileMetaData.DownloadUrl, nil
	}

	downloadUrl, err := port.storageService.GetFileDownloadUrl(*fileMetaData)
	if err != nil {
		return "", fmt.Errorf("error at storageService.GetFileDownloadUrl. error: %w", err)
	}

	fileMetaData.DownloadUrl = &downloadUrl

	_, err = port.fileMetaDataRepo.Save(*fileMetaData)
	if err != nil {
		return "", fmt.Errorf("error at fileMetaDataRepo.Save in domainFileService.GetDownloadUrl. err: %w", err)
	}

	return downloadUrl, nil
}

// private

func (port *domainFileService) mustGetOneByHash(hash string) (*domain.FileMetaData, error) {
	fileMetaData, err := port.fileMetaDataRepo.FindByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("error at fileMetaDataRepo.FindByHash. error: %w", err)
	}
	if fileMetaData == nil {
		return nil, fmt.Errorf("fileMetaData required but not found for hash: %s", hash)
	}
	return fileMetaData, nil
}
