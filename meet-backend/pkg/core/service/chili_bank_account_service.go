package service

import (
	"fmt"

	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/repository"
	stringsutils "github.com/erodriguezg/meet/pkg/util/strings_utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChiliBankAccountService interface {
	GetBanks() ([]string, error)
	GetAccountTypes() ([]string, error)
	GetAccounts(modelNickName string) ([]dto.ChiliBankAccountDTO, error)
	Save(modelNickName string, account dto.ChiliBankAccountDTO) (dto.ChiliBankAccountDTO, error)
	Delete(modelNickName string, accountId string) error
}

type domainChileBankAccountService struct {
	modelService ModelService
	repository   repository.ChiliBankAccountRepository
}

func NewDomainChiliBankAccountService(
	modelService ModelService,
	repository repository.ChiliBankAccountRepository,
) ChiliBankAccountService {
	return &domainChileBankAccountService{modelService, repository}
}

// Create implements ChiliBankAccountService.
func (port *domainChileBankAccountService) Save(modelNickName string, account dto.ChiliBankAccountDTO) (dto.ChiliBankAccountDTO, error) {

	// validate bank name
	banks, err := port.GetBanks()
	if err != nil {
		return dto.ChiliBankAccountDTO{}, err
	}
	validBank := stringsutils.StringSliceContainsValue(banks, account.BankName)
	if !validBank {
		return dto.ChiliBankAccountDTO{}, fmt.Errorf("invalid bank name: %s", account.BankName)
	}

	// validate
	accountTypes, err := port.GetAccountTypes()
	if err != nil {
		return dto.ChiliBankAccountDTO{}, err
	}
	validAccountType := stringsutils.StringSliceContainsValue(accountTypes, account.AccountType)
	if !validAccountType {
		return dto.ChiliBankAccountDTO{}, fmt.Errorf("invalid account type: %s", account.AccountType)
	}
	model, err := port.modelService.FindModelByNickName(modelNickName)
	if err != nil {
		return dto.ChiliBankAccountDTO{}, err
	}

	// persist
	accountDomain, err := dto.MapChiliBankAccountToDomain(account)
	if err != nil {
		return dto.ChiliBankAccountDTO{}, err
	}
	accountDomain.ModelId = *model.Id

	var accountObjectID *primitive.ObjectID
	if accountDomain.Id == nil {
		accountObjectID, err = port.repository.Persist(accountDomain)
	} else {
		accountObjectID, err = port.repository.Update(accountDomain)
	}

	if err != nil {
		return dto.ChiliBankAccountDTO{}, err
	}

	auxId := accountObjectID.Hex()
	account.Id = &auxId
	return account, nil
}

// Delete implements ChiliBankAccountService.
func (port *domainChileBankAccountService) Delete(modelNickName string, accountId string) error {

	accountObjectId, err := primitive.ObjectIDFromHex(accountId)
	if err != nil {
		return err
	}

	model, err := port.modelService.FindModelByNickName(modelNickName)
	if err != nil {
		return err
	}

	accountFound, err := port.repository.FindOneByIdAndModelId(accountObjectId, *model.Id)
	if err != nil {
		return err
	}

	if accountFound == nil {
		return fmt.Errorf("account not found")
	}

	return port.repository.Delete(accountObjectId)
}

// GetAccounts implements ChiliBankAccountService.
func (port *domainChileBankAccountService) GetAccounts(modelNickName string) ([]dto.ChiliBankAccountDTO, error) {
	model, err := port.modelService.FindModelByNickName(modelNickName)
	if err != nil {
		return nil, err
	}
	accounts, err := port.repository.FindByModelId(*model.Id)
	if err != nil {
		return nil, err
	}

	var dtos []dto.ChiliBankAccountDTO
	for _, account := range accounts {
		dtos = append(dtos, dto.MapChiliBankAccountToDTO(account))
	}
	return dtos, nil

}

// GetAccountTypes implements ChiliBankAccountService.
func (port *domainChileBankAccountService) GetAccountTypes() ([]string, error) {
	types := []string{
		"Cuenta Corriente",
		"Cuenta Vista",
		"Cuenta RUT",
	}
	return types, nil
}

// GetBanks implements ChiliBankAccountService.
func (port *domainChileBankAccountService) GetBanks() ([]string, error) {
	banks := []string{
		"Banco BBVA",
		"Banco BCI",
		"Banco BICE",
		"Banco Citibank",
		"Banco de Chile",
		"Banco del Desarrollo",
		"Banco Edwards",
		"Banco Estado",
		"Banco Falabella",
		"Banco Itaú Chile",
		"Banco Ripley",
		"Banco Santander Santiago",
		"Banco Scotiabank",
		"Banco Security",
		"Banco Tbanc",
		"Banco Mercado Pago",
	}
	return banks, nil
}
