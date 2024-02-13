package api

import (
	"github.com/cs-681-cloud-native-software-engineering/todo-api/voterApi/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

func (vApi *VoterAPI) ListAllVoters(content *gin.Context) {
	voterList, err := vApi.db.GetAllVoters()
	if err != nil {
		log.Println("Error Getting All Items: ", err)
		content.AbortWithStatus(http.StatusNotFound)
	}
}
