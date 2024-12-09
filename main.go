package main

import (
	"fmt"

	"os"

	_ "github.com/raynaluzier/go-artifactory/common"
	"github.com/raynaluzier/go-artifactory/operations"
	_ "github.com/raynaluzier/go-artifactory/search"
	"github.com/raynaluzier/go-artifactory/tasks"
	"github.com/raynaluzier/go-artifactory/util"
)

func main(){

	// -------------- TESTING --------------------------------------------------
	// mimicking passing in vars and then assigning them to the global vars
	serverApi 	:= os.Getenv("ARTIFACTORY_SERVER")
	token 		:= os.Getenv("ARTIFACTORY_TOKEN")
	logLevel 	:= os.Getenv("ARTIFACTORY_LOGGING")
	//outputDir 	:= os.Getenv("ARTIFACTORY_OUTPUTDIR")

	util.ServerApi = serverApi
	util.Token     = token
	util.Logging   = logLevel
	//util.OutputDir = outputDir
	// -------------------------------------------------------------------------
	
	kvProps := []string{}
	kvProps = append(kvProps, "release=latest-stable")
	//kvProps = append(kvProps, "release=stable")
	
	ext := "vmxt"
	artifName := "win-22"

	fmt.Println(tasks.GetImageDetails(util.ServerApi, util.Token, util.Logging, artifName, ext, kvProps))
	fmt.Println(operations.ListRepos())
	/*
	//artifName := "W22"
	listArtifacts, err := search.GetArtifactsByName(serverApi, token, artifName)
	if err != nil {
		fmt.Println("some error")
	}

	newArtifacts, err := (search.FilterListByFileType(ext, listArtifacts))
	if err != nil {
		fmt.Println("some other error")
	}

	result, err := (operations.FilterListByProps(token, newArtifacts, kvProps))
	if err != nil {
		fmt.Println("some other error")
	}
	fmt.Println(operations.GetArtifactNameFromUri(result))
	*/

}
