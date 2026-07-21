package repository

import "ivpn.net/auth/services/verifier/model"

func (d *Database) GetLatestManifestLog() (model.ManifestLog, error) {
	var log model.ManifestLog
	err := d.Client.Table(d.TableName).Order("created_at DESC").First(&log).Error
	if err != nil {
		return model.ManifestLog{}, err
	}
	return log, nil
}
