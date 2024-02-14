package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// VoterHistory struct to keep track how many
// times a voter has voted
type VoterHistory struct {
	PollId   uint      `json:"pollId"`
	VoterId  uint      `json:"voterId"`
	VoteDate time.Time `json:"voteDate"`
}

// VoterData struct to keep track of unique
// voter information
type VoterData struct {
	VoterId      uint           `json:"voterId"`
	FirstName    string         `json:"firstName"`
	LastName     string         `json:"lastName"`
	IsDone       bool           `json:"isDone"`
	VoterHistory []VoterHistory `json:"voterHistory,omitempty"`
}

// vMap is an alias for a map of VoterData and the
// key will be VoterData.VoterId
type vMap map[uint]VoterData

// Voter struct to store db data in memory
type Voter struct {
	voterMap vMap
}

// New creates a new map of the database

func New() (*Voter, error) {
	voter := &Voter{
		voterMap: make(map[uint]VoterData),
	}

	return voter, nil
}

// AddVoter allows voter information to be added to the DB
func (v *Voter) AddVoter(voter VoterData) error {
	_, ok := v.voterMap[voter.VoterId]
	if ok {
		return errors.New("voter exists")
	}

	v.voterMap[voter.VoterId] = voter

	return nil
}

// DeleteVoter allows deletion of voter by VoterId
func (v *Voter) DeleteVoter(voterId uint) error {
	delete(v.voterMap, voterId)

	return nil
}

// DeleteAll removes all items from the DB
// to be exposed via /voters
func (v *Voter) DeleteAll() error {
	v.voterMap = make(map[uint]VoterData)
	return nil
}

// UpdateVoter changes voter information
// before it changes it checks to see if voter exists
func (v *Voter) UpdateVoter(voter VoterData) error {
	_, ok := v.voterMap[voter.VoterId]
	if !ok {
		return errors.New("voter does not exist")
	}
	v.voterMap[voter.VoterId] = voter

	return nil
}

// GetVoter gets voter based on id passed
func (v *Voter) GetVoter(voterId uint) (VoterData, error) {

	voter, ok := v.voterMap[voterId]
	if !ok {
		return VoterData{}, errors.New("voter does not exist")
	}

	return voter, nil
}

// ChangeItemDoneStatus is not yet implemented
func (v *Voter) ChangeItemDoneStatus(voterId uint, isDone bool) error {
	voter, ok := v.voterMap[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	voter.IsDone = isDone

	return nil
}

// GetAllVoters grabs all voters in the database
func (v *Voter) GetAllVoters() ([]VoterData, error) {
	var getAllVoterData []VoterData

	for _, voter := range v.voterMap {
		getAllVoterData = append(getAllVoterData, voter)
	}

	return getAllVoterData, nil
}

// PrintVoter outputs voter information to console in pretty format
func (v *Voter) PrintVoter(voter VoterData) {
	jsonBytes, _ := json.MarshalIndent(voter, "", " ")
	fmt.Println(string(jsonBytes))
}

// PrintAllVoters outputs all voter data in pretty format
// the PrintVoter is called per voter data
func (v *Voter) PrintAllVoters(voter []VoterData) {
	for _, voterInfo := range voter {
		v.PrintVoter(voterInfo)
	}
}

// JsonToVoter is a function that allows JSON to be take in as a VoterData
// a string passed in from the CLI.
func (v *Voter) JsonToVoter(jsonString string) (VoterData, error) {
	var voter VoterData
	err := json.Unmarshal([]byte(jsonString), &voter)
	if err != nil {
		return VoterData{}, err
	}

	return voter, nil
}
