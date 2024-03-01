package db_test

import (
	"drexel.edu/voter-api/db"
	"testing"
)

func TestWithEmptyDb(t *testing.T) {
	emptyDb, err := db.New()
	if err != nil {
		t.Fatalf("Creating new DB: %v", err)
	}

	t.Run("GetAllVoters", func(t *testing.T) {
		voters, err := emptyDb.GetAllVoters()

		if err != nil {
			t.Errorf("Errored: %v", err)
		}

		if len(voters) > 0 {
			t.Errorf("Wanted [] got %v", voters)
		}
	})

	t.Run("GetVoter", func(t *testing.T) {
		_, err := emptyDb.GetVoter(1)

		if err == nil {
			t.Error("Succeeded but shouldn't have")
		}

	})

	t.Run("DeleteVoter", func(t *testing.T) {
		err := emptyDb.DeleteVoter(1)

		if err == nil {
			t.Error("Succeeded but shouldn't have")
		}

	})

	t.Run("UpdateVoter", func(t *testing.T) {
		err := emptyDb.UpdateVoter(db.Voter{
			VoterId:     1,
			Name:        "Ohno Wontwork",
			VoteHistory: make([]db.VoterHistory, 0),
		})

		if err == nil {
			t.Error("Succeeded but shouldn't have")
		}

	})

	// This has to run last in this function because after this the DB will no longer be empty
	t.Run("AddVoter", func(t *testing.T) {
		theCount := db.Voter{
			VoterId:     1,
			Name:        "Count Chocula",
			VoteHistory: make([]db.VoterHistory, 0),
		}

		err := emptyDb.AddVoter(theCount)

		if err != nil {
			t.Errorf("Failed: %v", err)
		}

		v, ok := emptyDb.Voters[theCount.VoterId]
		if !ok {
			t.Error("Put a voter in the DB but it didn't stay")
		}

		if v.VoterId != theCount.VoterId || v.Name != theCount.Name {
			t.Errorf("Expected %v, got %v", theCount, v)
		}
	})
}
