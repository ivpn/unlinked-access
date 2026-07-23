package repository

import (
	"fmt"
	"strings"

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

	n := len(subs)
	var (
		activeUntilCase strings.Builder
		tierCase        strings.Builder
		idPlaceholders  strings.Builder
	)
	// 2 args per sub for active_until CASE, 2 for tier CASE, 1 for WHERE IN
	args := make([]any, 0, n*5)

	for _, sub := range subs {
		activeUntilCase.WriteString("WHEN ? THEN ? ")
		args = append(args, sub.ID, sub.ActiveUntil)
	}

	for _, sub := range subs {
		tierCase.WriteString("WHEN ? THEN ? ")
		args = append(args, sub.ID, sub.Tier)
	}

	for i, sub := range subs {
		if i > 0 {
			idPlaceholders.WriteString(",")
		}
		idPlaceholders.WriteString("?")
		args = append(args, sub.ID)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET active_until = CASE id %s END, tier = CASE id %s END WHERE id IN (%s)",
		d.TableName,
		activeUntilCase.String(),
		tierCase.String(),
		idPlaceholders.String(),
	)

	return d.Client.Exec(query, args...).Error
}

func joinInt64s(ids []int64) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = fmt.Sprintf("%d", id)
	}
	return strings.Join(parts, ",")
}
