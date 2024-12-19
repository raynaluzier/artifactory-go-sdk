package main

import (
	"fmt"
	"os"

	_ "github.com/raynaluzier/artifactory-go-sdk/common"
	_ "github.com/raynaluzier/artifactory-go-sdk/operations"
	_ "github.com/raynaluzier/artifactory-go-sdk/search"
	"github.com/raynaluzier/artifactory-go-sdk/tasks"
	"github.com/raynaluzier/artifactory-go-sdk/util"
)

func main(){

	// -------------- TESTING --------------------------------------------------
	// mimicking passing in vars and then assigning them to the global vars
	serverApi 	:= os.Getenv("ARTIFACTORY_SERVER")
	token 		:= os.Getenv("ARTIFACTORY_TOKEN")
	logLevel 	:= os.Getenv("ARTIFACTORY_LOGGING")
	//outputDir 	:= os.Getenv("ARTIFACTORY_OUTPUTDIR")   //c:\lab\output-testing\ or /lab/output-testing

	util.ServerApi = serverApi
	util.Token     = token
	util.Logging   = logLevel
	//util.OutputDir = outputDir
	// -------------------------------------------------------------------------
	
	/*
    testArtifact := "C:\\lab\\test-artifact.txt"
    fileSuffix := ""
    kvProps := []string{}
    kvProps = append(kvProps, "release=latest-stable")
    fmt.Println(tasks.SetupTest(util.ServerApi, util.Token, testArtifact, fileSuffix, kvProps))
    */

    
    //artifactUri := "https://riverpointtechnology.jfrog.io/artifactory/api/storage/rpt-libs-local/ecp/win/win-22-4444444.vmxt"
    fmt.Println(tasks.TeardownTest(util.ServerApi, util.Token))
    
}
