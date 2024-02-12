package db

import "time"

type VoterHistory struct {
	PollId   uint      `json:"poll_id"`
	VoteId   uint      `json:"vote_id"`
	VoteDate time.Time `json:"vote_date"`
}

type Voter struct {
	VoterId     uint           `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	VoteHistory []VoterHistory `json:"history"`
}
