package events

type Choice struct {
	Name     string `json:"name,omitempty" db:"name"`
	ChoiceId int    `json:"choice_id,omitempty" db:"choice_id"`
}

type PollEvent struct {
	Choices []Choice `json:"choices"`
}

type PollFinishingEvent struct {
	PollId int `json:"poll_id"`
}

type VoteEvent struct {
	PollId   int `json:"poll_id"`
	ChoiceId int `json:"choice_id"`
}
