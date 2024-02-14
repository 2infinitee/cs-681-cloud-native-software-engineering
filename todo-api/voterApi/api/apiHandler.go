package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/cs-681-cloud-native-software-engineering/todo-api/voterApi/db"
	"github.com/gin-gonic/gin"
)

// VoterAPI creates and maintains a reference to the data handler
type VoterAPI struct {
	db *db.Voter
}

// New allows the start of a new api handler
func New() (*VoterAPI, error) {
	dbHandler, err := db.New()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{db: dbHandler}, nil
}

// ListAllVoters implements a GET /voter to grab all voters and their data
func (api *VoterAPI) ListAllVoters(ctx *gin.Context) {
	voterList, err := api.db.GetAllVoters()
	if err != nil {
		log.Println("Error getting all voters: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if voterList == nil {
		voterList = make([]db.VoterData, 0)
	}

	ctx.JSON(http.StatusOK, voterList)
}

// ListSelectVoters implements GET /v2/voter
// and returns voters that are either done or not done
// depending on the value set /v2/voter?done=true
func (api *VoterAPI) ListSelectVoters(ctx *gin.Context) {

	// load data into memory
	voterList, err := api.db.GetAllVoters()
	if err != nil {
		log.Println("Error getting all voters: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	// if the db is empty make an empty slice
	// so that JSON marshaling works
	if voterList == nil {
		voterList = make([]db.VoterData, 0)
	}

	doneStatus := ctx.Query("done")

	// if doneStatus is empty, then return all voters
	if doneStatus == "" {
		ctx.JSON(http.StatusOK, voterList)
		return
	}

	// if doneStatus is not empty then we need to filter
	// based on the doneStatus list
	done, err := strconv.ParseBool(doneStatus)
	if err != nil {
		log.Println("Error converting done to bool: ", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// filter list based on done value
	var filteredList []db.VoterData
	for _, voter := range voterList {
		if voter.IsDone == done {
			filteredList = append(filteredList, voter)
		}
	}

	// if db returns a nil there are no voters in db
	// need to convert this to an empty slice to return
	if filteredList == nil {
		filteredList = make([]db.VoterData, 0)
	}

	ctx.JSON(http.StatusOK, filteredList)
}

// GetVoter implements GET method /voter/:voterId
// returns a single voter
func (api *VoterAPI) GetVoter(ctx *gin.Context) {

	voterId := ctx.Param("voterId")
	convertIdToInt64, err := strconv.ParseInt(voterId, 10, 32)
	if err != nil {
		log.Println("Error Converting voterId to int64", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := api.db.GetVoter(int(convertIdToInt64))
	if err != nil {
		log.Println("Voter not found: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, voter)
}

// AddVoter implements POST method /voter
// adds a new voter to the db through the api
func (api *VoterAPI) AddVoter(ctx *gin.Context) {
	var voterData db.VoterData

	if err := ctx.ShouldBindJSON(&voterData); err != nil {
		log.Println("Error binding JSON: ", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := api.db.AddVoter(voterData); err != nil {
		log.Println("Error adding voter", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, voterData)
}

// UpdateVoter implements PUT method /voter
// adds a new voter to the db through the api
func (api *VoterAPI) UpdateVoter(ctx *gin.Context) {
	var voterData db.VoterData
	if err := ctx.ShouldBindJSON(&voterData); err != nil {
		log.Println("Error binding JSON: ", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := api.db.UpdateVoter(voterData); err != nil {
		log.Println("Error adding voter", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, voterData)
}

// DeleteVoter implements DELETE /voter/:voterId
// deletes a single voter
func (api *VoterAPI) DeleteVoter(ctx *gin.Context) {
	voterId := ctx.Param("voterId")
	convertIdToInt64, _ := strconv.ParseInt(voterId, 10, 32)

	if err := api.db.DeleteVoter(int(convertIdToInt64)); err != nil {
		log.Println("Error deleting voter: ", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)

}

// DeleteAllVoters implements DELETE /voter
// deletes all voter
func (api *VoterAPI) DeleteAllVoters(ctx *gin.Context) {

	if err := api.db.DeleteAll(); err != nil {
		log.Println("Error deleting all items: ", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func (api *VoterAPI) CrashSimulator(ctx *gin.Context) {
	panic("Simulating a unexpected crash")
}

func (api *VoterAPI) HealthCheck(ctx *gin.Context) {

	ctx.JSON(http.StatusOK,
		gin.H{
			"status:":            "ok",
			"version":            "0.0.1",
			"uptime":             1,
			"voters_processed":   1,
			"errors_encountered": 1,
		})
}
