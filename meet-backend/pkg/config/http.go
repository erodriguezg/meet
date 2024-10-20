package config

import (
	"github.com/erodriguezg/meet/docs"
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/addons"
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/handler"
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/wshandler"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/goccy/go-json"
	"github.com/gofiber/swagger"
)

var (
	validate    *validator.Validate
	configFiber *fiber.Config
	appFiber    *fiber.App
	wsRoot      string
)

func configHttp() {
	wsRoot = "/ws"
	configValidator()
	appFiber = configAppFiber()
	v1Router := appFiber.Group("/api/v1")
	configFiberMiddlewares()
	configFiberHandlers(&v1Router)
	wshandler.InitWebSocketsHandlers(wsRoot, appFiber, log)
	configFiberStatic()
}

func configValidator() {
	validate = validator.New()
}

func configAppFiber() *fiber.App {

	customErrorHandler := addons.NewCustomFiberErrorHandler(log)

	auxConfig := fiber.Config{
		ErrorHandler: customErrorHandler.CustomFiberErrorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	}

	configFiber = &auxConfig

	return fiber.New(*configFiber)
}

func configFiberMiddlewares() {
	// panic recover
	appFiber.Use(recover.New(
		recover.Config{
			Next:              nil,
			EnableStackTrace:  true,
			StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
		},
	))

	// cors
	appFiber.Use(
		cors.New(cors.Config{
			AllowOrigins:     propUtils.GetProp("FIBER_CORS_ORIGINS"),
			AllowCredentials: true,
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		}))

	// swagger
	configSwaggerMiddleware()

	// web sockets
	appFiber.Use(wsRoot, wshandler.NewMiddlewareFunction(httpSecurityService))
}

func configFiberHandlers(v1 *fiber.Router) {

	panicIfAnyNil(personService, httpSecurityService, profileService, modelService,
		fileService, packService, buyPackService, chiliBankService, packPaymentMethodService, log)

	v1Handlers := [...]handler.FiberHandler{
		handler.NewHealthCheckHandler(log),
		handler.NewPersonFiberHandler(personService, log),
		handler.NewSecurityHandler(httpSecurityService, log),
		handler.NewProfileFiberHandler(profileService, log),
		handler.NewModelFiberHandler(modelService, httpSecurityService, log),
		handler.NewFileFiberHandler(fileService, httpSecurityService, log),
		handler.NewPackFiberHandler(packService, httpSecurityService, validate, log),
		handler.NewBuyPackHandler(buyPackService, httpSecurityService, log),
		handler.NewChiliBankAccountFiberHandler(chiliBankService, httpSecurityService, log),
		handler.NewPackPaymentMethodFiberHandler(packPaymentMethodService, httpSecurityService, log),
	}
	for _, fHandler := range v1Handlers {
		fHandler.RegisterRoutes(v1)
	}
}

func configFiberStatic() {
	appFiber.Static("/", "./public")
	// for match anything in static
	appFiber.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("public/index.html")
	})
}

func configSwaggerMiddleware() {
	if !propUtils.GetBoolProp("SWAGGER_ENABLED") {
		return
	}
	docs.SwaggerInfo.Version = version
	appFiber.Get("/swagger/*", swagger.New(swagger.Config{}))
}
