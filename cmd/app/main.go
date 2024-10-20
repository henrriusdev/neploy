package main

import (
	"neploy.dev/config"
	"neploy.dev/neploy"
	"neploy.dev/pkg/store"
)

func main() {
	config.LoadEnv()

	db, _ := store.NewConnection(config.Env)
	npy := neploy.Neploy{
		DB:   db,
		Port: config.Env.Port,
	}

	neploy.Start(npy)
}
