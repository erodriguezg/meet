package dto

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChiliBankAccountDTO struct {
	Id            *string `json:"id,omitempty" validate:"max=255"`
	Rut           int     `json:"rut" validate:"required,max=10"`
	HolderName    string  `json:"holderName" validate:"required,max=255"`
	BankName      string  `json:"bankName" validate:"required,max=255"`
	AccountType   string  `json:"accountType" validate:"required,max=50"`
	AccountNumber int     `json:"accountNumber" validate:"required,max=30"`
}

func MapChiliBankAccountToDTO(in domain.ChiliBankAccount) ChiliBankAccountDTO {
	var idOut *string
	if in.Id != nil {
		hexAux := in.Id.Hex()
		idOut = &hexAux
	}
	return ChiliBankAccountDTO{
		Id:            idOut,
		Rut:           in.Rut,
		HolderName:    in.HolderName,
		BankName:      in.BankName,
		AccountType:   in.AccountType,
		AccountNumber: in.AccountNumber,
	}
}

func MapChiliBankAccountToDomain(in ChiliBankAccountDTO) (domain.ChiliBankAccount, error) {

	outDomain := domain.ChiliBankAccount{
		Rut:           in.Rut,
		HolderName:    in.HolderName,
		BankName:      in.BankName,
		AccountType:   in.AccountType,
		AccountNumber: in.AccountNumber,
	}

	if in.Id != nil {
		auxId, err := primitive.ObjectIDFromHex(*in.Id)
		if err != nil {
			return domain.ChiliBankAccount{}, err
		}

		outDomain.Id = &auxId
	}

	return outDomain, nil
}
