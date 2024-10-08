package config

import (
	"github.com/erodriguezg/meet/pkg/core/repository"
	"github.com/erodriguezg/meet/pkg/infrastructure/awscli"
	"github.com/erodriguezg/meet/pkg/infrastructure/dropboxcli"
	"github.com/erodriguezg/meet/pkg/infrastructure/mongodb"
	"github.com/erodriguezg/meet/pkg/infrastructure/paypalcli"
)

var (
	personRepository            repository.PersonRepository
	profileRepository           repository.ProfileRepository
	storageRepository           repository.StorageRepository
	modelRepository             repository.ModelRepository
	fileMetaDataRepository      repository.FileMetaDataRepository
	ownedResourceRepository     repository.OwnedResourceRepository
	packRepository              repository.PackRepository
	paymentClientRepository     repository.PaymentClientRepository
	paymentOrderRepository      repository.PaymentOrderRepository
	chiliBankRepository         repository.ChiliBankAccountRepository
	packPaymentMethodRepository repository.PackPaymentMethodRepository
)

func configRepositories() {
	personRepository = configPersonRepository()
	profileRepository = configProfileRepository()
	storageRepository = configStorageRepository()
	modelRepository = configModelRepository()
	fileMetaDataRepository = configFileMetaDataRepository()
	ownedResourceRepository = configOwnedResourceRepository()
	packRepository = configPackRepository()
	paymentClientRepository = configPaymentClientRepository()
	paymentOrderRepository = configPaymentOrderRepository()
	chiliBankRepository = configChiliBankRepository()
	packPaymentMethodRepository = configPackPaymentMethodRepository()
}

func configPersonRepository() repository.PersonRepository {
	return mongodb.NewPersonMongoDB(mongoDB)
}

func configProfileRepository() repository.ProfileRepository {
	return mongodb.NewProfileMongoDB(mongoDB)
}

func configModelRepository() repository.ModelRepository {
	return mongodb.NewModelMongoDB(mongoDB)
}

func configStorageRepository() repository.StorageRepository {
	storageType := propUtils.GetProp("STORAGE_TYPE")
	if storageType == "S3" {
		return configS3StorageRepository()
	} else if storageType == "DROPBOX" {
		return configDropboxStorageRepository()
	} else {
		panic("incompatible storage type: " + storageType)
	}
}

func configDropboxStorageRepository() repository.StorageRepository {
	appKey := propUtils.GetProp("DROPBOX_APP_KEY")
	appSecret := propUtils.GetProp("DROPBOX_APP_SECRET")
	refreshToken := propUtils.GetProp("DROPBOX_REFRESH_TOKEN")
	panicIfAnyNil(httpClient, log)
	return dropboxcli.NewStorageDropbox(
		appKey,
		appSecret,
		refreshToken,
		httpClient,
		log,
	)
}

func configS3StorageRepository() repository.StorageRepository {
	config := awscli.S3StorageConfig{
		RestEndPoint:           propUtils.GetProp("S3_REST_ENDPOINT"),
		AccessKey:              propUtils.GetProp("S3_ACCESS_KEY"),
		SecretAccessKey:        propUtils.GetProp("S3_SECRET_ACCESS_KEY"),
		Region:                 propUtils.GetProp("S3_REGION"),
		BucketName:             propUtils.GetProp("S3_BUCKET"),
		PathStyleAccessEnabled: propUtils.GetBoolProp("S3_PATH_STYLE_ENABLED"),
	}
	panicIfAnyNil(log)
	return awscli.NewS3StorageClient(config, log)
}

func configFileMetaDataRepository() repository.FileMetaDataRepository {
	panicIfAnyNil(mongoDB)
	return mongodb.NewFileMetaDataMongoDB(mongoDB)
}

func configOwnedResourceRepository() repository.OwnedResourceRepository {
	panicIfAnyNil(mongoDB)
	return mongodb.NewOwnedResourceMongoDB(mongoDB)
}

func configPackRepository() repository.PackRepository {
	panicIfAnyNil(mongoDB)
	return mongodb.NewPackMongoDB(mongoDB)
}

func configPaymentClientRepository() repository.PaymentClientRepository {
	panicIfAnyNil(httpClient)
	apiUrl := propUtils.GetProp("PAYPAL_API_URL")
	clientId := propUtils.GetProp("PAYPAL_CLIENT_ID")
	clientSecret := propUtils.GetProp("PAYPAL_APP_SECRET")
	return paypalcli.NewPayPalPaymentClientRepository(httpClient, apiUrl, clientId, clientSecret)
}

func configPaymentOrderRepository() repository.PaymentOrderRepository {
	panicIfAnyNil(mongoDB)
	return mongodb.NewPaymentOrderMongoDB(mongoDB)
}

func configChiliBankRepository() repository.ChiliBankAccountRepository {
	panicIfAnyNil(mongoDB)
	return mongodb.NewChiliBankAccountMongoDB(mongoDB)
}

func configPackPaymentMethodRepository() repository.PackPaymentMethodRepository {
	panicIfAnyNil(mongoDB)
	return mongodb.NewPackPaymentMethodRepository(mongoDB)
}
