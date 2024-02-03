package repository

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

// Отвечает за запросы к базе данных
type DBAdaptor interface {
	InsertHero(name string) (int, error)
	SelectHero(id int) (string, error)
	Close() error
}

type DB struct {
	conn *sqlx.DB
}

func NewDB(cfg Config) (*DB, error) {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		slog.Error("error during connection to db error: ", err)
		return nil, err
	}
	return &DB{conn: db}, nil
}

func (db *DB) InsertHero(name string) (int, error) {
	query := `
		INSERT INTO heros (name) values ($1)
		RETURNING id
	`
	var ids []int
	err := db.conn.Select(&ids, query, name)
	if err != nil {
		return 0, err
	}
	return ids[0], nil
}

func (db *DB) SelectHero(id int) (string, error) {
	return "", nil
}

func (db *DB) Close() error {
	err := db.conn.Close()
	if err != nil {
		slog.Error("error during closing db, error: ", err)
		return err
	}
	return nil
}
