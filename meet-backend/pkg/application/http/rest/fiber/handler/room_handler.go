package handler

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/domain"
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
	err := port.securityService.MustHavePermission(domain.PermissionCodeManageSystem, c)
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
	return nil
}

// ShowAccount godoc
// @Summary      Create Room
// @Description  Create New Room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  rest.ApiResponse[dto.RoomDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/room/new [post]
func (port *roomFiberHandler) createRoom(c *fiber.Ctx) error {
	return nil
}

// ShowAccount godoc
// @Summary      Change Visibility Room
// @Description  Change visibility room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  rest.ApiResponse[dto.RoomDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/room/visibility [post]
func (port *roomFiberHandler) changeVisibilityRoom(c *fiber.Ctx) error {
	return nil
}

func (port *roomFiberHandler) deleteRoom(c *fiber.Ctx) error {
	return nil
}

func (port *roomFiberHandler) deleteExpiredRooms(c *fiber.Ctx) error {
	return nil
}
