package service

import (
	"fmt"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/exception"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelService interface {
	SearchModels(filters domain.FilterSearchModel, first int, last int) (*domain.SearchModelResponse, error)
	FindModelByPersonId(personId string) (*domain.Model, error)
	RegisterModel(registerData dto.ModelRegisterDto) error
	FindModelByNickName(modelNickName string) (*domain.Model, error)
	PrepareUploadUrlForProfileImage(modelNickName string) ([]dto.ResourceUploadUrlDto, error)
}

type domainModelService struct {
	personService PersonService
	fileService   FileService
	repository    repository.ModelRepository
}

func NewDomainModelService(
	personService PersonService,
	fileService FileService,
	repository repository.ModelRepository) ModelService {
	return &domainModelService{personService, fileService, repository}
}

func (port *domainModelService) SearchModels(filters domain.FilterSearchModel, first int, last int) (*domain.SearchModelResponse, error) {
	totalCount, err := port.repository.SearchModelsCount(filters)
	if err != nil {
		return nil, err
	}
	var response domain.SearchModelResponse
	if totalCount > 0 {
		models, err := port.repository.SearchModels(filters, first, last)
		if err != nil {
			return nil, err
		}
		response.Models = models
	}
	response.TotalCount = totalCount
	return &response, nil
}

func (port *domainModelService) FindModelByPersonId(personId string) (*domain.Model, error) {
	return port.repository.FindModelByPersonId(personId)
}

func (port *domainModelService) RegisterModel(registerData dto.ModelRegisterDto) error {

	// check if the person already have registered a model

	existingModel, err := port.repository.FindModelByPersonId(registerData.PersonId)
	if err != nil {
		return err
	}

	if existingModel != nil {
		return exception.NewModelAlreadyRegisteredException(existingModel)
	}

	// check if model nick name is available

	existingModel, err = port.repository.FindModelByNickName(registerData.NickName)
	if err != nil {
		return err
	}

	if existingModel != nil {
		return exception.NewModelNickNameNotAvailable(registerData.NickName)
	}

	// update person

	person, err := port.personService.FindById(registerData.PersonId)
	if err != nil {
		return err
	}

	if person == nil {
		return fmt.Errorf("the person does not exist")
	}

	person.ProfileCode = domain.ProfileCodeModel

	_, err = port.personService.Save(*person)
	if err != nil {
		return err
	}

	// persist model

	personIdObjectId, err := primitive.ObjectIDFromHex(registerData.PersonId)
	if err != nil {
		return err
	}

	_, err = port.repository.SaveModel(domain.Model{
		PersonId: personIdObjectId,
		NickName: registerData.NickName,
	})

	return err
}

func (port *domainModelService) FindModelByNickName(modelNickName string) (*domain.Model, error) {
	model, err := port.repository.FindModelByNickName(modelNickName)
	if err != nil {
		return nil, fmt.Errorf("error at FindModelByNickName. Model: %s, error: %w", modelNickName, err)
	}
	return model, nil
}

func (port *domainModelService) PrepareUploadUrlForProfileImage(modelNickName string) ([]dto.ResourceUploadUrlDto, error) {

	model, err := port.mustGetModelByNickName(modelNickName)
	if err != nil {
		return nil, err
	}

	modelIdHex := model.Id.Hex()

	pathNormalFile := fmt.Sprintf("models/%s/profile-img.png", modelIdHex)
	pathThumbnailFile := fmt.Sprintf("models/%s/profile-img-thumbnail.png", modelIdHex)

	profileImageFile, normalUploadUrl, err := port.fileService.CreateForUpload(pathNormalFile, []string{modelIdHex, modelNickName, "profileImage"})
	if err != nil {
		return nil, fmt.Errorf("errot at CreateForUpload. error: %w", err)
	}

	profileImageThumbnailFile, thumbnailUploadUrl, err := port.fileService.CreateForUpload(pathThumbnailFile, []string{modelIdHex, modelNickName, "profileImageThumbnail"})
	if err != nil {
		return nil, fmt.Errorf("errot at CreateForUpload. error: %w", err)
	}

	model.ProfileImageFileHash = &profileImageFile.Hash
	model.ProfileImageThumbnailFileHash = &profileImageThumbnailFile.Hash

	_, err = port.repository.SaveModel(*model)
	if err != nil {
		return nil, fmt.Errorf("error at SaveModel for PrepareUploadUrlForProfileImage. error: %w", err)
	}

	return []dto.ResourceUploadUrlDto{
		{
			UploadUrl:   normalUploadUrl,
			FileHash:    *model.ProfileImageFileHash,
			IsThumbnail: false,
			IsBlurred:   false,
		},
		{
			UploadUrl:   thumbnailUploadUrl,
			FileHash:    *model.ProfileImageThumbnailFileHash,
			IsThumbnail: true,
			IsBlurred:   false,
		},
	}, nil

}

// private
func (port *domainModelService) mustGetModelByNickName(modelNickName string) (*domain.Model, error) {
	model, err := port.FindModelByNickName(modelNickName)
	if err != nil {
		return nil, fmt.Errorf("error at FindModelByNickName. model: %s, error: %w", modelNickName, err)
	}
	if model == nil {
		return nil, fmt.Errorf("model nickname %s does not exist", modelNickName)
	}
	return model, nil
}
