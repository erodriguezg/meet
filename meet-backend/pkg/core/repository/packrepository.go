package repository

import "github.com/erodriguezg/meet/pkg/core/domain"

type PackRepository interface {
	FindPackById(packId string) (*domain.Pack, error)

	FindPackByModelIdAndPackNumber(modelId string, packNumber int) (*domain.Pack, error)

	FindPackActiveByModelIdAndPackNumber(modelId string, packNumber int) (*domain.Pack, error)

	FindPacksActiveByModelId(modelId string) ([]domain.Pack, error)

	SavePack(pack domain.Pack) (*domain.Pack, error)
}
