package service

import (
	"fmt"
	"time"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"github.com/erodriguezg/meet/pkg/util/datetime"
	"github.com/erodriguezg/meet/pkg/util/hashutil"
)

type RoomService interface {
	FindAllRooms() ([]dto.RoomDTO, error)

	FindByOwnerPersonId(ownerPersonId string) ([]dto.RoomDTO, error)

	CreateRoom(params dto.CreateRoomDTO) (dto.RoomDTO, error)

	DeleteRoom(roomHash string) error

	DeleteAllExpiredRooms() error

	ChangeRoomVisibility(params dto.ChangeRoomVisibilityRoomDTO) (dto.RoomDTO, error)
}

type domainRoomService struct {
	roomRepository repository.RoomRepository
	personService  PersonService
}

func NewDomainRoomService(
	roomRepository repository.RoomRepository,
	personService PersonService) RoomService {
	return &domainRoomService{
		roomRepository,
		personService,
	}
}

// ChangeRoomVisibility implements RoomService.
func (port *domainRoomService) ChangeRoomVisibility(params dto.ChangeRoomVisibilityRoomDTO) (dto.RoomDTO, error) {
	panic("unimplemented")
}

// CreateRoom implements RoomService.
func (port *domainRoomService) CreateRoom(params dto.CreateRoomDTO) (dto.RoomDTO, error) {
	owner, err := port.findPerson(params.OwnerPersonId)
	if err != nil {
		return dto.RoomDTO{}, err
	}

	presentTime := time.Now()

	newRoom := domain.Room{
		OwnerPersonId:       *owner.Id,
		CreationDate:        datetime.NewFromTime(presentTime),
		AnonymousAccess:     params.AnonymousAccess,
		LastInteractionDate: datetime.NewFromTime(presentTime),
	}

	persistedRoom, err := port.roomRepository.Persist(newRoom)
	if err != nil {
		return dto.RoomDTO{}, nil
	}

	return port.roomToDTO(persistedRoom, &owner), nil
}

// DeleteAllExpiredRooms implements RoomService.
func (port *domainRoomService) DeleteAllExpiredRooms() error {
	panic("unimplemented")
}

// DeleteRoom implements RoomService.
func (port *domainRoomService) DeleteRoom(roomHash string) error {
	panic("unimplemented")
}

// FindAllRooms implements RoomService.
func (port *domainRoomService) FindAllRooms() ([]dto.RoomDTO, error) {
	panic("unimplemented")
}

// FindByOwnerPersonId implements RoomService.
func (port *domainRoomService) FindByOwnerPersonId(ownerPersonId string) ([]dto.RoomDTO, error) {
	panic("unimplemented")
}

// private

func (port *domainRoomService) findPerson(personId string) (domain.Person, error) {
	person, err := port.personService.FindById(personId)
	if err != nil {
		return domain.Person{}, err
	}
	if person == nil {
		return domain.Person{}, fmt.Errorf("person not found with id: %s", personId)
	}
	return *person, nil
}

func (port *domainRoomService) roomToDTO(room *domain.Room, person *domain.Person) dto.RoomDTO {

	owner := dto.RoomOwnerDTO{
		PersonId:  person.Id.Hex(),
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Email:     person.Email,
	}

	return dto.RoomDTO{
		RoomHash:            hashutil.MD5HashB64URLEncoding(room.Id.Hex()),
		Owner:               owner,
		CreationDate:        room.CreationDate,
		AnonymousAccess:     room.AnonymousAccess,
		LastInteractionDate: room.LastInteractionDate,
	}

}
