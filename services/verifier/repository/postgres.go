package repository

import (
	"fmt"
	"log"
	"time"

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

// postgresManifestLog mirrors model.ManifestLog but widens Version to int64,
// since Postgres stores it as bigint while other stores keep using int.
type postgresManifestLog struct {
	ID             string    `json:"id"`
	Version        int64     `json:"version"`
	CreatedAt      time.Time `json:"created_at"`
	SignatureValid bool      `json:"signature_valid"`
	Status         string    `json:"status"`
}

func newPostgresManifestLog(log model.ManifestLog) postgresManifestLog {
	return postgresManifestLog{
		ID:             log.ID,
		Version:        int64(log.Version),
		CreatedAt:      log.CreatedAt,
		SignatureValid: log.SignatureValid,
		Status:         log.Status,
	}
}

func (l postgresManifestLog) toModel() model.ManifestLog {
	return model.ManifestLog{
		ID:             l.ID,
		Version:        int(l.Version),
		CreatedAt:      l.CreatedAt,
		SignatureValid: l.SignatureValid,
		Status:         l.Status,
	}
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
		sslMode = "require"
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

	now := time.Now()
	return d.Client.Transaction(func(tx *gorm.DB) error {
		for _, sub := range subs {
			result := tx.Table(d.TableName).
				Where("id = ?", sub.ID).
				Updates(map[string]any{
					"active_until": sub.ActiveUntil,
					"tier":         sub.Tier,
					"updated_at":   now,
				})
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func (d *PostgresDB) GetLatestManifestLog() (model.ManifestLog, error) {
	var logEntry postgresManifestLog
	err := d.Client.Table("manifest_logs").Order("version DESC").First(&logEntry).Error
	if err != nil {
		return model.ManifestLog{Version: 0}, nil
	}
	return logEntry.toModel(), nil
}

func (d *PostgresDB) AddManifestLog(log model.ManifestLog) error {
	entry := newPostgresManifestLog(log)
	err := d.Client.Table("manifest_logs").Create(&entry).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *PostgresDB) CleanupManifestLogs() error {
	expirationTime := time.Now().AddDate(0, 0, -7)
	err := d.Client.Table("manifest_logs").Where("created_at < ?", expirationTime).Delete(&postgresManifestLog{}).Error
	if err != nil {
		return err
	}
	return nil
}
