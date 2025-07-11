package repository

import (
	"ivpn.net/auth/services/generator/model"
)

func (d *Database) GetAccounts() ([]*model.Account, error) {
	var accounts []*model.Account
	err := d.Client.Where("is_new = ?", false).Find(&accounts).Error
	return accounts, err
}
