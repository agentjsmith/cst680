package db

type Voter struct {
	VoterId     uint           `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	VoteHistory []VoterHistory `json:"history"`
}
