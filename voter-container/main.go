package main

import (
	"drexel.edu/voter-api/api"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

// Global variables to hold the command line flags to drive the voter CLI
// application
var (
	redisAddr  string
	listenAddr string
)

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func processEnvVars() {
	redisAddr = getEnvOrDefault("VOTER_API_REDIS_ADDR", "0.0.0.0:6379")
	listenAddr = getEnvOrDefault("VOTER_API_LISTEN_ADDR", "0.0.0.0:1080")
}
func addRoutes(app *fiber.App, apiHandler *api.VoterAPI) {
	//HTTP Standards for "REST" APIS
	//GET - Read/Query
	//POST - Create
	//PUT - Update
	//DELETE - Delete

	app.Get("/voters", apiHandler.GetAllVoters)
	app.Delete("/voters", apiHandler.DeleteAllVoters)

	app.Get("/voters/:id<int;min(0)>", apiHandler.GetVoter)
	app.Post("/voters/:id<int;min(0)>", apiHandler.AddVoter)
	app.Put("/voters/:id<int;min(0)>", apiHandler.UpdateVoter)
	app.Delete("/voters/:id<int;min(0)>", apiHandler.DeleteVoter)

	app.Get("/voters/:id<int;min(0)>/polls", apiHandler.GetVoterHistory)

	app.Get("/voters/:id<int;min(0)>/polls/:pollid<int;min(0)>", apiHandler.GetVoterHistoryPoll)
	app.Post("/voters/:id<int;min(0)>/polls/:pollid<int;min(0)>", apiHandler.AddVoterHistoryPoll)
	app.Put("/voters/:id<int;min(0)>/polls/:pollid<int;min(0)>", apiHandler.UpdateVoterHistoryPoll)
	app.Delete("/voters/:id<int;min(0)>/polls/:pollid<int;min(0)>", apiHandler.DeleteVoterHistoryPoll)

	app.Get("/voters/health", apiHandler.HealthCheck)

}

// main is the entry point for our voter API application.  It processes
// the command line flags and then uses the db package to perform the
// requested operation
func main() {
	processEnvVars()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New())

	log.Println("Connecting to Redis on ", redisAddr)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	apiHandler, err := api.New(redisClient)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	addRoutes(app, apiHandler)

	log.Println("Starting server on ", listenAddr)
	app.Listen(listenAddr)
}
