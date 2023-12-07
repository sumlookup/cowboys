package intermediate

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// todo: fix import cycle
func shootAtCowboyTx(ctx context.Context, tx pgx.Tx, shooterId uuid.UUID, shooterName string, shooterDmg int32, st map[string]interface{}) error {

	var victimId string
	var victimHealth int
	var victimName string
	var shooterHealth string
	var healthAfterDmg int32

	if err := tx.QueryRow(ctx,
		"SELECT health from cowboys where id = $1 AND health > 0", shooterId).Scan(&shooterHealth); err != nil {
		return err
	}

	if err := tx.QueryRow(ctx,
		"SELECT id, health, name from cowboys Where health > 0 AND id != $1 order by random() LIMIT 1;", shooterId.String()).Scan(&victimId, &victimHealth, &victimName); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		"UPDATE cowboys SET health = (health - $1), updated_at = NOW() WHERE id = $2;", shooterDmg, victimId); err != nil {
		return err
	}

	if err := tx.QueryRow(ctx,
		"SELECT health from cowboys where id = $1;", victimId).Scan(&healthAfterDmg); err != nil {
		return err
	}

	st["shooter"] = shooterName
	st["shooter_health"] = shooterHealth
	st["victim_name"] = victimName
	st["shooter_damage"] = shooterDmg
	st["victim_health"] = victimHealth
	st["victim_health_after"] = healthAfterDmg

	return nil
}
