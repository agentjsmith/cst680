package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

type VoterList struct {
	Voters map[uint]Voter //A map of VoterIDs as keys and Voter structs as values
}

var voterList VoterList

// New is a constructor function that returns a pointer to a new
// VoterList struct.  It takes a single string argument that is the
// name of the file that will be used to store the VoterList items.
// If the file doesn't exist, it will be created.  If the file
// does exist, it will be loaded into the VoterList struct.
func New() (*VoterList, error) {

	voterList := &VoterList{
		Voters: make(map[uint]Voter),
	}

	return voterList, nil
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR TODO APP
//------------------------------------------------------------

// AddItem accepts a Voter and adds it to the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must not already exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if so, return an error
//
// Postconditions:
//
//	 (1) The item will be added to the DB
//		(2) The DB file will be saved with the item added
//		(3) If there is an error, it will be returned
func (vl *VoterList) AddVoter(item Voter) error {

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	_, ok := vl.Voters[item.VoterId]
	if ok {
		return errors.New("voter already exists")
	}

	//Now that we know the item doesn't exist, lets add it to our map
	vl.Voters[item.VoterId] = item

	//If everything is ok, return nil for the error
	return nil
}

// DeleteItem accepts an item id and removes it from the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be removed from the DB
//		(2) The DB file will be saved with the item removed
//		(3) If there is an error, it will be returned
func (vl *VoterList) DeleteVoter(id uint) error {

	// we should if item exists before trying to delete it
	// this is a good practice, return an error if the
	// item does not exist

	//Now lets use the built-in go delete() function to remove
	//the item from our map
	delete(vl.Voters, id)

	return nil
}

// DeleteAll removes all items from the DB.
// It will be exposed via a DELETE /todo endpoint
func (vl *VoterList) DeleteAll() error {
	//To delete everything, we can just create a new map
	//and assign it to our existing map.  The garbage collector
	//will clean up the old map for us
	vl.Voters = make(map[uint]Voter)

	return nil
}

// UpdateItem accepts a Voter and updates it in the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be updated in the DB
//		(2) The DB file will be saved with the item updated
//		(3) If there is an error, it will be returned
func (vl *VoterList) UpdateVoter(item Voter) error {

	// Check if item exists before trying to update it
	// this is a good practice, return an error if the
	// item does not exist
	_, ok := vl.Voters[item.VoterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	//Now that we know the item exists, lets update it
	vl.Voters[item.VoterId] = item

	return nil
}

func (vl *VoterList) GetHistoryByPollId(userId, pollId uint) (VoterHistory, error) {
	v, ok := vl.Voters[userId]
	if !ok {
		return VoterHistory{}, errors.New("voter does not exist")
	}

	for i := range v.VoteHistory {
		if v.VoteHistory[i].PollId == pollId {
			return v.VoteHistory[i], nil
		}
	}
	return VoterHistory{}, errors.New("poll not found in voter history")
}

func (vl *VoterList) AddHistoryByPollId(userId, pollId uint, newHistory VoterHistory) (VoterHistory, error) {
	v, ok := vl.Voters[userId]
	if !ok {
		return VoterHistory{}, errors.New("voter does not exist")
	}

	for i := range v.VoteHistory {
		if v.VoteHistory[i].PollId == pollId {
			return VoterHistory{}, errors.New("voter history already exists for that poll")
		}
	}

	v.VoteHistory = append(v.VoteHistory, newHistory)
	vl.Voters[userId] = v

	return newHistory, nil
}

func (vl *VoterList) UpdateHistoryByPollId(userId, pollId uint, newHistory VoterHistory) (VoterHistory, error) {
	v, ok := vl.Voters[userId]
	if !ok {
		return VoterHistory{}, errors.New("voter does not exist")
	}

	for i := range v.VoteHistory {
		if v.VoteHistory[i].PollId == pollId {
			vl.Voters[userId].VoteHistory[i] = newHistory
			return newHistory, nil
		}
	}

	return VoterHistory{}, errors.New("poll not found in voter history")
}

func (vl *VoterList) DeleteHistoryByPollId(userId, pollId uint) error {
	v, ok := vl.Voters[userId]
	if !ok {
		return errors.New("voter does not exist")
	}

	for i := range v.VoteHistory {
		if v.VoteHistory[i].PollId == pollId {
			newHistory := slices.Delete(v.VoteHistory, i, i+1)
			v.VoteHistory = newHistory
			vl.Voters[userId] = v

			return nil
		}
	}

	return errors.New("poll not found in voter history")
}

// GetItem accepts an item id and returns the item from the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be returned, if it exists
//		(2) If there is an error, it will be returned
//			along with an empty Voter
//		(3) The database file will not be modified
func (vl *VoterList) GetVoter(id uint) (Voter, error) {

	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	item, ok := vl.Voters[id]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}

	return item, nil
}

// GetAllItems returns all items from the DB.  If successful it
// returns a slice of all of the items to the caller
// Preconditions:   (1) The database file must exist and be a valid
//
// Postconditions:
//
//	 (1) All items will be returned, if any exist
//		(2) If there is an error, it will be returned
//			along with an empty slice
//		(3) The database file will not be modified
func (vl *VoterList) GetAllVoters() ([]Voter, error) {

	//Now that we have the DB loaded, lets crate a slice
	var voters []Voter

	//Now lets iterate over our map and add each item to our slice
	for _, item := range vl.Voters {
		voters = append(voters, item)
	}

	//Now that we have all of our items in a slice, return it
	return voters, nil
}

// PrintItem accepts a Voter and prints it to the console
// in a JSON pretty format. As some help, look at the
// json.MarshalIndent() function from our in class go tutorial.
func (vl *VoterList) PrintVoter(item Voter) {
	jsonBytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Println(string(jsonBytes))
}

// PrintAllItems accepts a slice of Voters and prints them to the console
// in a JSON pretty format.  It should call PrintItem() to print each item
// versus repeating the code.
func (vl *VoterList) PrintAllVoters(itemList []Voter) {
	for _, item := range itemList {
		vl.PrintVoter(item)
	}
}

// JsonToItem accepts a json string and returns a Voter
// This is helpful because the CLI accepts todo items for insertion
// and updates in JSON format.  We need to convert it to a Voter
// struct to perform any operations on it.
func (vl *VoterList) JsonToVoter(jsonString string) (Voter, error) {
	var item Voter
	err := json.Unmarshal([]byte(jsonString), &item)
	if err != nil {
		return Voter{}, err
	}

	return item, nil
}
