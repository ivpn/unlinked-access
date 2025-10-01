package repository

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"ivpn.net/auth/services/generator/model"
)

func (d *Database) GetAccounts() ([]*model.Account, error) {
	var accounts []*model.Account
	var err error

	if d.Cfg.Service.Mock {
		err = d.Client.Find(&accounts).Error
	} else {
		start := time.Now()
		err = d.Client.
			Where("is_new = ?", false).
			Where("EXISTS (SELECT 1 FROM services WHERE services.accounting_id = accounts.accounting_id AND accounts.is_active = true)").
			Find(&accounts).Error

		elapsed := time.Since(start)
		log.Printf("GetAccounts() query completed in %s", elapsed)
	}

	return accounts, err
}

func (d *Database) GetAccountsMock(count int) ([]*model.Account, error) {
	accounts := make([]*model.Account, count)
	for i := range count {
		accounts[i] = &model.Account{
			ID:          randomId(),
			CreatedAt:   time.Now(),
			IsActive:    true,
			ActiveUntil: time.Now().AddDate(0, i%12+1, 0), // Active for x months
			Product:     fmt.Sprintf("Tier %d", i%3+1),    // Mocking different tiers
		}
	}

	return accounts, nil
}

func (d *Database) CreateAccountsMock(count int) error {
	accounts, err := d.GetAccountsMock(count)
	if err != nil {
		log.Printf("error generating mock accounts: %v", err)
		return err
	}

	for _, account := range accounts {
		err := d.PostAccount(account)
		if err != nil {
			log.Printf("error posting mock account %s: %v", account.ID, err)
			return err
		}
		log.Printf("mock account created: %s", account.ID)
	}

	return nil
}

func (d *Database) PostAccount(account *model.Account) error {
	return d.Client.Create(account).Error
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
