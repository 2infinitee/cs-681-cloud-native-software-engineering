package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cs-681-cloud-native-software-engineering/todo-api/voterApi/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// global variables from the cli as flags
var (
	hostFlag string
	portFlag uint
)

// initializeClientFlags parses flags provided from the cli
func initializeFlags() {

	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listens on all interfaces")
	flag.UintVar(&portFlag, "p", 8080, "Default port is set to 8080")

	flag.Parse()
}

func main() {
	initializeFlags()

	// gin-contrib/cors is the gin-cors is a middleware
	// it implements Cross Origin Resource Sharing specs
	// from WC3 which enables pages within a browser to
	// consume resources such as REST APIs
	instance := gin.Default()
	instance.Use(cors.Default())

	apiHandler, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	instance.GET("/voter", apiHandler.ListAllVoters)
	instance.GET("/voter/:voterId", apiHandler.GetVoter)
	instance.POST("/voter", apiHandler.AddVoter)
	instance.PUT("/voter", apiHandler.UpdateVoter)
	instance.DELETE("/voter/:voterId", apiHandler.DeleteVoter)
	instance.DELETE("/voter", apiHandler.DeleteAllVoters)

	instance.GET("/crash", apiHandler.CrashSimulator)
	instance.GET("/health", apiHandler.HealthCheck)

	v2 := instance.Group("/v2")
	v2.GET("/voter", apiHandler.ListSelectVoters)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	instance.Run(serverPath)
}
