package config

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func ConfigAll() {

	configBase()

	configLoggers()

	configDatabases()

	configUtils()

	configRepositories()

	configServices()

	configHttp()

}

func CloseAll() {

	// wait for all goroutines to end

	globalWaitGroup.Wait()

	// close mongodb

	if err := mongoDB.Client().Disconnect(context.Background()); err != nil {
		log.Error("failed to close mongodb correctly", zap.Error(err))
	} else {
		log.Debug("the mongodb was closed correctly")
	}

	// panic recovery

	errPanic := recover()
	if errPanic != nil {
		err := errPanic.(error)
		log.Error(err.Error(), zap.Error(err))
	}
	err := log.Sync()
	if err != nil {
		fmt.Printf("log.Sync error: %v", err)
	}

}

func StartFiber() {
	fiberPort := propUtils.GetIntProp("FIBER_PORT")
	log.Info("starting! ", zap.String("app", appName), zap.String("version", version))
	err := appFiber.Listen(fmt.Sprintf("0.0.0.0:%d", fiberPort))
	if err != nil {
		panic(err)
	}
}
