package tests

//Introduction to testing.  Note that testing is built into go and we will be using
//it extensively in this class. Below is a starter for your testing code.  In
//addition to what is built into go, we will be using a few third party packages
//that improve the testing experience.  The first is testify.  This package brings
//asserts to the table, that is much better than directly interacting with the
//testing.T object.  Second is gofakeit.  This package provides a significant number
//of helper functions to generate random data to make testing easier.

import (
	"encoding/json"
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
	DEFAULT_DB_FILE_NAME        = "../data/todo.json"
	DEFAULT_DB_BACKUP_FILE_NAME = "../data/todo.json.bak"
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

	//TODO: finish this test, add an item to the database and then
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

	//TODO: Now finish the test case by looking up the item in the DB
	//and making sure it matches the item that you put in the DB above

	fileContents, _ := DB.GetItem(item.Id)
	assert.Equal(t, item, fileContents, "found added item.id")
}

func TestAddRandomStructItem(t *testing.T) {
	//You can also use the Stuct() fake function to create a random struct
	//Not going to do anyting
	item := db.ToDoItem{}
	err := fake.Struct(&item)
	t.Log("Testing Adding a Randomly Generated Struct: ", item)

	assert.NoError(t, err, "Created fake item OK")

	//TODO: Complete the test
}

func TestAddRandomItem(t *testing.T) {
	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	t.Log("Testing Adding an Item with Random Fields: ", item)

}

//TODO: Create additional tests to showcase the correct operation of your program
//for example getting an item, getting all items, updating items, and so on. Be
//creative here.

// RestoreDB func test
func TestRestoreDB(t *testing.T) {

	assert.FileExists(t, DEFAULT_DB_BACKUP_FILE_NAME, "todo.json.back file in ../data does not exist")

	// remove the current db
	err := os.Remove(DEFAULT_DB_FILE_NAME)
	assert.NoError(t, err, "Found error while removing file in TestRestoreDB")
	assert.NoFileExists(t, DEFAULT_DB_FILE_NAME, "todo.json file in ../data exist")

	// use the restoreDB function and see if the file was created
	err = DB.RestoreDB()
	assert.NoError(t, err, "Found error while running RestoreDB in TestRestoreDB")
	assert.FileExists(t, DEFAULT_DB_FILE_NAME, "todo.json file in ../data does not exist")
}

func TestGetAllItems(t *testing.T) {
	items, err := DB.GetAllItems()
	var fileContents []db.ToDoItem
	data, _ := os.ReadFile(DEFAULT_DB_BACKUP_FILE_NAME)
	err = json.Unmarshal(data, &fileContents)

	assert.NoError(t, err, "Found error while running TestGetAllItems")
	assert.Equal(t, fileContents, items, "Did not find expected value.")
}

func TestDeleteItem(t *testing.T) {

	// add a new item.id
	item := db.ToDoItem{
		Id:     21,
		Title:  "try something new",
		IsDone: false,
	}
	err := DB.AddItem(item)
	assert.NoError(t, err, "Found error while running TestAddItem")

	// delete newly created item.id
	err = DB.DeleteItem(21)
	assert.NoError(t, err, "The func DeleteItem returned a error")

	// test to see if item.id exists
	err = DB.DeleteItem(21)
	assert.Error(t, err, "Did not pass because TestDeleteItem Id should return error")
}

func TestGetItem(t *testing.T) {

	// test an actual existing id
	id := 1
	data, err := DB.GetItem(id)
	assert.NoError(t, err, "Ran into error when running GetItem func")
	assert.Equal(t, id, data.Id, " Did not find matching ID during GetItem Run")

	// test a non-existing id, error should return as error
	id = 10
	_, err = DB.GetItem(id)
	assert.Error(t, err, "Ran into nil when running GetItem func")
}

func TestUpdateItem(t *testing.T) {

	item := db.ToDoItem{
		Id:     4,
		Title:  "Fury Max",
		IsDone: false,
	}

	err := DB.UpdateItem(item)
	if err != nil {
		return
	}

	var fileContents []db.ToDoItem
	data, _ := os.ReadFile(DEFAULT_DB_FILE_NAME)
	err = json.Unmarshal(data, &fileContents)

	assert.NoError(t, err, "Ran into error while running UpdateItem func")
	assert.Equal(t, "Fury Max", fileContents[3].Title, "Did not find matching update.")

}
