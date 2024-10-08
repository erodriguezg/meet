package handler

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type fileFiberHandler struct {
	fileService     service.FileService
	securityService security.HttpSecurityService
	log             *zap.Logger
}

func NewFileFiberHandler(
	fileService service.FileService,
	securityService security.HttpSecurityService,
	log *zap.Logger,
) FiberHandler {
	return &fileFiberHandler{fileService, securityService, log}
}

func (port *fileFiberHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/file")
	group.Get("/storage-type", port.getStorageType)
	group.Get("/redirect/:hash", port.redirectDownloadUrl)
	group.Get("/get/:hash", port.getDownloadUrl)
	group.Post("/confirm/:hash", port.confirmUploaded)
}

// ShowAccount godoc
// @Summary      Get Storage Type
// @Description  Get the storage type
// @Tags         File
// @Accept       json
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/file/storage-type [get]
func (port *fileFiberHandler) getStorageType(c *fiber.Ctx) error {
	port.log.Debug("-> getStorageType")
	storageType := port.fileService.GetStorageType()
	output := map[string]string{"storageType": storageType}
	return c.JSON(rest.ApiOk(&output))
}

// ShowAccount godoc
// @Summary      Get Download Url
// @Description  Get the download url from hash
// @Tags         File
// @Accept       json
// @Param        hash  path     string  true  "unique hash for the file"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/file/get/{hash} [get]
func (port *fileFiberHandler) getDownloadUrl(c *fiber.Ctx) error {
	hashParam := c.Params("hash")
	port.log.Debug("-> getDownloadUrl", zap.String("hash", hashParam))

	downloadUrl, err := port.fileService.GetDownloadUrl(hashParam)
	if err != nil {
		return err
	}
	output := map[string]string{"url": downloadUrl}
	return c.JSON(rest.ApiOk(&output))
}

// ShowAccount godoc
// @Summary      Redirect To File
// @Description  Redirect for download the file
// @Tags         File
// @Accept       json
// @Param        hash  path     string  true  "unique hash for the file"
// @Success      301  {object}  string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/file/redirect/{hash} [get]
func (port *fileFiberHandler) redirectDownloadUrl(c *fiber.Ctx) error {
	hashParam := c.Params("hash")
	port.log.Debug("-> redirectDownloadUrl", zap.String("hash", hashParam))

	downloadUrl, err := port.fileService.GetDownloadUrl(hashParam)
	if err != nil {
		return err
	}
	return c.Redirect(downloadUrl)
}

// ShowAccount godoc
// @Summary      Confirm File Uploaded
// @Description  Confirm the file was uploaded
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        hash  path     string  true  "unique hash for the file"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/file/confirm/{hash} [post]
func (port *fileFiberHandler) confirmUploaded(c *fiber.Ctx) error {

	_, err := port.securityService.GetIdentity(c)
	if err != nil {
		return err
	}

	hashParam := c.Params("hash")

	err = port.fileService.ConfirmUploaded(hashParam)
	if err != nil {
		return err
	}

	return c.JSON(rest.ApiOkEmpty())
}
