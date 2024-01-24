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
	"errors"
	"fmt"
	"io"
	"log"
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

	//TODO: finish this test, add an item to the database and then
	//check that it was added correctly by looking it back up
	//use assert.NoError() to ensure errors are not returned.
	//explore other useful asserts in the testify package, see
	//https://github.com/stretchr/testify.  Specifically look
	//at things like assert.Equal() and assert.Condition()

	//I will get you started, uncomment the lines below to add to the DB
	//and ensure no errors:
	//---------------------------------------------------------------
	//err := DB.AddItem(item)
	//assert.NoError(t, err, "Error adding item to DB")

	//TODO: Now finish the test case by looking up the item in the DB
	//and making sure it matches the item that you put in the DB above

	data, err := os.ReadFile("../data/todo.json")
	if err != nil {
		log.Fatalln(err)
	}

	var dataList []db.ToDoItem

	err = json.Unmarshal(data, &dataList)
	if err != nil {
		log.Fatalln(err)
	}

	for _, i := range dataList {
		if item.Id == i.Id {
			err = errors.New("shit broken")
		}
	}

	append(dataList, item)

	jsonString, _ := json.Marshal(dataList)

	err := os.WriteFile(jsonString, , 0644)

	assert.Equal(t, 999, item.Id, "item.Id did not match with expected value.")
	assert.Equal(t, "This is a test case item", item.Title, "item.Title did not match with expected value.")
	assert.Equal(t, false, item.IsDone, "item.IsDone did not match with expected value.")
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

func TestRestoreDatabase(t *testing.T) {

	// define file paths
	primaryFile := "../data/todo.json"
	backupFile := "../data/todo.json.bak"

	// remove the file
	err := os.Remove(primaryFile)

	// test if there were any errors
	assert.NoError(t, err, "ToDo.json Exists.")

	// read the backup file
	backupData, err := os.Open(backupFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer backupData.Close()

	assert.NoError(t, err, "Backup File error.")

	// create a new one
	newDB, err := os.Create(primaryFile)
	if err != nil {
		fmt.Println(err)
	}

	defer newDB.Close()

	bufferReader := make([]byte, 32)

	for {
		data, err := backupData.Read(bufferReader)
		if err != nil && err != io.EOF {
			fmt.Println(err)
		}

		if data == 0 {
			break
		}

		if _, err := newDB.Write(bufferReader[:data]); err != nil {
			fmt.Println(err)
		}
	}

	assert.FileExists(t, primaryFile, "ToDo.json file not found.")

	assert.NoError(t, err, "File exists.")
}
