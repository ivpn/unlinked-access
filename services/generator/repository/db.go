package repository

import (
	"log"

	mysqldrv "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ivpn.net/auth/services/generator/config"
	"ivpn.net/auth/services/generator/model"
)

type Database struct {
	Client *gorm.DB
	Cfg    config.Config
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
		Cfg:    cfg,
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
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	dsnCfg := mysqldrv.Config{
		User:                 cfg.User,
		Passwd:               cfg.Password,
		Net:                  "tcp",
		Addr:                 cfg.Host + ":" + cfg.Port,
		DBName:               cfg.Name,
		Params:               map[string]string{"charset": "utf8mb4"},
		ParseTime:            true,
		Loc:                  nil, // use UTC
		AllowNativePasswords: true,
	}
	dsn := dsnCfg.FormatDSN()

	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
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
