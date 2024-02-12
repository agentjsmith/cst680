package db

import "errors"

type Voter struct {
	VoterId     uint           `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	VoteHistory []VoterHistory `json:"history"`
}

func (v *Voter) GetHistoryByPollId(pollId uint) (VoterHistory, error) {
	for i := range v.VoteHistory {
		if v.VoteHistory[i].PollId == pollId {
			return v.VoteHistory[i], nil
		}
	}
	return VoterHistory{}, errors.New("poll not found in voter history")
}
