package api

import (
	"log"
	"net/http"

	"drexel.edu/voter-api/db"
	"github.com/gofiber/fiber/v2"
)

// The api package creates and maintains a reference to the data handler
// this is a good design practice
type VoterAPI struct {
	db *db.VoterList
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.New()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{db: dbHandler}, nil
}

// implementation for GET /todo
// returns all todos
func (va *VoterAPI) GetAllVoters(c *fiber.Ctx) error {

	todoList, err := va.db.GetAllVoters()
	if err != nil {
		log.Println("Error Getting All Voters: ", err)
		return fiber.NewError(http.StatusNotFound,
			"Error Getting All Voters")
	}
	//Note that the database returns a nil slice if there are no items
	//in the database.  We need to convert this to an empty slice
	//so that the JSON marshalling works correctly.  We want to return
	//an empty slice, not a nil slice. This will result in the json being []
	if todoList == nil {
		todoList = make([]db.Voter, 0)
	}

	return c.JSON(todoList)
}

// implementation for GET /todo/:id
// returns a single todo
func (va *VoterAPI) GetVoter(c *fiber.Ctx) error {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	voter, err := va.db.GetVoter(uint(id))
	if err != nil {
		log.Println("Voter not found: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	return c.JSON(voter)
}

func (va *VoterAPI) GetVoterHistory(c *fiber.Ctx) error {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	voter, err := va.db.GetVoter(uint(id))
	if err != nil {
		log.Println("Voter not found: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	return c.JSON(voter.VoteHistory)
}

func (va *VoterAPI) GetVoterHistoryPoll(c *fiber.Ctx) error {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	id, err := c.ParamsInt("id")
	poll_id, err2 := c.ParamsInt("pollid")
	if err != nil || err2 != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	pollHistory, err := va.db.GetHistoryByPollId(uint(id), uint(poll_id))
	if err != nil {
		log.Println("Voter history not found: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	return c.JSON(pollHistory)
}

func (va *VoterAPI) AddVoterHistoryPoll(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	poll_id, err2 := c.ParamsInt("pollid")
	if err != nil || err2 != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	var newPollHistory db.VoterHistory

	if err = c.BodyParser(&newPollHistory); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	result, err := va.db.AddHistoryByPollId(uint(id), uint(poll_id), newPollHistory)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError)
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	return c.Status(http.StatusCreated).JSON(result)

}

func (va *VoterAPI) UpdateVoterHistoryPoll(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	poll_id, err2 := c.ParamsInt("pollid")
	if err != nil || err2 != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	var newPollHistory db.VoterHistory
	if err = c.BodyParser(&newPollHistory); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	history, err := va.db.UpdateHistoryByPollId(uint(id), uint(poll_id), newPollHistory)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError)
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	return c.Status(http.StatusOK).JSON(history)

}

func (va *VoterAPI) DeleteVoterHistoryPoll(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	poll_id, err2 := c.ParamsInt("pollid")
	if err != nil || err2 != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	err = va.db.DeleteHistoryByPollId(uint(id), uint(poll_id))
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(struct{}{})

}

// implementation for POST /todo
// adds a new todo
func (va *VoterAPI) AddVoter(c *fiber.Ctx) error {
	var voter db.Voter
	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Println("id param missing or not an int")
		return fiber.NewError(http.StatusBadRequest)
	}

	if uint(id) != voter.VoterId {
		log.Println("id param does not match payload")
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := va.db.AddVoter(voter); err != nil {
		log.Println("Error adding item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(voter)
}

// implementation for PUT /todo
// Web api standards use PUT for Updates
func (va *VoterAPI) UpdateVoter(c *fiber.Ctx) error {
	var voter db.Voter
	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Println("id param missing or not an int")
		return fiber.NewError(http.StatusBadRequest)
	}

	if uint(id) != voter.VoterId {
		log.Println("id param does not match payload")
		return fiber.NewError(http.StatusBadRequest)
	}

	// This function is supposed to update the voter details only,
	// not the history, so we need to save the old history and use
	// it to replace whatever was passed in.
	oldVoter, err := va.db.GetVoter(uint(id))
	if err != nil {
		log.Println("User not found for update: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	voter.VoteHistory = oldVoter.VoteHistory

	if err := va.db.UpdateVoter(voter); err != nil {
		log.Println("Error updating item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(voter)
}

// implementation for DELETE /todo/:id
// deletes a todo
func (va *VoterAPI) DeleteVoter(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := va.db.DeleteVoter(uint(id)); err != nil {
		log.Println("Error deleting item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).SendString("Delete OK")
}

// implementation for DELETE /todo
// deletes all todos
func (va *VoterAPI) DeleteAllVoters(c *fiber.Ctx) error {

	if err := va.db.DeleteAll(); err != nil {
		log.Println("Error deleting all items: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).SendString("I hope you meant to do that!")
}

// implementation of GET /health. It is a good practice to build in a
// health check for your API.  Below the results are just hard coded
// but in a real API you can provide detailed information about the
// health of your API with a Health Check
// TODO: make values realistic for extra credit

type HealthCheckResult struct {
	Status       string `json:"status"`
	Version      string `json:"version"`
	Uptime       uint   `json:"uptime"`
	Transactions uint   `json:"transaction_count"`
	Errors       uint   `json:"error_count"`
}

func (va *VoterAPI) HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).
		JSON(HealthCheckResult{
			Status:       "ok",
			Version:      "1.0.0",
			Uptime:       100,
			Transactions: 1000,
			Errors:       10,
		})
}
