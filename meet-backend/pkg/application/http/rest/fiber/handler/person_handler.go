package handler

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type personFiberHandler struct {
	personService service.PersonService
	log           *zap.Logger
}

func NewPersonFiberHandler(personService service.PersonService, log *zap.Logger) FiberHandler {
	return &personFiberHandler{personService, log}
}

func (port *personFiberHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/person")
	group.Get("/all", port.findAllPersons)
	group.Get("/filter", port.filterPersons)
	group.Get("/:uuid", port.findById)
	group.Post("/save", port.savePerson)
	group.Delete("/:uuid", port.deletePerson)
}

// ShowAccount godoc
// @Summary      Filter Persons
// @Description  Search persons filtered
// @Tags         Person
// @Accept       json
// @Produce      json
// @Param        id   query     string  false  "id of person"
// @Param        rut  query     string  false  "rut of person"
// @Param        nameLike   query     string  false  "text search on names"
// @Param        birthdayLower   query     time.Time  false "birthdayLower"
// @Param        birthdayUpper   query     time.Time  false "birthdayUpper"
// @Success      200  {object}  rest.ApiResponse[[]domain.Person]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/person/filter [get]
func (port *personFiberHandler) filterPersons(c *fiber.Ctx) error {
	var filters domain.PersonFilter
	err := c.QueryParser(&filters)
	if err != nil {
		return err
	}
	port.log.Debug("-> filterPersons", zap.Any("filters", filters))
	persons, err := port.personService.FilterPaginated(filters)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&persons))
}

// ShowAccount godoc
// @Summary      Find All Persons
// @Description  Get all persons
// @Tags         Person
// @Accept       json
// @Produce      json
// @Success      200  {object}  rest.ApiResponse[[]domain.Person]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/person/all [get]
func (port *personFiberHandler) findAllPersons(c *fiber.Ctx) error {
	port.log.Debug("-> findAllPersons")
	persons, err := port.personService.FindAll()
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&persons))
}

// ShowAccount godoc
// @Summary      Find By Id
// @Description  Find person by uuid id
// @Tags         Person
// @Accept       json
// @Produce      json
// @Param        uuid   path     string  false  "id of person"
// @Success      200  {object}  rest.ApiResponse[domain.Person]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/person/{uuid} [get]
func (port *personFiberHandler) findById(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	port.log.Debug("-> findById", zap.String("uuid", uuidParam))
	person, err := port.personService.FindById(uuidParam)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&person))
}

// ShowAccount godoc
// @Summary      Save Person
// @Description  Insert or update a person
// @Tags         Person
// @Accept       json
// @Produce      json
// @Param        data body domain.Person true "The input person"
// @Success      200  {object}  rest.ApiResponse[domain.Person]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/person/save [post]
func (port *personFiberHandler) savePerson(c *fiber.Ctx) error {
	var person domain.Person
	err := c.BodyParser(&person)
	if err != nil {
		return err
	}
	port.log.Debug("-> savePerson", zap.Any("person", &person))
	updatedPerson, err := port.personService.Save(person)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&updatedPerson))
}

// ShowAccount godoc
// @Summary      Delete Person By Id
// @Description  Delete one person by his uuid
// @Tags         Person
// @Accept       json
// @Produce      json
// @Param        uuid   path     string  false  "id of person"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/person/{uuid} [delete]
func (port *personFiberHandler) deletePerson(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	port.log.Debug("-> deletePerson", zap.String("uuid", uuidParam))
	err := port.personService.Delete(uuidParam)
	if err != nil {
		return err
	}
	okString := "ok"
	return c.JSON(rest.ApiOk(&okString))
}
