package repository

import (
	"log"

	mysqldrv "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

type Database struct {
	Client                *gorm.DB
	SubscriptionTableName string
	ManifestLogTableName  string
}

func NewDB(cfg config.Config) (*Database, error) {
	db, err := connect(cfg.DB)
	if err != nil {
		return nil, err
	}

	if cfg.Service.SampleData {
		err = migrateSampleData(db, cfg.DB.Table)
		if err != nil {
			return nil, err
		}
	}

	err = migrateManifestLog(db, "manifest_logs")
	if err != nil {
		return nil, err
	}

	return &Database{
		Client:                db,
		SubscriptionTableName: cfg.DB.Table,
		ManifestLogTableName:  "manifest_logs",
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
		Loc:                  nil,
		AllowNativePasswords: true,
	}
	dsn := dsnCfg.FormatDSN()

	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
	if err != nil {
		log.Println("DB connection ERROR:", err)
		return nil, err
	}

	log.Println("DB connection OK")

	return db, nil
}

func migrateManifestLog(db *gorm.DB, tableName string) error {
	err := db.Table(tableName).AutoMigrate(
		&model.ManifestLog{},
	)
	if err != nil {
		return err
	}

	log.Println("DB migration OK")

	return nil
}

func migrateSampleData(db *gorm.DB, tableName string) error {
	err := db.Table(tableName).AutoMigrate(
		&model.Subscription{},
	)
	if err != nil {
		return err
	}

	log.Println("DB migration OK")

	return nil
}
