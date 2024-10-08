package service

import (
	"fmt"
	"time"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"github.com/erodriguezg/meet/pkg/util/hashutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PackAccessLevelEdit   = 4
	PackAccessLevelModel  = 3
	PackAccessLevelView   = 2
	PackAccessLevelLocked = 1
	PackAccessLevelDenied = 0
)

type PackService interface {
	CreateNewPack(modelNickName string) (*dto.PackDto, error)

	DeletePack(modelNickName string, packNumber int) error

	PrepareUploadForPackItem(modelNickName string, packNumber int, typeCode string, isPublic bool) ([]dto.ResourceUploadUrlDto, error)

	GetPackInfo(modelNickName string, packNumber int, personIdRequester *string) (*dto.PackInfoDto, error)

	GetItemsFromPack(modelNickName string, packNumber int, personIdRequester *string) ([]dto.PackItemDto, error)

	GetPacksFromModel(modelNickName string, personIdRequester *string) ([]dto.PackDto, error)

	DeletePackItem(modelNickName string, packNumber int, itemNumber int) error

	ReadyToPublishPack(modelNickName string, packNumber int) error

	PublishPack(modelNickName string, packNumber int) error

	FindPackById(packId string) (*domain.Pack, error)

	FindPackByModelNicknameAndPackNumber(modelNickname string, packNumber int) (*domain.Pack, error)

	FindPackByModelIdAndPackNumber(modelId string, packNumber int) (*domain.Pack, error)

	FindActivePackByModelIdAndPackNumber(modelId string, packNumber int) (*domain.Pack, error)

	EditPackTitle(modelNickName string, packNumber int, title string) error

	EditPackDescription(modelNickName string, packNumber int, description string) error
}

type domainPackService struct {
	personService        PersonService
	profileService       ProfileService
	modelService         ModelService
	ownerResourceService OwnedResourceService
	fileService          FileService
	repository           repository.PackRepository
}

func NewDomainPackService(
	personService PersonService,
	profileService ProfileService,
	modelService ModelService,
	ownerResourceService OwnedResourceService,
	fileService FileService,
	repository repository.PackRepository,
) PackService {
	return &domainPackService{
		personService,
		profileService,
		modelService,
		ownerResourceService,
		fileService,
		repository,
	}
}

func (port *domainPackService) CreateNewPack(modelNickName string) (*dto.PackDto, error) {

	modelObjectId, err := port.getModelId(modelNickName)
	if err != nil {
		return nil, err
	}
	actualDate := time.Now()

	packs, err := port.repository.FindPacksActiveByModelId(modelObjectId.Hex())
	if err != nil {
		return nil, fmt.Errorf("error in packService FindPacksActiveByModelId. model nick: %s, error: %w", modelNickName, err)
	}

	newPackNumber := 1
	if len(packs) > 0 {
		for _, pack := range packs {
			if pack.PackNumber > newPackNumber {
				newPackNumber = pack.PackNumber
			}
		}
		newPackNumber = newPackNumber + 1
	}

	newPack := domain.Pack{
		ModelId:        *modelObjectId,
		PackNumber:     newPackNumber,
		ReadyToPublish: false,
		Published:      false,
		CreationDate:   actualDate,
		PackItems:      []domain.PackItem{},
		Active:         true,
	}

	savedPack, err := port.repository.SavePack(newPack)
	if err != nil {
		return nil, fmt.Errorf("error at PackService SavePack. newPack: %v, error: %w", newPack, err)
	}

	dto := dto.PackDto{
		PackNumber:         savedPack.PackNumber,
		Title:              savedPack.Title,
		CoverImageFileHash: nil,
		IsLocked:           false,
	}

	return &dto, nil
}

func (port *domainPackService) DeletePack(modelNickName string, packNumber int) error {

	pack, err := port.mustGetPackActive(modelNickName, packNumber)
	if err != nil {
		return fmt.Errorf("error pack service: delete pack: mustGetPack. error: %w", err)
	}

	pack.Active = false

	_, err = port.repository.SavePack(*pack)

	if err != nil {
		return fmt.Errorf("error at PackService DeletePack SavePack. error: %w", err)
	}

	return nil
}

func (port *domainPackService) DeletePackItem(modelNickName string, packNumber int, itemNumber int) error {
	pack, err := port.mustGetPackActive(modelNickName, packNumber)
	if err != nil {
		return fmt.Errorf("error pack service: DeletePackItem: mustGetPackActive. error: %w", err)
	}

	if pack.PackItems == nil {
		return fmt.Errorf("error packService: DeletePackItem: the pack don't have items. modelNick: %s, packNumber: %d", modelNickName, packNumber)
	}

	for idx := range pack.PackItems {
		packItem := &pack.PackItems[idx]
		if packItem.Active && packItem.ItemNumber == itemNumber {
			packItem.Active = false
			_, err := port.repository.SavePack(*pack)
			if err != nil {
				return fmt.Errorf("error at packService: DelePackItem: SavePack. error: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("error packService: DeletePackItem: item active not found. modelNickName: %s, packNumber: %d, itemNumber: %d", modelNickName, packNumber, itemNumber)
}

func (port *domainPackService) GetPackInfo(modelNickName string, packNumber int, personIdRequester *string) (*dto.PackInfoDto, error) {
	pack, err := port.mustGetPackActive(modelNickName, packNumber)
	if err != nil {
		return nil, err
	}
	return &dto.PackInfoDto{
		Title:       pack.Title,
		Description: pack.Description,
	}, nil
}

func (port *domainPackService) GetItemsFromPack(modelNickName string, packNumber int, personIdRequester *string) ([]dto.PackItemDto, error) {
	pack, err := port.mustGetPackActive(modelNickName, packNumber)
	if err != nil {
		return nil, fmt.Errorf("error pack service: GetItemsFromPack: mustGetPackActive. error: %w", err)
	}

	var itemsDto []dto.PackItemDto

	if len(pack.PackItems) == 0 {
		return itemsDto, nil
	}

	packAccessLevel, err := port.getAccessLevelToPack(pack, modelNickName, personIdRequester)
	if err != nil {
		return nil, err
	}

	if packAccessLevel == PackAccessLevelDenied {
		return []dto.PackItemDto{}, nil
	}

	for i := range pack.PackItems {

		packItem := pack.PackItems[i]

		var resourceFileHash *string
		var thumbnailFileHash string
		var isLocked bool

		if packItem.PublicItem || packAccessLevel >= PackAccessLevelView {
			resourceFileHash = &packItem.ResourceFileHash
			thumbnailFileHash = packItem.ThumbnailFileHash
			isLocked = false
		} else {
			resourceFileHash = nil
			thumbnailFileHash = packItem.ThumbnailLockedFileHash
			isLocked = true
		}

		itemDto := dto.PackItemDto{
			TypeCode:          packItem.TypeCode,
			ItemNumber:        packItem.ItemNumber,
			ResourceFileHash:  resourceFileHash,
			ThumbnailFileHash: thumbnailFileHash,
			IsLocked:          isLocked,
		}

		itemsDto = append(itemsDto, itemDto)
	}

	return itemsDto, nil
}

func (port *domainPackService) GetPacksFromModel(modelNickName string, personIdRequester *string) ([]dto.PackDto, error) {
	modelObjectId, err := port.getModelId(modelNickName)
	if err != nil {
		return nil, err
	}

	packs, err := port.repository.FindPacksActiveByModelId(modelObjectId.Hex())
	if err != nil {
		return nil, fmt.Errorf("error at PackService: FindPacksActiveByModelId: error: %w", err)
	}

	var packsDto []dto.PackDto

	if len(packs) == 0 {
		return packsDto, nil
	}

	for i := range packs {
		packAux := &packs[i]

		accessLevel, err := port.getAccessLevelToPack(packAux, modelNickName, personIdRequester)
		if err != nil {
			return nil, err
		}

		var isLocked bool
		// requester has no access to pack
		if accessLevel == PackAccessLevelDenied {
			continue
		} else if accessLevel == PackAccessLevelLocked {
			isLocked = true
		} else {
			isLocked = false
		}

		coverImgFileHash := port.getCoverImageFromPack(packAux, isLocked)

		dto := dto.PackDto{
			PackNumber:         packAux.PackNumber,
			Title:              packAux.Title,
			IsLocked:           isLocked,
			CoverImageFileHash: coverImgFileHash,
		}

		packsDto = append(packsDto, dto)
	}

	return packsDto, nil
}

func (port *domainPackService) PrepareUploadForPackItem(modelNickName string, packNumber int, typeCode string, isPublic bool) ([]dto.ResourceUploadUrlDto, error) {

	pack, err := port.mustGetPackActive(modelNickName, packNumber)
	if err != nil {
		return nil, err
	}

	actualDate := time.Now()
	actualDateFormat := actualDate.Format("20060102150405")

	itemNumber, err := port.getNextItemNumber(pack)
	if err != nil {
		return nil, err
	}

	var extension string
	var thumbnailExtension string
	var lockExtension string

	if typeCode == domain.PackItemTypeCodeImgJpg {
		extension = "jpg"
		thumbnailExtension = "jpg"
		lockExtension = "jpg"
	} else if typeCode == domain.PackItemTypeCodeImgPng {
		extension = "png"
		thumbnailExtension = "jpg"
		lockExtension = "jpg"
	} else if typeCode == domain.PackItemTypeCodeVideoMp4 {
		extension = "mp4"
		thumbnailExtension = "jpg"
		lockExtension = "jpg"
	} else if typeCode == domain.PackItemTypeCodeVideoOgg {
		extension = "ogg"
		thumbnailExtension = "jpg"
		lockExtension = "jpg"
	} else {
		return nil, fmt.Errorf("error at PackService: PrepareUploadForPackItem: typeCode: %s not supported", typeCode)
	}

	normalFileHash := hashutil.SHA256HashB64UrlEncoding(fmt.Sprintf("%s-%d-%d-%s-%s",
		pack.ModelId.Hex(),
		pack.PackNumber,
		itemNumber,
		actualDateFormat,
		extension,
	))

	thumbnailFileHash := hashutil.SHA256HashB64UrlEncoding(fmt.Sprintf("%s-%d-%d-%s-%s-thumbnail",
		pack.ModelId.Hex(),
		pack.PackNumber,
		itemNumber,
		actualDateFormat,
		thumbnailExtension,
	))

	lockFileHash := hashutil.SHA256HashB64UrlEncoding(fmt.Sprintf("%s-%d-%d-%s-%s-lock",
		pack.ModelId.Hex(),
		pack.PackNumber,
		itemNumber,
		actualDateFormat,
		lockExtension,
	))

	normalPath := fmt.Sprintf("models/%s/packs/pack-%d/item-%d/%s.%s",
		pack.ModelId.Hex(),
		pack.PackNumber,
		itemNumber,
		normalFileHash,
		extension)

	thumbnailPath := fmt.Sprintf("models/%s/packs/pack-%d/item-%d/%s.%s",
		pack.ModelId.Hex(),
		pack.PackNumber,
		itemNumber,
		thumbnailFileHash,
		thumbnailExtension)

	lockPath := fmt.Sprintf("models/%s/packs/pack-%d/item-%d/%s.%s",
		pack.ModelId.Hex(),
		pack.PackNumber,
		itemNumber,
		lockFileHash,
		lockExtension)

	normalFile, normalUploadUrl, err := port.fileService.CreateForUpload(normalPath,
		[]string{pack.ModelId.Hex(), pack.Id.Hex(), fmt.Sprint(itemNumber), "normal", extension, actualDateFormat})
	if err != nil {
		return nil, fmt.Errorf("error at packService: PrepareUploadForPackItem: create normal file. error: %w ", err)
	}

	thumbnailFile, thumbnailUploadUrl, err := port.fileService.CreateForUpload(thumbnailPath,
		[]string{pack.ModelId.Hex(), pack.Id.Hex(), fmt.Sprint(itemNumber), "thumbnail", thumbnailExtension, actualDateFormat})
	if err != nil {
		return nil, fmt.Errorf("error at packService: PrepareUploadForPackItem: create thumbnail file. error: %w ", err)
	}

	lockFile, lockUploadUrl, err := port.fileService.CreateForUpload(lockPath,
		[]string{pack.ModelId.Hex(), pack.Id.Hex(), fmt.Sprint(itemNumber), "lock", lockExtension, actualDateFormat})
	if err != nil {
		return nil, fmt.Errorf("error at packService: PrepareUploadForPackItem: create lock file. error: %w ", err)
	}

	newPackItem := domain.PackItem{
		TypeCode:                typeCode,
		ItemNumber:              itemNumber,
		ResourceFileHash:        normalFile.Hash,
		ThumbnailFileHash:       thumbnailFile.Hash,
		ThumbnailLockedFileHash: lockFile.Hash,
		PublicItem:              isPublic,
		CreationDate:            actualDate,
		Active:                  true,
	}

	pack.PackItems = append(pack.PackItems, newPackItem)

	_, err = port.repository.SavePack(*pack)
	if err != nil {
		return nil, fmt.Errorf("error at packService: PrepareUploadForPackItem: SavePack. error: %w", err)
	}

	return []dto.ResourceUploadUrlDto{
		{
			UploadUrl:   normalUploadUrl,
			FileHash:    normalFile.Hash,
			IsThumbnail: false,
			IsBlurred:   false,
		},
		{
			UploadUrl:   thumbnailUploadUrl,
			FileHash:    thumbnailFile.Hash,
			IsThumbnail: true,
			IsBlurred:   false,
		},
		{
			UploadUrl:   lockUploadUrl,
			FileHash:    lockFile.Hash,
			IsThumbnail: true,
			IsBlurred:   true,
		},
	}, nil

}

func (port *domainPackService) ReadyToPublishPack(modelNickName string, packNumber int) error {
	modelId, err := port.getModelId(modelNickName)
	if err != nil {
		return err
	}

	pack, err := port.repository.FindPackActiveByModelIdAndPackNumber(modelId.Hex(), packNumber)
	if err != nil {
		return err
	}

	actualTime := time.Now()

	pack.ReadyToPublish = true
	pack.ReadyToPublishDate = &actualTime

	_, err = port.repository.SavePack(*pack)
	if err != nil {
		return fmt.Errorf("error at pack service: ReadyToPublishPack: SavePack. error: %w", err)
	}

	return nil
}

func (port *domainPackService) PublishPack(modelNickName string, packNumber int) error {
	modelId, err := port.getModelId(modelNickName)
	if err != nil {
		return err
	}

	pack, err := port.repository.FindPackActiveByModelIdAndPackNumber(modelId.Hex(), packNumber)
	if err != nil {
		return err
	}

	actualTime := time.Now()

	pack.Published = true
	pack.PublishedDate = &actualTime

	_, err = port.repository.SavePack(*pack)
	if err != nil {
		return fmt.Errorf("error at pack service: ReadyToPublishPack: SavePack. error: %w", err)
	}

	return nil
}

func (port *domainPackService) FindPackById(packId string) (*domain.Pack, error) {
	return port.repository.FindPackById(packId)
}

func (port *domainPackService) FindPackByModelNicknameAndPackNumber(modelNickname string, packNumber int) (*domain.Pack, error) {
	model, err := port.modelService.FindModelByNickName(modelNickname)
	if err != nil {
		return nil, err
	}

	if model == nil {
		return nil, fmt.Errorf("model %s not found", modelNickname)
	}

	return port.FindPackByModelIdAndPackNumber(model.Id.Hex(), packNumber)
}

func (port *domainPackService) FindPackByModelIdAndPackNumber(modelId string, packNumber int) (*domain.Pack, error) {
	return port.repository.FindPackByModelIdAndPackNumber(modelId, packNumber)
}

func (port *domainPackService) FindActivePackByModelIdAndPackNumber(modelId string, packNumber int) (*domain.Pack, error) {
	return port.repository.FindPackActiveByModelIdAndPackNumber(modelId, packNumber)
}

func (port *domainPackService) EditPackDescription(modelNickName string, packNumber int, description string) error {
	pack, err := port.mustGetPackActive(modelNickName, packNumber)
	if err != nil {
		return err
	}

	pack.Description = &description

	_, err = port.repository.SavePack(*pack)
	return err
}

func (port *domainPackService) EditPackTitle(modelNickName string, packNumber int, title string) error {
	pack, err := port.mustGetPackActive(modelNickName, packNumber)
	if err != nil {
		return err
	}

	pack.Title = &title

	_, err = port.repository.SavePack(*pack)
	return err
}

// private

func (port *domainPackService) getModel(modelNickName string) (*domain.Model, error) {
	model, err := port.modelService.FindModelByNickName(modelNickName)
	if err != nil {
		return nil, fmt.Errorf("error at packService mustGetModelByNickName. model nickname: %s, error: %w", modelNickName, err)
	}
	if model == nil {
		return nil, fmt.Errorf("model not found for nickname: %s", modelNickName)
	}
	return model, nil
}

func (port *domainPackService) getModelId(modelNickName string) (*primitive.ObjectID, error) {
	model, err := port.getModel(modelNickName)
	if err != nil {
		return nil, err
	}
	return model.Id, nil
}

func (port *domainPackService) mustGetPackActive(modelNickName string, packNumber int) (*domain.Pack, error) {

	modelId, err := port.getModelId(modelNickName)
	if err != nil {
		return nil, err
	}

	pack, err := port.repository.FindPackActiveByModelIdAndPackNumber(modelId.Hex(), packNumber)
	if err != nil {
		return nil, fmt.Errorf("error at packService mustGetPackActive. modelNickName: %s, packNumber: %d, error: %w", modelNickName, packNumber, err)
	}
	if pack == nil {
		return nil, fmt.Errorf("error packService: mustGetPackActive: pack does not exist. model nickname: %s, packNumber: %d", modelNickName, packNumber)
	}
	return pack, nil
}

func (port *domainPackService) getAccessLevelToPack(pack *domain.Pack, modelNickName string, personIdRequester *string) (int, error) {

	// anonymous access
	if personIdRequester == nil && pack.Published {
		return PackAccessLevelLocked, nil
	} else if personIdRequester == nil && !pack.Published {
		return PackAccessLevelDenied, nil
	}

	// owner pack access
	personRequester, err := port.personService.FindById(*personIdRequester)
	if err != nil {
		return PackAccessLevelDenied, fmt.Errorf("error at pack service: hasAccessToPack: find person. id: %s, error: %w", *personIdRequester, err)
	}
	if personRequester == nil {
		return PackAccessLevelDenied, fmt.Errorf("error at pack service: hasAccessToPack: person requester does not exist. person id: %s", *personIdRequester)
	}
	hasPack, err := port.ownerResourceService.PersonHasPack(personRequester.Id.Hex(), pack.Id.Hex())
	if err != nil {
		return PackAccessLevelDenied, fmt.Errorf("error at pack service: hasAccessToPack: ownerResourceService: PersonHasPack, error: %w", err)
	}
	if hasPack && pack.Published {
		return PackAccessLevelView, nil
	} else if hasPack && !pack.Published {
		return PackAccessLevelDenied, nil
	}

	// model access
	model, err := port.getModel(modelNickName)
	if err != nil {
		return PackAccessLevelDenied, err
	}

	if model.Id.Hex() == pack.ModelId.Hex() && model.PersonId.Hex() == personRequester.Id.Hex() {
		return PackAccessLevelModel, nil
	}

	// permission access
	profile, err := port.profileService.FindByCode(personRequester.ProfileCode)
	if err != nil {
		return PackAccessLevelDenied, fmt.Errorf("error checking if has access to pack. error: %w", err)
	}
	if profile.PermissionsCodes != nil {
		for _, permissionCode := range profile.PermissionsCodes {
			if permissionCode == domain.PermissionCodeEditAllModels {
				return PackAccessLevelEdit, nil
			}
		}
	}
	if pack.Published {
		return PackAccessLevelLocked, nil
	} else {
		return PackAccessLevelDenied, nil
	}
}

func (port *domainPackService) getCoverImageFromPack(pack *domain.Pack, isLocked bool) *string {
	if pack == nil {
		return nil
	}

	if len(pack.PackItems) == 0 {
		return nil
	}

	firstItem := pack.PackItems[0]

	if isLocked {
		return &firstItem.ThumbnailFileHash
	}

	// if has no access
	for _, packItem := range pack.PackItems {
		// search for first public image
		if packItem.PublicItem {
			return &packItem.ThumbnailFileHash
		}
	}

	// if no have public image then return the first blurred
	return &firstItem.ThumbnailLockedFileHash
}

func (port *domainPackService) getNextItemNumber(pack *domain.Pack) (int, error) {
	if pack == nil {
		return -1, fmt.Errorf("error at packService: getNextItemNumber: pack is nil")
	}

	if len(pack.PackItems) == 0 {
		return 1, nil
	}

	higherNumber := 0
	for _, item := range pack.PackItems {
		if item.Active && item.ItemNumber > higherNumber {
			higherNumber = item.ItemNumber
		}
	}

	return higherNumber + 1, nil
}
