package handler

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type roomFiberHandler struct {
	roomService     service.RoomService
	securityService security.HttpSecurityService
	log             *zap.Logger
}

func NewRoomFiberHandler(
	roomService service.RoomService,
	securityService security.HttpSecurityService,
	log *zap.Logger,
) FiberHandler {
	return &roomFiberHandler{
		roomService,
		securityService,
		log,
	}
}

// RegisterRoutes implements FiberHandler.
func (port *roomFiberHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/room")
	group.Get("/:hash", port.findRoomByHash)
	group.Get("/all", port.findAllRooms)
	group.Get("/owned", port.findOwnedRooms)
	group.Post("/new", port.createRoom)
	group.Post("/visibility", port.changeVisibilityRoom)
	group.Delete("/:hash", port.deleteRoom)
	group.Delete("/expired", port.deleteExpiredRooms)
}

// ShowAccount godoc
// @Summary      Find Room By Hash
// @Description  Find one room by hash
// @Tags         Room
// @Accept       json
// @Produce      json
// @Param        hash   path     string  true  "hash of room"
// @Success      200  {object}  rest.ApiResponse[dto.RoomDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/room/{hash} [get]
func (port *roomFiberHandler) findRoomByHash(c *fiber.Ctx) error {
	roomHash := c.Params("hash")
	identity, err := port.securityService.GetIdentity(c)
	if err != nil {
		return err
	}

	var personIdRequester *string
	if identity != nil {
		personIdRequester = &identity.PersonId
	}

	port.log.Debug("-> findRoomByHash", zap.String("hash", roomHash), zap.Any("personIdRequester", personIdRequester))

	room, err := port.roomService.FindRoomByHash(roomHash, personIdRequester)
	if err != nil {
		if accessErr, ok := err.(*service.RoomAccessDeniedError); ok {
			return c.Status(fiber.StatusUnauthorized).SendString(accessErr.Error())
		} else {
			return err
		}
	}

	return c.JSON(rest.ApiOk(&room))
}

// ShowAccount godoc
// @Summary      Find All Rooms
// @Description  Find all rooms
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  rest.ApiResponse[[]dto.RoomDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/room/all [get]
func (port *roomFiberHandler) findAllRooms(c *fiber.Ctx) error {
	_, err := port.securityService.MustHavePermission(domain.PermissionCodeManageSystem, c)
	if err != nil {
		return err
	}
	port.log.Debug("-> findAllRooms")
	rooms, err := port.roomService.FindAllRooms()
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkArray(rooms))
}

// ShowAccount godoc
// @Summary      Find Owned Rooms
// @Description  Find owned rooms (by identity)
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  rest.ApiResponse[[]dto.RoomDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/room/owned [get]
func (port *roomFiberHandler) findOwnedRooms(c *fiber.Ctx) error {
	identity, err := port.securityService.MustHavePermission(domain.PermissionCodeCreateRoom, c)
	if err != nil {
		return err
	}
	port.log.Debug("-> findOwnedRooms")
	rooms, err := port.roomService.FindByOwnerPersonId(identity.PersonId)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkArray(rooms))
}

// ShowAccount godoc
// @Summary      Create Room
// @Description  Create New Room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Param        data body dto.CreateRoomDTO true "Payload Data"
// @Success      200  {object}  rest.ApiResponse[dto.RoomDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/room/new [post]
func (port *roomFiberHandler) createRoom(c *fiber.Ctx) error {
	var payload dto.CreateRoomDTO
	err := c.BodyParser((&payload))
	if err != nil {
		return err
	}
	port.log.Debug("-> createRoom", zap.Any("payload", payload))

	identity, hasPermissionSys, err := port.securityService.HasPermission(domain.PermissionCodeManageSystem, c)
	if err != nil {
		return err
	}
	if identity == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
	}

	if !hasPermissionSys {
		if !identity.HasPermission(domain.PermissionCodeCreateRoom) {
			return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
		}
		payload.OwnerPersonId = identity.PersonId
	}

	roomDTO, err := port.roomService.CreateRoom(payload)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&roomDTO))
}

// ShowAccount godoc
// @Summary      Change Visibility Room
// @Description  Change visibility room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Param        data body dto.ChangeRoomVisibilityRoomDTO true "Payload Data"
// @Success      200  {object}  rest.ApiResponse[dto.RoomDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/room/visibility [post]
func (port *roomFiberHandler) changeVisibilityRoom(c *fiber.Ctx) error {
	var payload dto.ChangeRoomVisibilityRoomDTO
	err := c.BodyParser((&payload))
	if err != nil {
		return err
	}
	port.log.Debug("-> changeVisibilityRoom", zap.Any("payload", payload))

	identity, hasPermissionSys, err := port.securityService.HasPermission(domain.PermissionCodeManageSystem, c)
	if err != nil {
		return err
	}
	if identity == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
	}

	var roomDTO dto.RoomDTO
	if !hasPermissionSys {
		if !identity.HasPermission(domain.PermissionCodeCreateRoom) {
			return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
		}

		roomDTO, err = port.roomService.ChangeRoomVisibilityOwnRoom(payload, identity.PersonId)
	} else {
		roomDTO, err = port.roomService.ChangeRoomVisibility(payload)
	}

	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&roomDTO))
}

// ShowAccount godoc
// @Summary      Delete Room
// @Description  Delete room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Param        hash   path     string  true  "hash of room"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/room/{hash} [delete]
func (port *roomFiberHandler) deleteRoom(c *fiber.Ctx) error {
	roomHash := c.Params("hash")

	port.log.Debug("-> deleteRoom", zap.Any("roomHash", roomHash))

	identity, hasPermissionSys, err := port.securityService.HasPermission(domain.PermissionCodeManageSystem, c)
	if err != nil {
		return err
	}
	if identity == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
	}

	if !hasPermissionSys {
		if !identity.HasPermission(domain.PermissionCodeCreateRoom) {
			return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
		}

		err = port.roomService.DeleteOwnRoom(roomHash, identity.PersonId)
	} else {
		err = port.roomService.DeleteRoom(roomHash)
	}

	if err != nil {
		return err
	}

	okString := "ok"
	return c.JSON(rest.ApiOk(&okString))
}

// ShowAccount godoc
// @Summary      Delete Expired Rooms
// @Description  Delete expired rooms
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/room/expired [delete]
func (port *roomFiberHandler) deleteExpiredRooms(c *fiber.Ctx) error {

	_, err := port.securityService.MustHavePermission(domain.PermissionCodeManageSystem, c)
	if err != nil {
		return err
	}

	err = port.roomService.DeleteAllExpiredRooms()
	if err != nil {
		return err
	}
	okString := "ok"
	return c.JSON(rest.ApiOk(&okString))
}
