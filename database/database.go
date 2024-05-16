package database

import (
	"fmt"
	"log/slog"
	"main/configs"
	"main/types"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDb(sLog *slog.Logger) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=%s sslmode=disable timezone=UTC",
		configs.Envs.DBUser,
		configs.Envs.DBPassword,
		configs.Envs.DBName,
		configs.Envs.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		sLog.Error(err.Error())
	}

	sLog.Info("Connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)

	sLog.Info("Running migrations")
	err = db.AutoMigrate(&types.User{})
	if err != nil {
		sLog.Error(err.Error())
	}

	err = db.AutoMigrate(&types.Votes{})
	if err != nil {
		sLog.Error(err.Error())
	}

	err = db.AutoMigrate(&types.ACL{})
	if err != nil {
		sLog.Error(err.Error())
	}

	err = db.AutoMigrate(&types.UserRole{})
	if err != nil {
		sLog.Error(err.Error())
	}

	userRoles := types.GetUserRoleData()
	for _, v := range userRoles {
		err = db.Create(&v).Error
		if err != nil {
			sLog.Error(err.Error())
		}
	}

	return db
}
