package repository

import (
	"crypto/rand"
	"math/big"
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
	for i := 0; i < count; i++ {
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

	max := big.NewInt(int64(len(charset)))

	for i := range 4 {
		n, _ := rand.Int(rand.Reader, max)
		id.WriteByte(charset[n.Int64()])
		if i == 3 {
			id.WriteByte('-')
		}
	}
	for i := range 4 {
		n, _ := rand.Int(rand.Reader, max)
		id.WriteByte(charset[n.Int64()])
		if i == 3 {
			id.WriteByte('-')
		}
	}
	for range 4 {
		n, _ := rand.Int(rand.Reader, max)
		id.WriteByte(charset[n.Int64()])
	}

	return id.String()
}
