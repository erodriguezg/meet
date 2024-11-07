package service

import (
	"fmt"
	"time"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"github.com/erodriguezg/meet/pkg/util/datetime"
	"github.com/erodriguezg/meet/pkg/util/hashutil"
	"github.com/erodriguezg/meet/pkg/util/sliceutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomService interface {
	FindAllRooms() ([]dto.RoomDTO, error)

	FindByOwnerPersonId(ownerPersonId string) ([]dto.RoomDTO, error)

	CreateRoom(params dto.CreateRoomDTO) (dto.RoomDTO, error)

	DeleteRoom(roomHash string) error

	DeleteOwnRoom(roomHash string, ownerPersonId string) error

	DeleteAllExpiredRooms() error

	ChangeRoomVisibility(params dto.ChangeRoomVisibilityRoomDTO) (dto.RoomDTO, error)

	ChangeRoomVisibilityOwnRoom(params dto.ChangeRoomVisibilityRoomDTO, ownerPersonId string) (dto.RoomDTO, error)
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

	roomFound, err := port.findRoomByHash(params.RoomHash)
	if err != nil {
		return dto.RoomDTO{}, err
	}

	owner, err := port.findPerson(roomFound.OwnerPersonId.Hex())
	if err != nil {
		return dto.RoomDTO{}, err
	}

	return port.changeRoomVisibility(params, &roomFound, &owner)
}

// ChangeRoomVisibilityOwnRoom implements RoomService.
func (port *domainRoomService) ChangeRoomVisibilityOwnRoom(
	params dto.ChangeRoomVisibilityRoomDTO,
	ownerPersonId string) (dto.RoomDTO, error) {

	roomFound, owner, err := port.findOwnedRoom(params.RoomHash, ownerPersonId)
	if err != nil {
		return dto.RoomDTO{}, err
	}

	return port.changeRoomVisibility(params, &roomFound, &owner)
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

	roomHash := port.generateRoomHash(&newRoom)
	newRoom.RoomHash = &roomHash

	persistedRoom, err := port.roomRepository.Persist(newRoom)
	if err != nil {
		return dto.RoomDTO{}, nil
	}

	return port.roomToDTO(persistedRoom, &owner), nil
}

// DeleteAllExpiredRooms implements RoomService.
func (port *domainRoomService) DeleteAllExpiredRooms() error {

	sinceLastInteraction := time.Duration.Hours(4)

	timeLimit := time.Now().Add(-time.Duration(sinceLastInteraction))

	roomsToDelete, err := port.roomRepository.FindsWithoutInteractionSince(timeLimit)
	if err != nil {
		return err
	}

	roomsIdsToDelete := sliceutils.Map(roomsToDelete, func(r domain.Room) primitive.ObjectID {
		return *r.Id
	})

	return port.roomRepository.Deletes(roomsIdsToDelete)
}

// DeleteRoom implements RoomService.
func (port *domainRoomService) DeleteRoom(roomHash string) error {

	roomFound, err := port.findRoomByHash(roomHash)
	if err != nil {
		return err
	}

	return port.roomRepository.Delete(*roomFound.Id)
}

// DeleteOwnRoom implements RoomService.
func (port *domainRoomService) DeleteOwnRoom(roomHash string, ownerPersonId string) error {

	roomFound, _, err := port.findOwnedRoom(roomHash, ownerPersonId)
	if err != nil {
		return err
	}

	return port.roomRepository.Delete(*roomFound.Id)
}

// FindAllRooms implements RoomService.
func (port *domainRoomService) FindAllRooms() ([]dto.RoomDTO, error) {
	rooms, err := port.roomRepository.FindAll()
	if err != nil {
		return nil, err
	}

	roomsDTO, err := sliceutils.MapWithError(rooms, func(r domain.Room) (dto.RoomDTO, error) {
		p, err := port.findPerson(r.OwnerPersonId.Hex())
		if err != nil {
			return dto.RoomDTO{}, err
		}
		return port.roomToDTO(&r, &p), nil
	})

	if err != nil {
		return nil, err
	}

	return roomsDTO, nil
}

// FindByOwnerPersonId implements RoomService.
func (port *domainRoomService) FindByOwnerPersonId(ownerPersonId string) ([]dto.RoomDTO, error) {

	person, err := port.findPerson(ownerPersonId)
	if err != nil {
		return nil, err
	}

	rooms, err := port.roomRepository.FindByOwnerPersonId(*person.Id)
	if err != nil {
		return nil, err
	}

	roomsDTO, err := sliceutils.MapWithError(rooms, func(r domain.Room) (dto.RoomDTO, error) {
		p, err := port.findPerson(r.OwnerPersonId.Hex())
		if err != nil {
			return dto.RoomDTO{}, err
		}
		return port.roomToDTO(&r, &p), nil
	})

	if err != nil {
		return nil, err
	}

	return roomsDTO, nil
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

func (port *domainRoomService) findRoomByHash(roomHash string) (domain.Room, error) {
	roomFound, err := port.roomRepository.FindByRoomHash(roomHash)
	if err != nil {
		return domain.Room{}, err
	}
	if roomFound == nil {
		return domain.Room{}, fmt.Errorf("room with hash: %s not found", roomHash)
	}
	return *roomFound, nil
}

func (port *domainRoomService) findOwnedRoom(roomHash string, ownerPersonId string) (domain.Room, domain.Person, error) {

	owner, err := port.findPerson(ownerPersonId)
	if err != nil {
		return domain.Room{}, domain.Person{}, err
	}

	roomFound, err := port.findRoomByHash(roomHash)
	if err != nil {
		return domain.Room{}, domain.Person{}, err
	}

	if owner.Id.Hex() != roomFound.OwnerPersonId.Hex() {
		return domain.Room{}, domain.Person{}, fmt.Errorf("the room with hash: %s is not owned by user id: %s",
			roomHash, ownerPersonId,
		)
	}

	return roomFound, owner, nil
}

func (port *domainRoomService) roomToDTO(room *domain.Room, person *domain.Person) dto.RoomDTO {

	owner := dto.RoomOwnerDTO{
		PersonId:  person.Id.Hex(),
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Email:     person.Email,
	}

	return dto.RoomDTO{
		RoomHash:            *room.RoomHash,
		Owner:               owner,
		CreationDate:        room.CreationDate,
		AnonymousAccess:     room.AnonymousAccess,
		LastInteractionDate: room.LastInteractionDate,
	}

}

func (port *domainRoomService) generateRoomHash(room *domain.Room) string {
	dateStr := room.CreationDate.ToTime().Format("2006-01-02T15:04:05.000")
	data := dateStr + room.OwnerPersonId.Hex()
	return hashutil.SHA256HexEncodingTruncated(data)
}

func (port *domainRoomService) changeRoomVisibility(
	params dto.ChangeRoomVisibilityRoomDTO,
	roomFound *domain.Room,
	owner *domain.Person) (dto.RoomDTO, error) {

	roomFound.AnonymousAccess = params.NewAnonymousAccess
	roomPersisted, err := port.roomRepository.Update(*roomFound)
	if err != nil {
		return dto.RoomDTO{}, err
	}

	return port.roomToDTO(roomPersisted, owner), nil
}
