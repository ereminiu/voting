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
	Poll(poll events.VoteEvent) error
	GetResults(finishEvent events.PollFinishingEvent) (int, string, error)
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
		addChoiceQuery := fmt.Sprintf(`
			INSERT INTO choices (name) VALUES ($1)
			ON CONFLICT DO NOTHING
			RETURNING choice_id
		`)
		var choiceId int
		row := tx.QueryRow(addChoiceQuery, choice.Name)
		err = row.Scan(&choiceId)
		if err != nil {
			selectChoiceId := fmt.Sprintf(`
				SELECT choice_id FROM choices
				WHERE name = $1
			`)
			row := tx.QueryRow(selectChoiceId, choice.Name)
			err = row.Scan(&choiceId)
			if err != nil {
				tx.Rollback()
				return -1, nil, err
			}
		}

		choiceIds = append(choiceIds, choiceId)
	}

	for _, choiceId := range choiceIds {
		addChoiceToPollRelation := fmt.Sprintf(`
			INSERT INTO choices_to_polls (hero_id, poll_id)
			VALUES ($1, $2)
		`)
		_, err = tx.Exec(addChoiceToPollRelation, choiceId, pollId)
		if err != nil {
			tx.Rollback()
			return -1, nil, err
		}
	}

	return pollId, choiceIds, tx.Commit()
}

func (db *DB) Poll(poll events.VoteEvent) error {
	// TODO: провять что голосование открыто
	updAmount := `
		UPDATE choices_to_polls 
		SET amount = amount + 1
		WHERE hero_id = $1 AND poll_id = $2
	`
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(updAmount, poll.ChoiceId, poll.PollId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) GetResults(finishEvent events.PollFinishingEvent) (int, string, error) {
	closePoll := `
		UPDATE polls 
		SET closed = true
		WHERE poll_id = $1
	`

	tx, err := db.conn.Begin()
	if err != nil {
		return -1, "", err
	}

	_, err = tx.Exec(closePoll, finishEvent.PollId)
	if err != nil {
		tx.Rollback()
		return -1, "", err
	}

	selectWinner := `
		SELECT ch.choice_id, ch.name 
		FROM choices ch 
		JOIN choices_to_polls ctp 
		ON ctp.hero_id = ch.choice_id
		JOIN polls p 
		ON p.poll_id = ctp.poll_id
		WHERE p.poll_id = $1
		ORDER BY ctp.amount DESC
		LIMIT 1
	`
	row := tx.QueryRow(selectWinner, finishEvent.PollId)

	var id int
	var name string
	err = row.Scan(&id, &name)
	if err != nil {
		tx.Rollback()
		return -1, "", err
	}

	return id, name, tx.Commit()
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
