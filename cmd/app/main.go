package main

import (
	"neploy.dev/config"
	"neploy.dev/neploy"
	"neploy.dev/pkg/store"
)

// @title Neploy API
// @version 1.0
// @description Neploy API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email 6oMwz@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	config.LoadEnv()

	db, _ := store.NewConnection(config.Env)
	npy := neploy.Neploy{
		DB:   db,
		Port: config.Env.Port,
	}

	neploy.Start(npy)
}
