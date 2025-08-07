package repository

import (
	"ivpn.net/auth/services/verifier/model"
)

func (d *Database) GetSubscriptions() ([]model.Subscription, error) {
	var subs []model.Subscription
	err := d.Client.Find(&subs).Error
	return subs, err
}

func (d *Database) UpdateSubscription(s model.Subscription) error {
	return d.Client.Model(&s).Updates(map[string]any{
		"is_active":    s.IsActive,
		"active_until": s.ActiveUntil,
		"tier":         s.Tier,
	}).Error
}
