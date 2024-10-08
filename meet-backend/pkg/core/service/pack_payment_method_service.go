package service

import (
	"fmt"
	"time"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackPaymentMethodService interface {
	Save(modelNickName string, packNumber int, packPaymentMethod dto.PackPaymentMethodDTO) (dto.PackPaymentMethodDTO, error)
	GetFromPack(modelNickName string, packNumber int) (*dto.PackPaymentMethodDTO, error)
}

type domainPPMSService struct {
	packService PackService
	repository  repository.PackPaymentMethodRepository
}

func NewPackPaymentMethodService(packService PackService,
	repository repository.PackPaymentMethodRepository) PackPaymentMethodService {
	return &domainPPMSService{
		packService,
		repository,
	}
}

func (port *domainPPMSService) Save(modelNickName string, packNumber int, packPaymentMethodDTO dto.PackPaymentMethodDTO) (dto.PackPaymentMethodDTO, error) {

	pack, err := port.packService.FindPackByModelNicknameAndPackNumber(modelNickName, packNumber)
	if err != nil {
		return dto.PackPaymentMethodDTO{}, err
	}
	if pack == nil {
		return dto.PackPaymentMethodDTO{}, fmt.Errorf("no pack found for modelNickName: %s and packNumber: %d", modelNickName, packNumber)
	}

	existingPackPayment, err := port.repository.FindByPackId(*pack.Id)
	if err != nil {
		return dto.PackPaymentMethodDTO{}, err
	}

	var id *primitive.ObjectID
	if existingPackPayment != nil {
		id = existingPackPayment.Id
	} else {
		id = nil
	}

	updatedTime := time.Now()

	var chiliBankReceiptAccountId *primitive.ObjectID
	if packPaymentMethodDTO.ChiliBankReceiptAccountId != nil {
		auxObjectId, err := primitive.ObjectIDFromHex(*packPaymentMethodDTO.ChiliBankReceiptAccountId)
		if err != nil {
			return dto.PackPaymentMethodDTO{}, err
		}
		chiliBankReceiptAccountId = &auxObjectId
	} else {
		chiliBankReceiptAccountId = nil
	}

	toUpdateDomain := domain.PackPaymentMethod{
		Id:                            id,
		ChiliBankReceiptMethodEnabled: packPaymentMethodDTO.ChiliBankReceiptMethodEnabled,
		PackId:                        *pack.Id,
		ChiliBankReceiptAccountId:     chiliBankReceiptAccountId,
		ChiliBankReceiptCLPPrice:      packPaymentMethodDTO.ChiliBankReceiptCLPPrice,
		PaypalReceiptMethodEnabled:    packPaymentMethodDTO.PaypalReceiptMethodEnabled,
		PaypalReceiptRecipientEmail:   packPaymentMethodDTO.PaypalReceiptRecipientEmail,
		PaypalReceiptUSDPrice:         packPaymentMethodDTO.PaypalOnlineUSDPrice,
		PaypalOnlineMethodEnabled:     packPaymentMethodDTO.PaypalOnlineMethodEnabled,
		PaypalOnlineRecipientEmail:    packPaymentMethodDTO.PaypalOnlineRecipientEmail,
		PaypalOnlineUSDPrice:          packPaymentMethodDTO.PaypalOnlineUSDPrice,
		UpdateDate:                    updatedTime,
	}

	savedDomain, err := port.repository.Save(toUpdateDomain)
	if err != nil {
		return dto.PackPaymentMethodDTO{}, err
	}

	newDto := port.mapDomainToDTO(&savedDomain)
	return newDto, nil

}

func (port *domainPPMSService) GetFromPack(modelNickName string, packNumber int) (*dto.PackPaymentMethodDTO, error) {

	pack, err := port.packService.FindPackByModelNicknameAndPackNumber(modelNickName, packNumber)
	if err != nil {
		return nil, err
	}
	if pack == nil {
		return nil, fmt.Errorf("no pack found for modelNickName: %s and packNumber: %d", modelNickName, packNumber)
	}

	packPaymentMethodFound, err := port.repository.FindByPackId(*pack.Id)
	if err != nil {
		return nil, err
	}
	if packPaymentMethodFound == nil {
		return nil, nil
	} else {
		dto := port.mapDomainToDTO(packPaymentMethodFound)
		return &dto, nil
	}
}

// private

func (port *domainPPMSService) mapDomainToDTO(packPaymentMethod *domain.PackPaymentMethod) dto.PackPaymentMethodDTO {

	var chiliBankReceiptAccountId *string
	if packPaymentMethod.ChiliBankReceiptAccountId != nil {
		auxText := packPaymentMethod.ChiliBankReceiptAccountId.Hex()
		chiliBankReceiptAccountId = &auxText
	}

	dto := dto.PackPaymentMethodDTO{
		ChiliBankReceiptMethodEnabled: packPaymentMethod.ChiliBankReceiptMethodEnabled,
		ChiliBankReceiptAccountId:     chiliBankReceiptAccountId,
		ChiliBankReceiptCLPPrice:      packPaymentMethod.ChiliBankReceiptCLPPrice,
		PaypalReceiptMethodEnabled:    packPaymentMethod.PaypalReceiptMethodEnabled,
		PaypalReceiptRecipientEmail:   packPaymentMethod.PaypalReceiptRecipientEmail,
		PaypalReceiptUSDPrice:         packPaymentMethod.PaypalReceiptUSDPrice,
		PaypalOnlineMethodEnabled:     packPaymentMethod.PaypalOnlineMethodEnabled,
		PaypalOnlineRecipientEmail:    packPaymentMethod.PaypalOnlineRecipientEmail,
		PaypalOnlineUSDPrice:          packPaymentMethod.PaypalOnlineUSDPrice,
	}

	return dto
}
