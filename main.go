package main

import (
	"fmt"
	"os"

	"github.com/raynaluzier/artifactory-go-sdk/common"
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
	logLevel 	:= os.Getenv("LOGGING")
	//outputDir 	:= //os.Getenv("OUTPUTDIR")   //c:\lab\output-testing\ or /lab/output-testing

	util.ServerApi = serverApi
	util.Token     = token
	util.Logging   = logLevel
	//util.OutputDir = outputDir
	// -------------------------------------------------------------------------

	const testDirName      = "test-directory"
	const testArtifactName = "test-artifact.txt"   // this gets renamed to an ova file later
	const artifactContents = "Just some test content."
	var kvProps []string
	const setAsOva = true
	testDirPath        := common.CreateTestDirectory(testDirName)
	// Returns $HOME_DIR/test-directory/test-artifact.ova
	testArtifactPath   := common.CreateTestFile(testDirPath, testArtifactName, artifactContents, setAsOva)  // the artifact will be renamed to test-artifact.ova
	uploadTestArtifact := false  // Don't need to upload artifact as part of test setup; test itself will do this.
	kvProps = append(kvProps,"release=latest-stable")

	status, err := tasks.SetupTest(serverApi, token, testArtifactPath, kvProps, uploadTestArtifact)		
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Status of setup: " + status)

	imageType := "ova"
	imageName := "test-artifact"
	sourcePath := common.CheckAddSlashToPath(testDirPath)
	targetPath := "/test-packer-plugin"
	result := tasks.UploadArtifacts(serverApi, token, imageType, imageName, sourcePath, targetPath)
	fmt.Println(result)

	artifactUri, artifactName, createDate, downloadUri, err := tasks.GetImageDetails(serverApi, token, "test-artifact", "ova", kvProps)
	fmt.Println("Download URI: " + downloadUri)
	fmt.Println("Artifact URI: " + artifactUri)
	fmt.Println("Created Date: " + createDate)
	fmt.Println("Artifact Name: " + artifactName)
	if err != nil {
		fmt.Println(err)
	}
}
