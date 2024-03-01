package tests

import (
	"fmt"
	"testing"
	"time"

	"drexel.edu/voter-api/api"
	"drexel.edu/voter-api/db"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

var (
	BASE_API = "http://localhost:1080"
	cli      = resty.New()

	testVoters []db.Voter = []db.Voter{
		{
			VoterId: 1,
			Name:    "Count Chocula",
			VoteHistory: []db.VoterHistory{
				{PollId: 1, VoteId: 1, VoteDate: time.Date(2020, time.March, 17, 15, 00, 00, 00, time.UTC)},
			},
		},
		{
			VoterId:     2,
			Name:        "Captain Crunch",
			VoteHistory: make([]db.VoterHistory, 0),
		},
		{
			VoterId:     3,
			Name:        "Tony the Tiger",
			VoteHistory: make([]db.VoterHistory, 0),
		},
	}
)

func voterUrl(v db.Voter) string {
	return fmt.Sprintf("%s/voters/%d", BASE_API, v.VoterId)
}

func voterUrlById(vid uint) string {
	return fmt.Sprintf("%s/voters/%d", BASE_API, vid)
}

func voterHistoryUrlById(vid uint) string {
	return fmt.Sprintf("%s/voters/%d/polls", BASE_API, vid)
}

func voterPollUrlById(vid, pid uint) string {
	return fmt.Sprintf("%s/voters/%d/polls/%d", BASE_API, vid, pid)
}

var storedHealth api.HealthCheckResult

func Test_HealthBeforeActivity(t *testing.T) {
	var health api.HealthCheckResult
	rsp, err := cli.R().SetResult(&health).Get(BASE_API + "/voters/health")

	// store the results of this health check so we can test that the
	// counters increased at the end
	storedHealth = health

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())
	assert.Equal(t, "ok", health.Status)
}

func Test_SetupDb(t *testing.T) {
	rsp, err := cli.R().Delete(BASE_API + "/voters")
	if err != nil {
		t.Fatalf("can not clear database: %v\n", err)
	}

	if rsp.StatusCode() != 200 {
		t.Fatalf("bad status clearing db: %d\n", rsp.StatusCode())
	}

	for _, v := range testVoters {
		rsp, err := cli.R().SetBody(v).Post(voterUrl(v))
		if err != nil {
			t.Fatalf("can not initialize db: %v\n", err)
		}

		if rsp.StatusCode() >= 400 {
			t.Fatalf("bad status populating db: %d\n", rsp.StatusCode())
		}

	}
}

func Test_GetAllVoters(t *testing.T) {
	var voters []db.Voter
	rsp, err := cli.R().SetResult(&voters).Get(BASE_API + "/voters")

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())
	assert.Equal(t, 3, len(voters))
}

func Test_GetVoter(t *testing.T) {
	t.Run("GetVoter1", func(t *testing.T) {
		var voter db.Voter
		rsp, err := cli.R().SetResult(&voter).Get(voterUrlById(1))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())

		assert.Equal(t, voter, testVoters[0])
	})

	t.Run("GetVoter2", func(t *testing.T) {
		var voter db.Voter
		rsp, err := cli.R().SetResult(&voter).Get(voterUrlById(2))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())

		assert.Equal(t, voter, testVoters[1])
	})

	t.Run("GetVoter3", func(t *testing.T) {
		var voter db.Voter
		rsp, err := cli.R().SetResult(&voter).Get(voterUrlById(3))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())

		assert.Equal(t, voter, testVoters[2])
	})

	t.Run("GetNonExistingVoter", func(t *testing.T) {
		var voter db.Voter
		rsp, err := cli.R().SetResult(&voter).Get(voterUrlById(999))

		assert.Nil(t, err)
		assert.Equal(t, 404, rsp.StatusCode())
	})

	t.Run("GetVoterNegative1", func(t *testing.T) {
		var voter db.Voter
		rsp, err := cli.R().SetResult(&voter).Get(BASE_API + "/voters/-1")

		assert.Nil(t, err)
		assert.Equal(t, 404, rsp.StatusCode())
	})
}

func Test_GetVoterHistory(t *testing.T) {
	t.Run("GetVoter1", func(t *testing.T) {
		var vh []db.VoterHistory
		rsp, err := cli.R().SetResult(&vh).Get(voterHistoryUrlById(1))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())

		assert.Equal(t, vh, testVoters[0].VoteHistory)
	})

	t.Run("GetVoter2", func(t *testing.T) {
		var vh []db.VoterHistory
		rsp, err := cli.R().SetResult(&vh).Get(voterHistoryUrlById(2))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())

		assert.Equal(t, vh, testVoters[1].VoteHistory)
	})

	t.Run("GetNonExistingVoter", func(t *testing.T) {
		var vh []db.VoterHistory
		rsp, err := cli.R().SetResult(&vh).Get(voterHistoryUrlById(999))

		assert.Nil(t, err)
		assert.Equal(t, 404, rsp.StatusCode())
	})
}

func Test_GetVoterHistoryPoll(t *testing.T) {
	t.Run("GetVoter1Poll1", func(t *testing.T) {
		var vh db.VoterHistory
		rsp, err := cli.R().SetResult(&vh).Get(voterPollUrlById(1, 1))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())

		assert.Equal(t, vh, testVoters[0].VoteHistory[0])
	})

	t.Run("VoterExistsPollDoesNot", func(t *testing.T) {
		var vh []db.VoterHistory
		rsp, err := cli.R().SetResult(&vh).Get(voterPollUrlById(2, 1))

		assert.Nil(t, err)
		assert.Equal(t, 404, rsp.StatusCode())
	})

	t.Run("GetNonExistingVoter", func(t *testing.T) {
		var vh []db.VoterHistory
		rsp, err := cli.R().SetResult(&vh).Get(voterPollUrlById(999, 1))

		assert.Nil(t, err)
		assert.Equal(t, 404, rsp.StatusCode())
	})
}

func Test_AddVoter(t *testing.T) {
	voter4 := db.Voter{
		VoterId: 4,
		Name:    "Lucky the Leprechaun",
		Email:   "lucky@luckycharms.com",
		VoteHistory: []db.VoterHistory{
			{PollId: 1, VoteId: 1, VoteDate: time.Date(2020, time.March, 17, 15, 00, 00, 00, time.UTC)},
			{PollId: 2, VoteId: 1, VoteDate: time.Date(2021, time.March, 17, 15, 00, 00, 00, time.UTC)},
			{PollId: 3, VoteId: 1, VoteDate: time.Date(2022, time.March, 17, 15, 00, 00, 00, time.UTC)},
		},
	}

	t.Run("AddVoter4", func(t *testing.T) {
		rsp, err := cli.R().SetBody(voter4).Post(voterUrl(voter4))

		assert.Nil(t, err)
		assert.Equal(t, 201, rsp.StatusCode())
	})

	t.Run("DuplicateVoter4", func(t *testing.T) {
		rsp, err := cli.R().SetBody(voter4).Post(voterUrl(voter4))

		assert.Nil(t, err)
		assert.Equal(t, 500, rsp.StatusCode())
	})

	t.Run("InvalidVoter4", func(t *testing.T) {
		rsp, err := cli.R().SetBody("this is not a voter").Post(voterUrl(voter4))

		assert.Nil(t, err)
		assert.Equal(t, 400, rsp.StatusCode())
	})

	t.Run("ReadBackVoter4", func(t *testing.T) {
		var voter db.Voter
		rsp, err := cli.R().SetResult(&voter).Get(voterUrlById(4))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())

		assert.EqualValues(t, voter, voter4)
	})
}

func Test_AddVoterHistoryPoll(t *testing.T) {
	newPoll := db.VoterHistory{PollId: 2, VoteId: 4, VoteDate: time.Date(2000, time.January, 01, 00, 00, 00, 00, time.UTC)}

	t.Run("AddVoter1Poll2", func(t *testing.T) {
		rsp, err := cli.R().SetBody(newPoll).Post(voterPollUrlById(1, 2))

		assert.Nil(t, err)
		assert.Equal(t, 201, rsp.StatusCode())
	})

	t.Run("ReadBackVoter1Poll2", func(t *testing.T) {
		var poll db.VoterHistory
		rsp, err := cli.R().SetResult(&poll).Get(voterPollUrlById(1, 2))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
		assert.EqualValues(t, newPoll, poll)
	})

	t.Run("AddDuplicatePoll", func(t *testing.T) {
		rsp, err := cli.R().SetBody(newPoll).Post(voterPollUrlById(1, 2))

		assert.Nil(t, err)
		assert.Equal(t, 500, rsp.StatusCode())
	})

	t.Run("AddPollToNonExistentUser", func(t *testing.T) {
		rsp, err := cli.R().SetBody(newPoll).Get(voterPollUrlById(999, 2))

		assert.Nil(t, err)
		assert.Equal(t, 404, rsp.StatusCode())
	})
}

func Test_UpdateVoter(t *testing.T) {
	voter4 := db.Voter{
		VoterId:     4,
		Name:        "Trix Rabbit",
		Email:       "rabbit@trix.com",
		VoteHistory: []db.VoterHistory{},
	}

	t.Run("UpdateVoter4", func(t *testing.T) {
		rsp, err := cli.R().SetBody(voter4).Put(voterUrl(voter4))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	})

	t.Run("NonExistentUser", func(t *testing.T) {
		rsp, err := cli.R().SetBody(voter4).Put(voterUrlById(999))

		assert.Nil(t, err)
		assert.Equal(t, 400, rsp.StatusCode())
	})

	t.Run("DuplicateVoter4", func(t *testing.T) {
		rsp, err := cli.R().SetBody(voter4).Put(voterUrl(voter4))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	})

	t.Run("InvalidVoter4", func(t *testing.T) {
		rsp, err := cli.R().SetBody("this is not a voter").Put(voterUrl(voter4))

		assert.Nil(t, err)
		assert.Equal(t, 400, rsp.StatusCode())
	})

	t.Run("ReadBackVoter4WithHistory", func(t *testing.T) {
		var voter db.Voter
		rsp, err := cli.R().SetResult(&voter).Get(voterUrlById(4))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())

		// Ensure that UpdateVoter affects the voter details only and does not affect
		// the voting history
		assert.Equal(t, voter.VoterId, voter4.VoterId)
		assert.Equal(t, voter.Name, voter4.Name)
		assert.Equal(t, voter.Email, voter4.Email)
		assert.NotEqualValues(t, voter.VoteHistory, voter4.VoteHistory)
	})
}

func Test_UpdateVoterHistoryPoll(t *testing.T) {
	newPoll := db.VoterHistory{PollId: 2, VoteId: 9, VoteDate: time.Date(2024, time.January, 01, 13, 15, 17, 00, time.UTC)}

	t.Run("UpdateVoter1Poll2", func(t *testing.T) {
		rsp, err := cli.R().SetBody(newPoll).Put(voterPollUrlById(1, 2))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	})

	t.Run("ReadBackVoter1Poll2", func(t *testing.T) {
		var poll db.VoterHistory
		rsp, err := cli.R().SetResult(&poll).Get(voterPollUrlById(1, 2))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
		assert.EqualValues(t, newPoll, poll)
	})

	t.Run("UpdatePollForNonExistentUser", func(t *testing.T) {
		rsp, err := cli.R().SetBody(newPoll).Put(voterPollUrlById(999, 2))

		assert.Nil(t, err)
		assert.Equal(t, 500, rsp.StatusCode())
	})

	t.Run("UserExistsButPollDoesNot", func(t *testing.T) {
		rsp, err := cli.R().SetBody(newPoll).Put(voterPollUrlById(1, 999))

		assert.Nil(t, err)
		assert.Equal(t, 500, rsp.StatusCode())
	})
}

func Test_DeleteVoter(t *testing.T) {
	t.Run("NonExistentVoter", func(t *testing.T) {
		rsp, err := cli.R().Delete(voterUrlById(999))

		assert.Nil(t, err)
		assert.Equal(t, 500, rsp.StatusCode())
	})

	t.Run("DeleteVoter3", func(t *testing.T) {
		rsp, err := cli.R().Delete(voterUrlById(3))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	})

	t.Run("DeleteUser3Again", func(t *testing.T) {
		rsp, err := cli.R().Delete(voterUrlById(3))

		assert.Nil(t, err)
		assert.Equal(t, 500, rsp.StatusCode())
	})
}

func Test_DeleteVoterHistoryPoll(t *testing.T) {
	t.Run("DeleteVoter1Poll2", func(t *testing.T) {
		rsp, err := cli.R().Delete(voterPollUrlById(1, 2))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	})

	t.Run("ReadBackVoter1Poll2", func(t *testing.T) {
		var poll db.VoterHistory
		rsp, err := cli.R().SetResult(&poll).Get(voterPollUrlById(1, 2))

		assert.Nil(t, err)
		assert.Equal(t, 404, rsp.StatusCode())
	})

	t.Run("DeletePollForNonExistentUser", func(t *testing.T) {
		rsp, err := cli.R().Delete(voterPollUrlById(999, 2))

		assert.Nil(t, err)
		assert.Equal(t, 500, rsp.StatusCode())
	})

	t.Run("UserExistsButPollDoesNot", func(t *testing.T) {
		rsp, err := cli.R().Delete(voterPollUrlById(1, 999))

		assert.Nil(t, err)
		assert.Equal(t, 500, rsp.StatusCode())
	})
}

func Test_HealthAfterActivity(t *testing.T) {
	var health api.HealthCheckResult
	rsp, err := cli.R().SetResult(&health).Get(BASE_API + "/voters/health")

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())
	assert.Equal(t, "ok", health.Status)

	assert.Equal(t, health.Version, storedHealth.Version)
	assert.GreaterOrEqual(t, health.Uptime, storedHealth.Uptime)
	assert.Greater(t, health.Transactions, storedHealth.Transactions)
	assert.Greater(t, health.Errors, storedHealth.Errors)
}
