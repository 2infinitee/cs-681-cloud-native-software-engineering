package tests

import (
	"fmt"
	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/cs-681-cloud-native-software-engineering/todo-api/voterApi/db"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	database *db.Voter
)

func init() {
	testDatabase, err := db.New()
	if err != nil {
		fmt.Println("Error Creating DB: ", err)
		os.Exit(1)
	}

	database = testDatabase
}

func createRandomPerson(voterId uint) db.VoterData {

	// create a random data
	person := db.VoterData{
		VoterId:   voterId,
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		IsDone:    false,
		VoterHistory: []db.VoterHistory{
			{
				PollId:   uint(fake.Number(100, 999)),
				VoterId:  voterId,
				VoteDate: fake.Date(),
			},
			{
				PollId:   uint(fake.Number(100, 999)),
				VoterId:  voterId,
				VoteDate: fake.Date(),
			},
		},
	}
	return person
}

// TestDbAddAndGetVoter creates a random person and adds them to the database
func TestDbAddAndGetVoter(t *testing.T) {

	var person db.VoterData
	size := 3
	for i := 0; i < size; i++ {
		person = createRandomPerson(uint(i))
		err := database.AddVoter(person)
		assert.Nil(t, err, "Error, was not able to add random data to database.")
	}

	data, err := database.GetVoter(2)
	assert.NoError(t, err, "Error, on GET when trying to GetVoter.")
	assert.Equal(t, person.VoterId, data.VoterId, "Error did not find matching database size.")
}

// TestDbDeleteVoter creates a random person, adds them to a database and deletes a voter from the database
func TestDbDeleteVoter(t *testing.T) {
	person := createRandomPerson(uint(1))

	err := database.AddVoter(person)
	assert.NoError(t, err, "Error, was not able to add random data to database.")

	err = database.DeleteVoter(person.VoterId)
	assert.NoError(t, err, "Was not able to delete voter from database.")
	assert.NotEqual(t, person, nil, "Found random person generated in database.")
}

// TestDbDeleteAll deletes random generated people added to the database
func TestDbDeleteAll(t *testing.T) {

	var people []db.VoterData
	size := 4
	for i := 0; i < size; i++ {
		person := createRandomPerson(uint(i))
		err := database.AddVoter(person)
		assert.Nil(t, err, "Error, was not able to add random data to database.")
	}

	err := database.DeleteAll()
	assert.NoError(t, err, "Error, was not able to delete all random data in database.")
	assert.NotEqual(t, 4, len(people), "Found random person generated in database.")

}

// TestDbGetAllVoter GET all random generated people added to the database
func TestDbGetAllVoter(t *testing.T) {
	personOne := createRandomPerson(uint(1))
	personTwo := createRandomPerson(uint(2))
	personThree := createRandomPerson(uint(3))
	personFour := createRandomPerson(uint(4))

	err := database.AddVoter(personOne)
	assert.NoError(t, err, "Error, was not able to add random data to database.")
	err = database.AddVoter(personTwo)
	assert.NoError(t, err, "Error, was not able to add random data to database.")
	err = database.AddVoter(personThree)
	assert.NoError(t, err, "Error, was not able to add random data to database.")
	err = database.AddVoter(personFour)
	assert.NoError(t, err, "Error, was not able to add random data to database.")

	_, err = database.GetAllVoters()
	assert.NoError(t, err, "Error, was not able to delete all random data in database.")
}

// TestDbUpdateVoter updates a voter in the database in this case it updates the first name only
func TestDbUpdateVoter(t *testing.T) {
	person := createRandomPerson(uint(1))

	err := database.AddVoter(person)
	assert.NoError(t, err, "Error, was not able to add random data to database.")

	originalFirstName := person.FirstName
	person.FirstName = fake.FirstName()

	err = database.UpdateVoter(person)
	assert.NoError(t, err, "Error, was not able to update random data to database.")

	vData, err := database.GetVoter(person.VoterId)
	assert.NotEqual(t, originalFirstName, vData.FirstName, "Error, first name did not update.")
}

// TestDbGetAllVoterPolls grabs the slice of voter history and
// makes sure it's been added in the database correctly
func TestDbGetAllVoterPolls(t *testing.T) {
	person := createRandomPerson(uint(1))

	err := database.AddVoter(person)
	assert.NoError(t, err, "Error, was not able to add random data to database.")

	polls, err := database.GetAllVoterPolls(person.VoterId)
	assert.NoError(t, err, "Error, was not able to GET voter polls.")
	assert.Equal(t, person.VoterHistory, polls, "Error, did not find matching voter history.")
}

// TestDbGetVoterPoll grabs the poll id in a voters history and
// makes sure it's been added in the database correctly
func TestDbGetVoterPoll(t *testing.T) {
	person := createRandomPerson(uint(1))

	err := database.AddVoter(person)
	assert.NoError(t, err, "Error, was not able to add random data to database.")

	var pollIds []uint

	for _, data := range person.VoterHistory {
		pollIds = append(pollIds, data.PollId)
	}
	poll := pollIds[0]

	polls, err := database.GetVoterPoll(person.VoterId, poll)
	assert.NoError(t, err, "Error, was not able to GET voter polls.")
	assert.Equal(t, poll, polls.PollId, "Error, did not find matching voter poll id.")
}

// TestDbChangeDoneStatus updates isDone field in VoterData struct and sees
// if changes has been updated in the database
func TestDbChangeDoneStatus(t *testing.T) {
	person := createRandomPerson(uint(1))

	err := database.AddVoter(person)
	assert.NoError(t, err, "Error, was not able to add random data to database.")

	person.IsDone = true

	err = database.ChangeDoneStatus(person.VoterId, person.IsDone)
	assert.NoError(t, err, "Error, was not able to update random data to database.")

	vData, err := database.GetVoter(person.VoterId)
	assert.Equal(t, true, vData.IsDone, "Error, did not find same isDone value.")
}

// TestDbPrintVoter creates a person and returns the information in a pretty JSON format
func TestDbPrintVoter(t *testing.T) {
	person := createRandomPerson(uint(1))

	err := database.PrintVoter(person)
	assert.NoError(t, err, "Error, when converting JSON string to pretty JSON format.")
}

// TestDbPrintAllVoters creates two random person and prints JSON string into pretty format
func TestDbPrintAllVoters(t *testing.T) {
	personOne := createRandomPerson(uint(1))
	personTwo := createRandomPerson(uint(2))

	err := database.DeleteAll()
	assert.NoError(t, err, "Error, was not able to delete all random data in database.")

	err = database.AddVoter(personOne)
	assert.NoError(t, err, "Error, was not able to add random data to database.")
	err = database.AddVoter(personTwo)
	assert.NoError(t, err, "Error, was not able to add random data to database.")

	voters, err := database.GetAllVoters()
	assert.NoError(t, err, "Error, was not able to delete all random data in database.")

	err = database.PrintAllVoters(voters)
	assert.NoError(t, err, "Error, when converting JSON string to pretty JSON format.")

	err = database.DeleteAll()
	assert.NoError(t, err, "Error, was not able to delete all random data in database.")
}
