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

	var ids []string
	var isActiveCases, activeUntilCases, tierCases strings.Builder

	for _, sub := range subs {
		id := sub.ID                               // assuming this is a string (e.g., UUID)
		ids = append(ids, fmt.Sprintf("'%s'", id)) // quote the string for SQL

		isActiveCases.WriteString(fmt.Sprintf("WHEN '%s' THEN %t ", id, sub.IsActive))
		activeUntilCases.WriteString(fmt.Sprintf("WHEN '%s' THEN '%s' ", id, sub.ActiveUntil.Format("2006-01-02 15:04:05")))
		tierCases.WriteString(fmt.Sprintf("WHEN '%s' THEN '%s' ", id, sub.Tier))
	}

	sql := fmt.Sprintf(`
		UPDATE %s
		SET 
			updated_at = NOW(),
			is_active = CASE id %s END,
			active_until = CASE id %s END,
			tier = CASE id %s END
		WHERE id IN (%s);
	`, d.TableName, isActiveCases.String(), activeUntilCases.String(), tierCases.String(), strings.Join(ids, ","))

	return d.Client.Exec(sql).Error
}

func joinInt64s(ids []int64) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = fmt.Sprintf("%d", id)
	}
	return strings.Join(parts, ",")
}
