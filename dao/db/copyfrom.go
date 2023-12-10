// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: copyfrom.go

package db

import (
	"context"
)

// iteratorForCreateManyCowboys implements pgx.CopyFromSource.
type iteratorForCreateManyCowboys struct {
	rows                 []CreateManyCowboysParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateManyCowboys) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateManyCowboys) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Health,
		r.rows[0].Damage,
	}, nil
}

func (r iteratorForCreateManyCowboys) Err() error {
	return nil
}

func (q *Queries) CreateManyCowboys(ctx context.Context, arg []CreateManyCowboysParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"cowboys"}, []string{"name", "health", "damage"}, &iteratorForCreateManyCowboys{rows: arg})
}
