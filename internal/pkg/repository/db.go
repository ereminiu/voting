package repository

import (
	"fmt"
	"log/slog"

	"github.com/ereminiu/voting/internal/config"
	"github.com/ereminiu/voting/internal/events"
	"github.com/jmoiron/sqlx"
)

// Отвечает за запросы к базе данных
type DBAdaptor interface {
	InsertHero(name string) (int, error)
	SelectHero(id int) (string, error)
	CreatePoll(pollEvent events.PollEvent) (int, []int, error)
	Close() error
}

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		slog.Error("error during connection to db error: ", err)
		return nil, err
	}
	return &DB{conn: db}, nil
}

func (db *DB) CreatePoll(pollEvent events.PollEvent) (int, []int, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return -1, nil, err
	}

	var pollId int
	createPollQuery := fmt.Sprintf("INSERT INTO polls (closed) VALUES(DEFAULT) RETURNING poll_id")
	row := tx.QueryRow(createPollQuery)
	err = row.Scan(&pollId)
	if err != nil {
		tx.Rollback()
		return -1, nil, err
	}

	// TODO: batch insert
	choiceIds := make([]int, 0, len(pollEvent.Choices))
	for _, choice := range pollEvent.Choices {
		var id int
		addChoiceQuery := fmt.Sprintf(`
			INSERT INTO choices (name) VALUES ($1)
			ON CONFLICT DO NOTHING
			RETURNING choice_id
		`)
		row := tx.QueryRow(addChoiceQuery, choice.Name)
		err = row.Scan(&id)
		if err != nil {
			selectChoiceId := fmt.Sprintf(`
				SELECT choice_id FROM choices
				WHERE name = $1
			`)
			row := tx.QueryRow(selectChoiceId, choice.Name)
			err = row.Scan(&id)
			if err != nil {
				tx.Rollback()
				return -1, nil, err
			}
			// tx.Rollback()
			// return -1, nil, err
		}

		choiceIds = append(choiceIds, id)
	}

	return pollId, choiceIds, tx.Commit()
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
