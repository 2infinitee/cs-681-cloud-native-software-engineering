package tests

import (
	"fmt"
	"github.com/cs-681-cloud-native-software-engineering/todo-api/voterApi/db"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	BaseApi = "http://localhost:8080"
	cli     = resty.New()
)

// TestMain deletes all data in the database of any data through an api request
func TestMain(m *testing.M) {
	response, err := cli.R().Delete(BaseApi + "/voter")

	if response.StatusCode() != 200 {
		fmt.Printf("Error clearing out the database, %v", err)
	}

	code := m.Run()

	os.Exit(code)
}

// TestApiLoadDatabase loads the database, aka add, through an api request
func TestApiLoadDatabase(t *testing.T) {
	maxPeople := 2
	for i := 0; i < maxPeople; i++ {
		person := createRandomPerson(uint(i))
		response, err := cli.R().SetBody(person).Post(BaseApi + "/voter")
		assert.Nil(t, err)
		assert.Equal(t, 200, response.StatusCode(), "Error, did not receive 200 status code from server.")
	}
}

// TestApiGetAllVoters gets all random people generated from the database through an api request
func TestApiGetAllVoters(t *testing.T) {
	var people []db.Voter

	response, err := cli.R().SetResult(&people).Get(BaseApi + "/voter")
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode(), "Error, did not receive 200 status code from server.")
	assert.Equal(t, 2, len(people))
}

func TestApiGetVoter(t *testing.T) {
	person := createRandomPerson(101)

	_, err := cli.R().SetBody(person).Post(BaseApi + "/voter")
	assert.Nil(t, err)

	addr := fmt.Sprintf("%v%v%v", BaseApi, "/voter/", person.VoterId)

	response, err := cli.R().SetBody(person).Get(addr)
	assert.Nil(t, err, "Error found when grabbing single voter")
	assert.Equal(t, 200, response.StatusCode(), "Error, was not able to grab single voter.")
}

// TestApiDeleteAll deletes all random generated data from the api request
func TestApiDeleteAll(t *testing.T) {
	response, err := cli.R().Delete(BaseApi + "/voter")
	assert.Nil(t, err, "Error found when deleting all voter data.")
	assert.Equal(t, 200, response.StatusCode(), "Error deleting all voters for api handler.")
}

func TestApiDeleteVoter(t *testing.T) {
	person := createRandomPerson(111)

	_, err := cli.R().SetBody(person).Post(BaseApi + "/voter")
	assert.Nil(t, err)

	response, err := cli.R().SetBody(person).Delete(BaseApi + "/voter")
	assert.Nil(t, err, "Error found when deleting single voter")
	assert.Equal(t, 200, response.StatusCode(), "Error, was not able to delete single voter.")
}

func TestApiGetVoterPolls(t *testing.T) {
	var voterStructure []db.VoterHistory
	person := createRandomPerson(11)

	_, err := cli.R().SetBody(person).Post(BaseApi + "/voter")
	assert.Nil(t, err)

	addr := fmt.Sprintf("%v%v%v%v", BaseApi, "/voter/", person.VoterId, "/polls")

	response, err := cli.R().SetBody(&voterStructure).Get(addr)
	assert.Nil(t, err, "Error found when grabbing single voter pools")
	assert.Equal(t, 200, response.StatusCode(), "Error, was not able to grab single voter.")
}

func TestApiGetSingleVoterPoll(t *testing.T) {
	var voterStructure []db.VoterHistory
	person := createRandomPerson(120)

	_, err := cli.R().SetBody(person).Post(BaseApi + "/voter")
	assert.Nil(t, err)

	var pollIds []uint

	for _, data := range person.VoterHistory {
		pollIds = append(pollIds, data.PollId)
	}
	pollNum := pollIds[0]

	addr := fmt.Sprintf("%v%v%v%v%v", BaseApi, "/voter/", person.VoterId, "/polls/", pollNum)

	response, err := cli.R().SetBody(&voterStructure).Get(addr)
	assert.Nil(t, err, "Error found when grabbing single voter pools")
	assert.Equal(t, 200, response.StatusCode(), "Error, was not able to grab single voter, single poll.")
}
