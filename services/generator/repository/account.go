package repository

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"
	"time"

	"ivpn.net/auth/services/generator/model"
)

func (d *Database) GetAccounts() ([]*model.Account, error) {
	var accounts []*model.Account

	start := time.Now()

	err := d.Client.
		Where("is_new = ?", false).
		Where("EXISTS (SELECT 1 FROM email_service WHERE email_service.accounting_id = accounts.accounting_id)").
		Find(&accounts).Error

	elapsed := time.Since(start)
	log.Printf("GetAccounts() query completed in %s", elapsed)

	return accounts, err
}

func (d *Database) GetAccountsMock(count int) ([]*model.Account, error) {
	accounts := make([]*model.Account, count)
	for i := range count {
		accounts[i] = &model.Account{
			ID:          randomId(),
			CreatedAt:   time.Now(),
			IsActive:    true,
			ActiveUntil: time.Now().AddDate(0, 1, 0),   // Active for one month
			Product:     "Tier " + string(rune(i%3+1)), // Mocking different tiers
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
