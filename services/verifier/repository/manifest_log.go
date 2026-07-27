package repository

import (
	"time"

	"ivpn.net/auth/services/verifier/model"
)

func (d *Database) GetLatestManifestLog() (model.ManifestLog, error) {
	var log model.ManifestLog
	err := d.Client.Table(d.ManifestLogTableName).Order("created_at DESC").First(&log).Error
	if err != nil {
		return model.ManifestLog{Version: 0}, nil
	}
	return log, nil
}

func (d *Database) AddManifestLog(log model.ManifestLog) error {
	err := d.Client.Table(d.ManifestLogTableName).Create(&log).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) CleanupManifestLogs() error {
	expirationTime := time.Now().AddDate(0, 0, -7)
	err := d.Client.Table(d.ManifestLogTableName).Where("created_at < ?", expirationTime).Delete(&model.ManifestLog{}).Error
	if err != nil {
		return err
	}
	return nil
}
