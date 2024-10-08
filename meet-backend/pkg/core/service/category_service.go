package service

import (
	"fmt"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryService interface {
	GetAllTree() ([]dto.CategoryDTO, error)
	Save(categoryDTO dto.CategoryDTO) (dto.CategoryDTO, error)
	Delete(categoryDTO dto.CategoryDTO) error
}

type categoryServiceImpl struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryServiceImpl(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryServiceImpl{
		categoryRepository,
	}
}

// Delete implements CategoryService.
func (port *categoryServiceImpl) Delete(categoryDTO dto.CategoryDTO) error {
	parentObjectId, err := primitive.ObjectIDFromHex(*categoryDTO.Id)
	if err != nil {
		return err
	}
	category := domain.Category{
		Id: &parentObjectId,
	}

	childrenCategories, err := port.categoryRepository.FindByParent(category)
	if err != nil {
		return err
	}

	if len(childrenCategories) > 0 {
		return fmt.Errorf("the category cannot be deleted as it has child categories")
	}

	return port.categoryRepository.Delete(category)
}

// Save implements CategoryService.
func (port *categoryServiceImpl) Save(categoryDTO dto.CategoryDTO) (dto.CategoryDTO, error) {

	var outputDto dto.CategoryDTO

	category, err := categoryDTO.ToDomain()
	if err != nil {
		return outputDto, err
	}

	if category.Id != nil {
		err = port.categoryRepository.Update(category)
		if err != nil {
			return outputDto, err
		}
	} else {
		oid, err := port.categoryRepository.Insert(category)
		if err != nil {
			return outputDto, err
		}
		category.Id = oid
	}

	err = outputDto.FromDomain(&category)
	return outputDto, err
}

// GetAll implements CategoryService.
func (port *categoryServiceImpl) GetAllTree() ([]dto.CategoryDTO, error) {

	categories, err := port.categoryRepository.FindAll()
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[string]*dto.CategoryDTO)
	for _, category := range categories {
		var dto dto.CategoryDTO
		err = dto.FromDomain(&category)
		if err != nil {
			return nil, err
		}
		if dto.Id != nil {
			dtoMap[*dto.Id] = &dto
		}
	}

	var tree []*dto.CategoryDTO
	for _, catDto := range dtoMap {
		if catDto.ParentId == nil {
			tree = append(tree, catDto)
		} else {
			parent := dtoMap[*catDto.ParentId]
			if parent != nil {
				parent.Children = append(parent.Children, catDto)
			}
		}
	}

	var output []dto.CategoryDTO
	for _, treeItem := range tree {
		output = append(output, *treeItem)
	}
	return output, nil
}
