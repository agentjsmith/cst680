package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type VoterDB struct {
	redisClient *redis.Client
}

// New is a constructor function that returns a pointer to a new
// VoterDB struct.  It takes a single string argument that is the
// name of the file that will be used to store the VoterDB items.
// If the file doesn't exist, it will be created.  If the file
// does exist, it will be loaded into the VoterDB struct.
func New(redisClient *redis.Client) (*VoterDB, error) {

	voterList := &VoterDB{
		redisClient: redisClient,
	}

	// ensure the connection actually works
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis session failed: %w", err)
	}

	return voterList, nil
}

// returns the Redis key of a voter given the voter struct
func voterKey(v Voter) string {
	return idKey(v.VoterId)
}

// returns the Redis key of a voter given the voter id
func idKey(id uint) string {
	return fmt.Sprintf("voter:%d", id)
}

func wildcardKey() string {
	return "voter:*"
}

// returns a JSONPath expression to extract the given poll id
// from a voter document
func pollIdPath(id uint) string {
	return fmt.Sprintf("$.history[?(@.poll_id==%d)]", id)
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR VOTER APP
//------------------------------------------------------------

func (db *VoterDB) AddVoter(ctx context.Context, item Voter) error {
	key := voterKey(item)

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	oldVoter, err := db.redisClient.JSONGet(ctx, key, "$").Result()
	if oldVoter != "" {
		return errors.New("voter already exists")
	}
	if err != nil {
		return fmt.Errorf("checking duplicate voter: %w", err)
	}

	//Now that we know the item doesn't exist, lets add it to our map
	_, err = db.redisClient.JSONSet(ctx, key, "$", item).Result()
	if err != nil {
		return fmt.Errorf("add voter: %w", err)
	}

	//If everything is ok, return nil for the error
	return nil
}

func (db *VoterDB) DeleteVoter(ctx context.Context, id uint) error {
	key := idKey(id)

	// we should if item exists before trying to delete it
	// this is a good practice, return an error if the
	// item does not exist
	_, err := db.fetchVoter(ctx, key)
	if err != nil {
		return errors.New("voter not found")
	}

	_, err = db.redisClient.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("deleting voter: %w", err)
	}

	return nil
}

// DeleteAll removes all items from the DB.
func (db *VoterDB) DeleteAll(ctx context.Context) error {
	//To delete everything, we can just create a new map
	//and assign it to our existing map.  The garbage collector
	//will clean up the old map for us
	db.redisClient.FlushDB(ctx)

	return nil
}

func (db *VoterDB) UpdateVoter(ctx context.Context, item Voter) error {
	key := voterKey(item)

	// Check if item exists before trying to update it
	// this is a good practice, return an error if the
	// item does not exist
	_, err := db.fetchVoter(ctx, key)
	if err != nil {
		return errors.New("voter does not exist")
	}

	//Now that we know the item exists, lets update it
	_, err = db.redisClient.JSONSet(ctx, key, "$", item).Result()
	if err != nil {
		return fmt.Errorf("updating voter: %w", err)
	}

	return nil
}

func (db *VoterDB) fetchHistory(ctx context.Context, userId, pollId uint) (VoterHistory, error) {
	key := idKey(userId)

	value, err := db.redisClient.JSONGet(ctx, key, pollIdPath(pollId)).Result()
	if err != nil {
		return VoterHistory{}, fmt.Errorf("get history by poll id: %w", err)
	}

	var vh []VoterHistory
	err = json.Unmarshal([]byte(value), &vh)
	if err != nil {
		return VoterHistory{}, fmt.Errorf("unmarshaling vote history: %w", err)
	}

	if len(vh) <= 0 {
		return VoterHistory{}, errors.New("poll not found in vote history")
	}
	return vh[0], nil
}

func (db *VoterDB) GetHistoryByPollId(ctx context.Context, userId, pollId uint) (VoterHistory, error) {
	return db.fetchHistory(ctx, userId, pollId)
}

func (db *VoterDB) AddHistoryByPollId(ctx context.Context, userId, pollId uint, newHistory VoterHistory) (VoterHistory, error) {
	key := idKey(userId)

	// ensure this history does not already exist
	_, err := db.fetchHistory(ctx, userId, pollId)
	if err == nil {
		return VoterHistory{}, errors.New("history already exists")
	}

	newHistoryJson, err := json.Marshal(newHistory)
	if err != nil {
		return VoterHistory{}, fmt.Errorf("marshalling history entry: %w", err)
	}

	_, err = db.redisClient.JSONArrAppend(ctx, key, "$.history", newHistoryJson).Result()
	if err != nil {
		return VoterHistory{}, fmt.Errorf("adding history by poll id: %w", err)
	}

	return newHistory, nil
}

func (db *VoterDB) UpdateHistoryByPollId(ctx context.Context, userId, pollId uint, newHistory VoterHistory) (VoterHistory, error) {
	key := idKey(userId)
	path := pollIdPath(pollId)

	// ensure this history exists
	_, err := db.fetchHistory(ctx, userId, pollId)
	if err != nil {
		return VoterHistory{}, errors.New("history does not exist")
	}

	_, err = db.redisClient.JSONSet(ctx, key, path, newHistory).Result()
	if err != nil {
		return VoterHistory{}, fmt.Errorf("updating history by poll id: %w", err)
	}

	return newHistory, nil
}

func (db *VoterDB) DeleteHistoryByPollId(ctx context.Context, userId, pollId uint) error {
	key := idKey(userId)
	path := pollIdPath(pollId)

	// ensure this history exists
	_, err := db.fetchHistory(ctx, userId, pollId)
	if err != nil {
		return errors.New("history does not exist")
	}

	_, err = db.redisClient.JSONDel(ctx, key, path).Result()
	if err != nil {
		return fmt.Errorf("deleting history by poll id: %w", err)
	}

	return nil
}

func (db *VoterDB) fetchVoter(ctx context.Context, key string) (Voter, error) {
	value, err := db.redisClient.JSONGet(ctx, key, "$").Result()
	if err != nil {
		return Voter{}, fmt.Errorf("get voter by id: %w", err)
	}

	var v []Voter
	err = json.Unmarshal([]byte(value), &v)
	if err != nil {
		return Voter{}, fmt.Errorf("unmarshaling voter: %w", err)
	}

	return v[0], nil
}

func (db *VoterDB) GetVoter(ctx context.Context, id uint) (Voter, error) {
	key := idKey(id)
	return db.fetchVoter(ctx, key)
}

// GetAllVoters returns a list of every registered voter
// Warning! runs N+1 gets
func (db *VoterDB) GetAllVoters(ctx context.Context) ([]Voter, error) {
	allKeys, err := db.redisClient.Keys(ctx, wildcardKey()).Result()
	if err != nil {
		return nil, fmt.Errorf("getting all voter keys: %w", err)
	}

	voters := make([]Voter, 0, len(allKeys))

	// fetch the JSON object for every matching key and append them to a list
	// to be returned to the caller
	for _, key := range allKeys {
		voter, err := db.fetchVoter(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("getting all voters (voter %s): %w", key, err)
		}
		voters = append(voters, voter)
	}

	//Now that we have all of our items in a slice, return it
	return voters, nil
}

func (db *VoterDB) HealthCheck(ctx context.Context) string {
	_, err := db.redisClient.Ping(ctx).Result()
	if err != nil {
		return err.Error()
	}
	return "ok"
}
