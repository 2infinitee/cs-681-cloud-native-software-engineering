package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"

	"github.com/gin-contrib/cors"
)

// global variables from the cli as flags
var (
	hostFlag string
	portFlag uint
)

// initializeClientFlags parses flags provided from the cli
func initializeClientFlags() {

	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listens on all interfaces")
	flag.UintVar(&portFlag, "p", 8080, "Default port is set to 8080")

	flag.Parse()
}

func main() {
	initializeClientFlags()

	// gin-contrib/cors is the gin-cors is a middleware
	// it implements Cross Origin Resource Sharing specs
	// from WC3 which enables pages within a browser to
	// consume resources such as REST APIs
	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/voter", apiHandler.)



}