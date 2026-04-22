package repository

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

type PostgresDB struct {
	Client    *gorm.DB
	TableName string
}

func NewPostgresDB(cfg config.Config) (*PostgresDB, error) {
	db, err := connectPostgres(cfg.PGDB)
	if err != nil {
		return nil, err
	}

	if cfg.Service.SampleData {
		err = migratePostgres(db, cfg.PGDB.Table)
		if err != nil {
			return nil, err
		}
	}

	return &PostgresDB{
		Client:    db,
		TableName: cfg.PGDB.Table,
	}, nil
}

func (d *PostgresDB) Close() error {
	db, err := d.Client.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func connectPostgres(cfg config.PGDBConfig) (*gorm.DB, error) {
	sslMode := cfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	log.Println("PostgresDB connection OK")
	return db, nil
}

func migratePostgres(db *gorm.DB, tableName string) error {
	err := db.Table(tableName).AutoMigrate(&model.Subscription{})
	if err != nil {
		return err
	}
	log.Println("PostgresDB migration OK")
	return nil
}

func (d *PostgresDB) GetSubscriptions() ([]model.Subscription, error) {
	var subs []model.Subscription
	err := d.Client.Table(d.TableName).Find(&subs).Error
	return subs, err
}

func (d *PostgresDB) UpdateSubscriptions(subs []model.Subscription) error {
	if len(subs) == 0 {
		return nil
	}

	var ids []string
	var isActiveCases, activeUntilCases, tierCases strings.Builder

	for _, sub := range subs {
		id := sub.ID
		ids = append(ids, fmt.Sprintf("'%s'", id))

		isActiveCases.WriteString(fmt.Sprintf("WHEN '%s' THEN %t ", id, sub.IsActive))
		activeUntilCases.WriteString(fmt.Sprintf("WHEN '%s' THEN '%s' ", id, sub.ActiveUntil.Format("2006-01-02 15:04:05")))
		tierCases.WriteString(fmt.Sprintf("WHEN '%s' THEN '%s' ", id, sub.Tier))
	}

	sql := fmt.Sprintf(`
		UPDATE %s
		SET
			updated_at = NOW(),
			is_active = CASE id %s END,
			active_until = CASE id %s END,
			tier = CASE id %s END
		WHERE id IN (%s);
	`, d.TableName, isActiveCases.String(), activeUntilCases.String(), tierCases.String(), strings.Join(ids, ","))

	return d.Client.Exec(sql).Error
}
