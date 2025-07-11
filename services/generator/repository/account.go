package repository

import (
	"math/rand"
	"strings"
	"time"

	"ivpn.net/auth/services/generator/model"
)

func (d *Database) GetAccounts() ([]*model.Account, error) {
	var accounts []*model.Account
	err := d.Client.Where("is_new = ?", false).Find(&accounts).Error

	return accounts, err
}

func (d *Database) GetAccountsMock(count int) ([]*model.Account, error) {
	accounts := make([]*model.Account, count)
	for i := range count {
		accounts[i] = &model.Account{
			ID:          randomId(),
			CreatedAt:   time.Now(),
			IsActive:    true,
			ActiveUntil: time.Now().AddDate(0, 1, 0), // Active for one month
		}
	}

	return accounts, nil
}

func randomId() string {
	// Generate a random ID, e.g., i-1234-ABCD-XYQZ
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var id strings.Builder
	id.WriteString("i-")

	for i := 0; i < 4; i++ {
		id.WriteByte(charset[rand.Intn(len(charset))])
		if i == 3 {
			id.WriteByte('-')
		}
	}
	for i := 0; i < 4; i++ {
		id.WriteByte(charset[rand.Intn(len(charset))])
		if i == 3 {
			id.WriteByte('-')
		}
	}
	for i := 0; i < 4; i++ {
		id.WriteByte(charset[rand.Intn(len(charset))])
	}

	return id.String()
}
