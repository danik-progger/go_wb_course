package db

import (
	"api/db/pg"
	"context"
	_ "embed"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	ctx context.Context
	q   *pg.Queries
}

func InitDB() (*DB, error) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	queries := pg.New(conn)

	return &DB{
		ctx: ctx,
		q:   queries,
	}, nil
}

func (d *DB) GetMsgById(id int64) (pg.Message, error) {
	return d.q.GetMsgById(d.ctx, id)
}

func (d *DB) GetAllMsgs() ([]pg.Message, error) {
	return d.q.GetAllMsgs(d.ctx)
}

func (d *DB) DeleteMsgById(id int64) error {
	return d.q.DeleteMsg(d.ctx, id)
}

func (d *DB) AddMsg(body string, send_at time.Time) (pg.Message, error) {
	return d.q.AddMsg(d.ctx, pg.AddMsgParams{
		Body:   body,
		SendAt: send_at.Format(time.RFC3339),
	})
}

func (d *DB) GetUnproccessed() ([]pg.Message, error) {
	return d.q.GetUnprocessed(d.ctx)
}

func (d *DB) Update(id int64, new_status string, send_at string) error {
	m, err := d.GetMsgById(id)
	if err != nil {
		return err
	}

	return d.q.Update(d.ctx, pg.UpdateParams{
		ID:        id,
		Status:    new_status,
		Tries:     m.Tries + 1,
		WasSentAt: send_at,
	})
}
