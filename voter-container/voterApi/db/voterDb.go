package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

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
	VoterHistory []VoterHistory `json:"voterHistory"`
}

// TODO REMOVE ALL BELOW HERE
//// vMap is an alias for a map of VoterData and the
//// key will be VoterData.VoterId
//type vMap map[uint]VoterData
//type pMap map[uint]VoterHistory

// Voter struct to store db data in memory
type Voter struct {
	cache
}

// VoterPolls struct to store db data in memory
type VoterPolls struct {
	cache
}

// New creates a new map of the database

func New() (*Voter, error) {

	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}

	return NewWithCacheInstance(redisUrl)
}

// NewWithCacheInstance is a constructor that returns a pointer
// to Voter struct
func NewWithCacheInstance(location string) (*Voter, error) {

	// the db is the client connecting to redis
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	// ctx is used to coordinate with redis
	ctx := context.Background()

	// ensure redis connection is working
	// highly recommended way
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error() + " 1-basic not available, continuing...")
		return nil, err
	}

	// allows json to be stored in redis
	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//Return a pointer to a new Voter struct
	return &Voter{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

// func to be used later
func isRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

// redis keys will be strings
func redisKeyFromId(id int) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

// a help to return VoterData from redis from a provided key
func (v *Voter) getVoterFromRedis(key string, voter *VoterData) error {

	// query an voter object
	voterObject, err := v.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(voterObject.([]byte), voter)
	if err != nil {
		return err
	}

	return nil
}

// AddVoter allows voter information to be added to the DB
func (v *Voter) AddVoter(voter VoterData) error {

	redisKey := redisKeyFromId(int(voter.VoterId))

	var existingVoter VoterData
	if err := v.getVoterFromRedis(redisKey, &existingVoter); err == nil {
		return errors.New("items already exists")
	}

	// add item to redis with JSON set
	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return err
	}

	return nil
}

// DeleteVoter allows deletion of voter by VoterId
func (v *Voter) DeleteVoter(voterId uint) error {

	pattern := redisKeyFromId(int(voterId))

	numDeleted, err := v.cacheClient.Del(v.context, pattern).Result()
	if err != nil {
		return nil
	}
	if numDeleted == 0 {
		return errors.New("attempted to delete non-existent item")
	}

	return nil
}

// DeleteAll removes all items from the DB
// to be exposed via /voters
func (v *Voter) DeleteAll() error {

	pattern := RedisKeyPrefix + "*"
	ks, _ := v.cacheClient.Keys(v.context, pattern).Result()

	numDeleted, err := v.cacheClient.Del(v.context, ks...).Result()

	if err != nil {
		return err
	}

	if numDeleted != int64(len(ks)) {
		return errors.New("one or more items could not be deleted")
	}

	return nil
}

// UpdateVoter changes voter information
// before it changes it checks to see if voter exists
func (v *Voter) UpdateVoter(voter VoterData) error {

	redisKey := redisKeyFromId(int(voter.VoterId))

	var existingItem VoterData

	if err := v.getVoterFromRedis(redisKey, &existingItem); err != nil {
		return errors.New("items does not exist")
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return nil
	}

	return nil
}

// GetVoter gets voter based on id passed
func (v *Voter) GetVoter(voterId uint) (VoterData, error) {

	var voter VoterData
	pattern := redisKeyFromId(int(voterId))
	err := v.getVoterFromRedis(pattern, &voter)

	if err != nil {
		return VoterData{}, err
	}

	return voter, nil
}

// GetAllVoterPolls gets voter based on id passed
func (v *Voter) GetAllVoterPolls(voterId uint) ([]VoterHistory, error) {

	var voter VoterData

	pattern := redisKeyFromId(int(voterId))
	err := v.getVoterFromRedis(pattern, &voter)
	if err != nil {
		return nil, err
	}

	return voter.VoterHistory, nil
}

// GetVoterPoll gets voter based on id passed
func (v *Voter) GetVoterPoll(voterId uint, pollId uint) (VoterHistory, error) {

	var voter VoterData

	pattern := redisKeyFromId(int(voterId))
	err := v.getVoterFromRedis(pattern, &voter)

	if err != nil {
		return VoterHistory{}, err
	}

	voterHistoryMap := make(map[uint]VoterHistory)

	for _, data := range voter.VoterHistory {
		voterHistoryMap[data.PollId] = data
	}

	return voterHistoryMap[pollId], nil
}

// ChangeDoneStatus is not yet implemented
func (v *Voter) ChangeDoneStatus(voterId uint, isDone bool) error {

	var voter VoterData

	redisKey := redisKeyFromId(int(voterId))
	err := v.getVoterFromRedis(redisKey, &voter)
	if err != nil {
		return errors.New("isDone status error")
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", isDone); err != nil {
		return err
	}

	return nil
}

// GetAllVoters grabs all voters in the database
func (v *Voter) GetAllVoters() ([]VoterData, error) {

	var voterList []VoterData
	var voterData VoterData

	pattern := RedisKeyPrefix + "*"
	ks, _ := v.cacheClient.Keys(v.context, pattern).Result()

	for _, key := range ks {
		err := v.getVoterFromRedis(key, &voterData)
		if err != nil {
			return nil, err
		}
		voterList = append(voterList, voterData)
	}

	return voterList, nil
}

// PrintVoter outputs voter information to console in pretty format
func (v *Voter) PrintVoter(voter VoterData) error {

	jsonBytes, err := json.MarshalIndent(voter, "", " ")
	if err != nil {
		return errors.New("could not convert data to pretty JSON format")
	}
	fmt.Println(string(jsonBytes))

	return nil
}

// PrintAllVoters outputs all voter data in pretty format
// the PrintVoter is called per voter data
func (v *Voter) PrintAllVoters(voter []VoterData) error {
	for _, voterInfo := range voter {
		err := v.PrintVoter(voterInfo)
		if err != nil {
			return err
		}
	}
	return nil
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
