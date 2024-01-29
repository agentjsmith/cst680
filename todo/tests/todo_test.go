package tests

//Introduction to testing.  Note that testing is built into go and we will be using
//it extensively in this class. Below is a starter for your testing code.  In
//addition to what is built into go, we will be using a few third party packages
//that improve the testing experience.  The first is testify.  This package brings
//asserts to the table, that is much better than directly interacting with the
//testing.T object.  Second is gofakeit.  This package provides a significant number
//of helper functions to generate random data to make testing easier.

import (
	"fmt"
	"os"
	"testing"

	"drexel.edu/todo/db"
	fake "github.com/brianvoe/gofakeit/v6" //aliasing package name
	"github.com/stretchr/testify/assert"
)

// Note the default file path is relative to the test package location.  The
// project has a /tests path where you are at and a /data path where the
// database file sits.  So to get there we need to back up a directory and
// then go into the /data directory.  Thus this is why we are setting the
// default file name to "../data/todo.json"
const (
	DEFAULT_DB_FILE_NAME = "../data/todo.json"
)

var (
	DB *db.ToDo
)

// note init() is a helpful function in golang.  If it exists in a package
// such as we are doing here with the testing package, it will be called
// exactly once.  This is a great place to do setup work for your tests.
func init() {
	//Below we are setting up the gloabal DB variable that we can use in
	//all of our testing functions to make life easier
	testdb, err := db.New(DEFAULT_DB_FILE_NAME)
	if err != nil {
		fmt.Print("ERROR CREATING DB:", err)
		os.Exit(1)
	}

	DB = testdb //setup the global DB variable to support test cases

	//Now lets start with a fresh DB with the sample test data
	testdb.RestoreDB()
}

// Sample Test, will always pass, comparing the second parameter to true, which
// is hard coded as true
func TestTrue(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestAddHardCodedItem(t *testing.T) {
	item := db.ToDoItem{
		Id:     999,
		Title:  "This is a test case item",
		IsDone: false,
	}
	t.Log("Testing Adding a Hard Coded Item: ", item)

	//DONE: finish this test, add an item to the database and then
	//check that it was added correctly by looking it back up
	//use assert.NoError() to ensure errors are not returned.
	//explore other useful asserts in the testify package, see
	//https://github.com/stretchr/testify.  Specifically look
	//at things like assert.Equal() and assert.Condition()

	//I will get you started, uncomment the lines below to add to the DB
	//and ensure no errors:
	//---------------------------------------------------------------
	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to DB")

	//DONE: Now finish the test case by looking up the item in the DB
	//and making sure it matches the item that you put in the DB above

	newItem, err := DB.GetItem(999)
	assert.NoError(t, err, "Error fetching item from DB")
	assert.Equal(t,item,newItem,"Item fetched from DB must match item inserted")
}

func TestAddRandomStructItem(t *testing.T) {
	//You can also use the Stuct() fake function to create a random struct
	//Not going to do anyting
	item := db.ToDoItem{}
	err := fake.Struct(&item)
	t.Log("Testing Adding a Randomly Generated Struct: ", item)

	assert.NoError(t, err, "Created fake item OK")

	//DONE: Complete the test
	err = DB.AddItem(item)
	assert.NoError(t, err, "Added fake item to DB")

	newItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Fetched item from DB")
	assert.Equal(t,item,newItem,"Item fetched from DB matches item inserted")
}

func TestAddRandomItem(t *testing.T) {
	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	t.Log("Testing Adding an Item with Random Fields: ", item)

	err := DB.AddItem(item)
	assert.NoError(t, err, "Added fake item to DB")

	newItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Fetched item from DB")
	assert.Equal(t,item,newItem,"Item fetched from DB matches item inserted")
}

//DONE: Create additional tests to showcase the correct operation of your program
//for example getting an item, getting all items, updating items, and so on. Be
//creative here.
func TestAddDuplicateItem(t *testing.T) {
	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	t.Log("Testing Adding a duplicate item", item)

	err := DB.AddItem(item)
	assert.NoError(t, err, "Added first item to DB")

	err = DB.AddItem(item)
	assert.Error(t, err, "Adding duplicate item should fail")
}

func TestUpdateItem(t *testing.T) {
	item,err := DB.GetItem(2)
	assert.NoError(t, err, "Fetch an item from the DB")

	newTitle := "Test updating items."
	newStatus := true
	
	item.Title = newTitle
	item.IsDone = newStatus

	err = DB.UpdateItem(item)
	assert.NoError(t, err, "Updating DB with new item")

	checkItem, err := DB.GetItem(2)
	assert.NoError(t, err, "After update, item should still be present")
	assert.Equal(t, item, checkItem, "All fields should match what we updated")
}

func TestChangeItemDoneState(t *testing.T) {
	err := DB.ChangeItemDoneStatus(1,true)
	assert.NoError(t, err, "Change item 1 done state")

	checkItem, err := DB.GetItem(1)
	assert.NoError(t, err, "Fetch item 1 from DB")
	assert.Equal(t, true, checkItem.IsDone, "Item should be done now")
}

func TestDeleteItem(t *testing.T) {
	err := DB.DeleteItem(1)
	assert.NoError(t, err, "Delete item 1")

	_, err = DB.GetItem(1)
	assert.Error(t, err, "Fetch item 1 from DB should fail now")
}

func TestGetAllItems(t *testing.T) {
	items,err := DB.GetAllItems()
	assert.NoError(t, err, "Getting all items should not fail")
	assert.GreaterOrEqual(t, len(items), 4, "At least 4 items should have been returned")
}
