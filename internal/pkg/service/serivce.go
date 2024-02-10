package service

import (
	"log/slog"

	"github.com/ereminiu/voting/internal/events"
	"github.com/ereminiu/voting/internal/pkg/repository"
)

type ServiceAdaptor interface {
	CreateHero(name string) (int, error)
	GetHeroByID(id int) (string, error)
	CreatePoll(pollEvent events.PollEvent) (int, []int, error)
	Poll(poll events.VoteEvent) error
	GetResults(finishEvent events.PollFinishingEvent) (int, string, error)
}

type PollService struct {
	log *slog.Logger
	db  *repository.DB
}

func NewPollService(db *repository.DB) (*PollService, error) {
	return &PollService{db: db}, nil
}

func (pc *PollService) GetResults(finishEvent events.PollFinishingEvent) (int, string, error) {
	id, name, err := pc.db.GetResults(finishEvent)
	if err != nil {
		return -1, "", err
	}
	return id, name, nil
}

func (pc *PollService) Poll(poll events.VoteEvent) error {
	err := pc.db.Poll(poll)
	if err != nil {
		return err
	}
	return nil
}

func (pc *PollService) CreatePoll(pollEvent events.PollEvent) (int, []int, error) {
	pollId, choiceIds, err := pc.db.CreatePoll(pollEvent)
	if err != nil {
		return -1, nil, err
	}
	return pollId, choiceIds, err
}

func (pc *PollService) CreateHero(name string) (int, error) {
	id, err := pc.db.InsertHero(name)
	if err != nil {
		return -1, err
	}
	return id, err
}

func (pc *PollService) GetHeroByID(id int) (string, error) {
	name, err := pc.db.SelectHero(id)
	if err != nil {
		return "", err
	}
	return name, nil
}
