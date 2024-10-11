package main

import (
	"github.com/erodriguezg/meet/pkg/config"
)

// @title meet API
// @version 1.0
// @description Api for meet application.
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api
func main() {
	defer config.CloseAll()

	config.ConfigAll()

	config.StartFiber()

}
