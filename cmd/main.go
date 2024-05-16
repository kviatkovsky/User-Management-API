package main

import (
	"main/cmd/api"
	"main/configs"
	"main/database"
	"main/logger"
	"main/redis"
)

func main() {
	log := logger.InitLogger()

	server := api.NewAPIServer(configs.Envs.Port, database.ConnectDb(log), redis.ConnectRedis(log))
	if err := server.Run(); err != nil {
		log.Error("Server run error. " + err.Error())
	}
}
