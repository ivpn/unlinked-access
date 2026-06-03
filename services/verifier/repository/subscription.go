package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"ivpn.net/auth/services/verifier/model"
)

func (d *Database) GetSubscriptions() ([]model.Subscription, error) {
	var subs []model.Subscription
	err := d.Client.Table(d.TableName).Find(&subs).Error
	return subs, err
}

func (d *Database) UpdateSubscriptions(subs []model.Subscription) error {
	if len(subs) == 0 {
		return nil
	}

	return d.Client.Transaction(func(tx *gorm.DB) error {
		for _, sub := range subs {
			result := tx.Table(d.TableName).
				Where("id = ?", sub.ID).
				Updates(map[string]any{
					"is_active":    sub.IsActive,
					"active_until": sub.ActiveUntil,
					"tier":         sub.Tier,
				})
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func joinInt64s(ids []int64) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = fmt.Sprintf("%d", id)
	}
	return strings.Join(parts, ",")
}
