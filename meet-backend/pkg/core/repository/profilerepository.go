package repository

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
)

type ProfileRepository interface {
	FindByCode(code int) (*domain.Profile, error)
	FindAll() ([]domain.Profile, error)
}
