package api

import (
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"

	"drexel.edu/voter-api/db"
	"github.com/gofiber/fiber/v2"
)

// The api package creates and maintains a reference to the data handler
// this is a good design practice
type VoterAPI struct {
	db           *db.VoterDB
	bootTime     time.Time
	transactions uint
	errors       uint
}

func New(redisClient *redis.Client) (*VoterAPI, error) {
	dbHandler, err := db.New(redisClient)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &VoterAPI{db: dbHandler, bootTime: now}, nil
}

func (va *VoterAPI) GetAllVoters(c *fiber.Ctx) error {
	va.transactions++

	voterList, err := va.db.GetAllVoters(c.Context())
	if err != nil {
		va.errors++
		log.Println("Error Getting All Voters: ", err)
		return fiber.NewError(http.StatusNotFound,
			"Error Getting All Voters")
	}

	return c.JSON(voterList)
}

func (va *VoterAPI) GetVoter(c *fiber.Ctx) error {
	va.transactions++

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	id, err := c.ParamsInt("id")
	if err != nil {
		va.errors++
		return fiber.NewError(http.StatusBadRequest)
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	voter, err := va.db.GetVoter(c.Context(), uint(id))
	if err != nil {
		va.errors++
		log.Println("Voter not found: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	return c.JSON(voter)
}

func (va *VoterAPI) GetVoterHistory(c *fiber.Ctx) error {
	va.transactions++

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	id, err := c.ParamsInt("id")
	if err != nil {
		va.errors++
		return fiber.NewError(http.StatusBadRequest)
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	voter, err := va.db.GetVoter(c.Context(), uint(id))
	if err != nil {
		va.errors++
		log.Println("Voter not found: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	return c.JSON(voter.VoteHistory)
}

func (va *VoterAPI) GetVoterHistoryPoll(c *fiber.Ctx) error {
	va.transactions++

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	id, err := c.ParamsInt("id")
	pollId, err2 := c.ParamsInt("pollid")
	if err != nil || err2 != nil {
		va.errors++
		return fiber.NewError(http.StatusBadRequest)
	}

	pollHistory, err := va.db.GetHistoryByPollId(c.Context(), uint(id), uint(pollId))
	if err != nil {
		va.errors++
		log.Println("Voter history not found: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	return c.JSON(pollHistory)
}

func (va *VoterAPI) AddVoterHistoryPoll(c *fiber.Ctx) error {
	va.transactions++

	id, err := c.ParamsInt("id")
	pollId, err2 := c.ParamsInt("pollid")
	if err != nil || err2 != nil {
		va.errors++
		return fiber.NewError(http.StatusBadRequest)
	}

	var newPollHistory db.VoterHistory

	if err = c.BodyParser(&newPollHistory); err != nil {
		va.errors++
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	if uint(pollId) != newPollHistory.PollId {
		va.errors++
		log.Printf("Duplicate poll id %d for voter %d", pollId, id)
		return fiber.NewError(http.StatusBadRequest)
	}

	result, err := va.db.AddHistoryByPollId(c.Context(), uint(id), uint(pollId), newPollHistory)
	if err != nil {
		va.errors++
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusCreated).JSON(result)

}

func (va *VoterAPI) UpdateVoterHistoryPoll(c *fiber.Ctx) error {
	va.transactions++

	id, err := c.ParamsInt("id")
	pollId, err2 := c.ParamsInt("pollid")
	if err != nil || err2 != nil {
		va.errors++
		return fiber.NewError(http.StatusBadRequest)
	}

	var newPollHistory db.VoterHistory
	if err = c.BodyParser(&newPollHistory); err != nil {
		va.errors++
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	history, err := va.db.UpdateHistoryByPollId(c.Context(), uint(id), uint(pollId), newPollHistory)
	if err != nil {
		va.errors++
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(history)

}

func (va *VoterAPI) DeleteVoterHistoryPoll(c *fiber.Ctx) error {
	va.transactions++

	id, err := c.ParamsInt("id")
	pollId, err2 := c.ParamsInt("pollid")
	if err != nil || err2 != nil {
		va.errors++
		return fiber.NewError(http.StatusBadRequest)
	}

	err = va.db.DeleteHistoryByPollId(c.Context(), uint(id), uint(pollId))
	if err != nil {
		va.errors++
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(struct{}{})

}

func (va *VoterAPI) AddVoter(c *fiber.Ctx) error {
	va.transactions++

	var voter db.Voter
	if err := c.BodyParser(&voter); err != nil {
		va.errors++
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		va.errors++
		log.Println("id param missing or not an int")
		return fiber.NewError(http.StatusBadRequest)
	}

	if uint(id) != voter.VoterId {
		va.errors++
		log.Println("id param does not match payload")
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := va.db.AddVoter(c.Context(), voter); err != nil {
		va.errors++
		log.Println("Error adding item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(voter)
}

func (va *VoterAPI) UpdateVoter(c *fiber.Ctx) error {
	va.transactions++

	var voter db.Voter
	if err := c.BodyParser(&voter); err != nil {
		va.errors++
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		va.errors++
		log.Println("id param missing or not an int")
		return fiber.NewError(http.StatusBadRequest)
	}

	if uint(id) != voter.VoterId {
		va.errors++
		log.Println("id param does not match payload")
		return fiber.NewError(http.StatusBadRequest)
	}

	// This function is supposed to update the voter details only,
	// not the history, so we need to save the old history and use
	// it to replace whatever was passed in.
	oldVoter, err := va.db.GetVoter(c.Context(), uint(id))
	if err != nil {
		va.errors++
		log.Println("User not found for update: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	voter.VoteHistory = oldVoter.VoteHistory

	if err := va.db.UpdateVoter(c.Context(), voter); err != nil {
		va.errors++
		log.Println("Error updating item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(voter)
}

func (va *VoterAPI) DeleteVoter(c *fiber.Ctx) error {
	va.transactions++

	id, err := c.ParamsInt("id")
	if err != nil {
		va.errors++
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := va.db.DeleteVoter(c.Context(), uint(id)); err != nil {
		va.errors++
		log.Println("Error deleting item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).SendString("Delete OK")
}

func (va *VoterAPI) DeleteAllVoters(c *fiber.Ctx) error {
	va.transactions++

	if err := va.db.DeleteAll(c.Context()); err != nil {
		va.errors++
		log.Println("Error deleting all items: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).SendString("I hope you meant to do that!")
}

// implementation of GET /health. It is a good practice to build in a
// health check for your API.

type HealthCheckResult struct {
	Status       string `json:"status"`
	Version      string `json:"version"`
	Uptime       uint   `json:"uptime_seconds"`
	Transactions uint   `json:"transaction_count"`
	Errors       uint   `json:"error_count"`
	DbHealth     string `json:"database_status"`
}

func (va *VoterAPI) HealthCheck(c *fiber.Ctx) error {
	dbh := va.db.HealthCheck(c.Context())

	return c.Status(http.StatusOK).
		JSON(HealthCheckResult{
			Status:       "ok",
			Version:      "1.0.0",
			Uptime:       uint(time.Now().Sub(va.bootTime).Seconds()),
			Transactions: va.transactions,
			Errors:       va.errors,
			DbHealth:     dbh,
		})
}
