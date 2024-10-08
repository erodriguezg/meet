package configpropertyutil

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type GoEnvConfigPropertiesUtil struct {
}

func InstanceGoEnvConfigPropertiesUtil() ConfigPropertiesUtil {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found!")
	}
	return &GoEnvConfigPropertiesUtil{}
}

func (util *GoEnvConfigPropertiesUtil) GetProp(name string) string {
	res := os.Getenv(name)
	if res == "" {
		fmt.Printf("no property %s found! \n", name)
	}
	return res
}

func (util *GoEnvConfigPropertiesUtil) GetIntProp(name string) int {
	return int(util.GetInt32Prop(name))
}

func (util *GoEnvConfigPropertiesUtil) GetInt32Prop(name string) int32 {
	result, err := strconv.ParseInt(util.GetProp(name), 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(result)
}

func (util *GoEnvConfigPropertiesUtil) GetInt64Prop(name string) int64 {
	result, err := strconv.ParseInt(util.GetProp(name), 10, 64)
	if err != nil {
		panic(err)
	}
	return result
}

func (util *GoEnvConfigPropertiesUtil) GetBoolProp(name string) bool {
	result, err := strconv.ParseBool(util.GetProp(name))
	if err != nil {
		panic(err)
	}
	return result
}

func (util *GoEnvConfigPropertiesUtil) GetFloat64Prop(name string) float64 {
	result, err := strconv.ParseFloat(util.GetProp(name), 64)
	if err != nil {
		panic(err)
	}
	return result
}

func (util *GoEnvConfigPropertiesUtil) GetStringArray(name string, splitBy string) []string {
	text := util.GetProp(name)
	return strings.Split(text, splitBy)
}
