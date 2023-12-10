package dao

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ShootAtCowboyTx performs a cockroachdb transaction
func ShootAtCowboyTx(ctx context.Context, tx pgx.Tx, shooterId uuid.UUID, shooterName string, shooterDmg int32, st map[string]interface{}) error {

	var receiverId string
	var receiverHealth int32
	var receiverName string
	var shooterHealth int32
	var healthAfterDmg int32

	if err := tx.QueryRow(ctx,
		"SELECT health from cowboys where id = $1 AND health > 0", shooterId).Scan(&shooterHealth); err != nil {
		return err
	}

	if err := tx.QueryRow(ctx,
		"SELECT id, health, name from cowboys Where health > 0 AND id != $1 order by random() LIMIT 1;", shooterId.String()).Scan(&receiverId, &receiverHealth, &receiverName); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		"UPDATE cowboys SET health = (health - $1), updated_at = NOW() WHERE id = $2;", shooterDmg, receiverId); err != nil {
		return err
	}

	if err := tx.QueryRow(ctx,
		"SELECT health from cowboys where id = $1;", receiverId).Scan(&healthAfterDmg); err != nil {
		return err
	}

	st["shooter_name"] = shooterName
	st["shooter_health"] = shooterHealth
	st["receiver_name"] = receiverName
	st["receiver_id"] = receiverId
	st["shooter_damage"] = shooterDmg
	st["receiver_health"] = receiverHealth
	st["receiver_health_after"] = healthAfterDmg

	return nil
}
