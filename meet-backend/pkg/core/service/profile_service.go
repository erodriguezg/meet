package service

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
)

type ProfileService interface {
	FindByCode(code int) (*domain.Profile, error)

	FindAll() ([]domain.Profile, error)
}

type domainProfileService struct {
	profileRepository repository.ProfileRepository
}

func NewDomainProfileService(profileRepository repository.ProfileRepository) ProfileService {
	return &domainProfileService{profileRepository}
}

// FindByCode implements ProfileService
func (port *domainProfileService) FindByCode(code int) (*domain.Profile, error) {
	return port.profileRepository.FindByCode(code)
}

// FindAll implements ProfileService
func (port *domainProfileService) FindAll() ([]domain.Profile, error) {
	return port.profileRepository.FindAll()
}
