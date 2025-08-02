package repository

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ivpn.net/auth/services/generator/config"
	"ivpn.net/auth/services/generator/model"
)

type Database struct {
	Client *gorm.DB
}

func NewDB(cfg config.Config) (*Database, error) {
	db, err := connect(cfg.DB)
	if err != nil {
		return nil, err
	}

	if cfg.Service.SampleData {
		err = migrate(db)
		if err != nil {
			return nil, err
		}
	}

	return &Database{
		Client: db,
	}, nil
}

func (d *Database) Close() error {
	db, err := d.Client.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func connect(cfg config.DBConfig) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.Name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	log.Println("DB connection OK")

	return db, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Account{},
	)
	if err != nil {
		return err
	}

	log.Println("DB migration OK")

	return nil
}
