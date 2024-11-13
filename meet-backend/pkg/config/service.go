package config

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/fiberidentity"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/erodriguezg/meet/pkg/util/openid"
)

var (
	openIdService            openid.OpenIdService
	personService            service.PersonService
	profileService           service.ProfileService
	modelService             service.ModelService
	httpSecurityService      security.HttpSecurityService
	storageService           service.StorageService
	fileService              service.FileService
	ownedResourceService     service.OwnedResourceService
	packService              service.PackService
	buyPackService           service.BuyPackService
	chiliBankService         service.ChiliBankAccountService
	packPaymentMethodService service.PackPaymentMethodService
	roomService              service.RoomService
)

func configServices() {
	storageService = configStorageService()
	fileService = configFileService()
	personService = configPersonService()
	modelService = configModelService()
	openIdService = configOpenIdService()
	profileService = configProfileService()
	httpSecurityService = configHttpSecurityService()
	ownedResourceService = configOwnedResourceService()
	packService = configPackService()
	buyPackService = configBuyPackService()
	chiliBankService = configChileBankService()
	packPaymentMethodService = configPackPaymentMethodService()
	roomService = configRoomService()
}

func configOpenIdService() openid.OpenIdService {
	config := openid.OpenIdConfig{
		ResponseType: "code",
		Scope:        "email openid profile",
		ClientId:     propUtils.GetProp("GOOGLE_OPENID_CLIENT_ID"),
		ClientSecret: propUtils.GetProp("GOOGLE_OPENID_CLIENT_SECRET"),
		RedirectUri:  propUtils.GetProp("GOOGLE_OPENID_REDIRECT_URL"),
	}
	panicIfAnyNil(config, googleOAuth2Api, jwtUtil, log)
	return openid.NewGoogleOpenIdService(config, googleOAuth2Api, jwtUtil, log)
}

func configPersonService() service.PersonService {
	panicIfAnyNil(personRepository)
	return service.NewDomainPersonService(personRepository)
}

func configProfileService() service.ProfileService {
	panicIfAnyNil(profileRepository)
	return service.NewDomainProfileService(profileRepository)
}

func configHttpSecurityService() security.HttpSecurityService {
	panicIfAnyNil(openIdService, personService, profileService, rsaPrivateKeyBytes, rsaPublicKeyBytes)

	openIdPassPhrase := propUtils.GetProp("SECURE_PASSPHRASE_OPENID")

	identityUtil := fiberidentity.NewFiberIdentityUtil(personService, profileService, modelService, rsaPublicKeyBytes)
	return security.NewDefaultHttpSecurityService(openIdService, personService, identityUtil, rsaPrivateKeyBytes, openIdPassPhrase)
}

func configStorageService() service.StorageService {
	panicIfAnyNil(storageRepository)
	return service.NewDomainStorageService(storageRepository)
}

func configFileService() service.FileService {
	panicIfAnyNil(storageService, fileMetaDataRepository)
	return service.NewFileService(storageService, fileMetaDataRepository)
}

func configModelService() service.ModelService {
	panicIfAnyNil(personService, fileService, modelRepository)
	return service.NewDomainModelService(personService, fileService, modelRepository)
}

func configOwnedResourceService() service.OwnedResourceService {
	panicIfAnyNil(ownedResourceRepository)
	return service.NewDomainOwnedResourceService(ownedResourceRepository)
}

func configPackService() service.PackService {
	panicIfAnyNil(personService, profileService, modelService, ownedResourceService,
		fileService, packRepository)
	return service.NewDomainPackService(personService, profileService, modelService, ownedResourceService,
		fileService, packRepository)
}

func configBuyPackService() service.BuyPackService {
	panicIfAnyNil(personService, modelService, packService, ownedResourceService,
		paymentClientRepository, paymentOrderRepository)
	return service.NewDomainBuyPackService(personService, modelService, packService, ownedResourceService,
		paymentClientRepository, paymentOrderRepository)
}

func configChileBankService() service.ChiliBankAccountService {
	panicIfAnyNil(modelService, chiliBankRepository)
	return service.NewDomainChiliBankAccountService(
		modelService, chiliBankRepository,
	)
}

func configPackPaymentMethodService() service.PackPaymentMethodService {
	panicIfAnyNil(packService, packPaymentMethodRepository)
	return service.NewPackPaymentMethodService(packService, packPaymentMethodRepository)
}

func configRoomService() service.RoomService {
	panicIfAnyNil(roomRepository, personRepository)
	return service.NewDomainRoomService(roomRepository, personService)
}
