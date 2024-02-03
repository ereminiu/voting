package service

import (
	"github.com/ereminiu/voting/internal/pkg/repository"
)

type ServiceAdaptor interface {
	CreateHero(name string) (int, error)
	GetHeroByID(id int) (string, error)
}

type PollService struct {
	db *repository.DB
}

func NewPollService(db *repository.DB) (*PollService, error) {
	return &PollService{db: db}, nil
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
