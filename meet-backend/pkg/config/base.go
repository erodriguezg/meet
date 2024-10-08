package config

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/erodriguezg/meet/pkg/util/base64util"
	"github.com/erodriguezg/meet/pkg/util/configpropertyutil"
)

var (
	version            string
	globalWaitGroup    *sync.WaitGroup
	propUtils          configpropertyutil.ConfigPropertiesUtil
	env                string
	appName            string
	rsaPrivateKeyBytes []byte
	rsaPublicKeyBytes  []byte
)

func configBase() {
	version = "0.2.0"
	appName = "go-env"
	globalWaitGroup = configGlobalWaitGroup()
	propUtils = configPropUtils()
	env = configEnviroment()
	rsaPrivateKeyBytes = configRsaPrivateKeyBytes()
	rsaPublicKeyBytes = configRsaPublicKeyBytes()
	configDefaultLocation()
}

func configGlobalWaitGroup() *sync.WaitGroup {
	var waitGroup sync.WaitGroup
	return &waitGroup
}

func configPropUtils() configpropertyutil.ConfigPropertiesUtil {
	return configpropertyutil.InstanceGoEnvConfigPropertiesUtil()
}

func configEnviroment() string {
	return propUtils.GetProp("ENV")
}

func configDefaultLocation() {
	if tz := propUtils.GetProp("DEFAULT_TIMEZONE"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			fmt.Printf("error loading location '%s': %v\n", tz, err)
		} else {
			fmt.Printf("default timezone %v \n", time.Local)
		}
	}
}

func configRsaPrivateKeyBytes() []byte {
	privateKeyB64 := propUtils.GetProp("SECURE_PRIVATE_KEY_B64")
	keyBytes, err := base64util.Decode(&privateKeyB64)
	if err != nil {
		panic(fmt.Errorf("error getting rsa private key: %s", err.Error()))
	}
	return keyBytes
}

func configRsaPublicKeyBytes() []byte {
	publicKeyB64 := propUtils.GetProp("SECURE_PUBLIC_KEY_B64")
	keyBytes, err := base64util.Decode(&publicKeyB64)
	if err != nil {
		panic(fmt.Errorf("error getting rsa public key: %s", err.Error()))
	}
	return keyBytes
}

func panicIfAnyNil(nilables ...interface{}) {
	if len(nilables) == 0 {
		return
	}
	for i, nilable := range nilables {
		if nilable == nil {
			var method string
			pc, _, _, ok := runtime.Caller(1)
			if ok {
				method = runtime.FuncForPC(pc).Name()
			} else {
				method = "unknown"
			}
			panic(fmt.Errorf("an config element is nil in the method '%s', argument number %d! ", method, i+1))
		}
	}
}
