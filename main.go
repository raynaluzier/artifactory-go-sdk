package main

import (
	"os"

	_ "github.com/raynaluzier/artifactory-go-sdk/common"
	_ "github.com/raynaluzier/artifactory-go-sdk/operations"
	_ "github.com/raynaluzier/artifactory-go-sdk/search"
	"github.com/raynaluzier/artifactory-go-sdk/util"
)

func main(){

	// -------------- TESTING --------------------------------------------------
	// mimicking passing in vars and then assigning them to the global vars
	serverApi 	:= os.Getenv("ARTIFACTORY_SERVER")
	token 		:= os.Getenv("ARTIFACTORY_TOKEN")
	logLevel 	:= os.Getenv("LOGGING")
	//outputDir 	:= //os.Getenv("OUTPUTDIR")   //c:\lab\output-testing\ or /lab/output-testing

	util.ServerApi = serverApi
	util.Token     = token
	util.Logging   = logLevel
	//util.OutputDir = outputDir
	// -------------------------------------------------------------------------
	



}
