package repository

import "github.com/erodriguezg/meet/pkg/core/domain"

type ModelRepository interface {
	SearchModelsCount(filters domain.FilterSearchModel) (int, error)
	SearchModels(filters domain.FilterSearchModel, first int, last int) ([]domain.Model, error)
	SaveModel(domain.Model) (*domain.Model, error)
	FindModelByPersonId(personId string) (*domain.Model, error)
	FindModelByNickName(nickName string) (*domain.Model, error)
}
