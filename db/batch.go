// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: batch.go

package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBatchAlreadyClosed = errors.New("batch already closed")
)

const batchCreate = `-- name: BatchCreate :batchone
INSERT INTO tournaments (
    name,
    entry_fee,
    start_time
) VALUES (
    $1, $2, $3
) RETURNING id, name, entry_fee, start_time
`

type BatchCreateBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type BatchCreateParams struct {
	Name      string
	EntryFee  int64
	StartTime pgtype.Timestamp
}

func (q *Queries) BatchCreate(ctx context.Context, arg []BatchCreateParams) *BatchCreateBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.Name,
			a.EntryFee,
			a.StartTime,
		}
		batch.Queue(batchCreate, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &BatchCreateBatchResults{br, len(arg), false}
}

func (b *BatchCreateBatchResults) QueryRow(f func(int, Tournament, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		var i Tournament
		if b.closed {
			if f != nil {
				f(t, i, ErrBatchAlreadyClosed)
			}
			continue
		}
		row := b.br.QueryRow()
		err := row.Scan(
			&i.ID,
			&i.Name,
			&i.EntryFee,
			&i.StartTime,
		)
		if f != nil {
			f(t, i, err)
		}
	}
}

func (b *BatchCreateBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}