package tasks

import (
	"fmt"
	"log"

	"github.com/raynaluzier/artifactory-go-sdk/common"
	"github.com/raynaluzier/artifactory-go-sdk/operations"
	"github.com/raynaluzier/artifactory-go-sdk/util"
)

// Must have Artifactory instance licensed at Pro or higher, access to create/remove repos and artifacts
func SetupTest(serverApi, token, testRepoName, configFilePath, testArtifact, fileSuffix string) (string, error) {
	util.ServerApi = serverApi
	util.Token	   = token

	testRepo, err  := common.CreateTestRepo(testRepoName, configFilePath)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	downloadUri, err := operations.UploadFile(testArtifact, testRepo, fileSuffix)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	artifactUri := common.SetArtifUriFromDownloadUri(downloadUri)

	return artifactUri, nil
}

func TeardownTest (serverApi, token, testRepoName, artifactUri string) (string) {
	util.ServerApi = serverApi
	util.Token	   = token

	// Delete test artifact
	statusCodeArtif, err := operations.DeleteArtifact(artifactUri)
	if statusCodeArtif == "204" {
		fmt.Println("Deletion of test artifact completed successfully.")
	} else {
		fmt.Println("Unable to delete test artifact - ", err)
	}

	// Delete test repo
	statusCodeRepo, err := common.DeleteTestRepo(testRepoName)
	if statusCodeRepo == "204" {
		fmt.Println("Deletion of test repo completed successfully.")
	} else {
		fmt.Println("Unable to delete test repo - ", err)
	}

	if statusCodeArtif == "204" && statusCodeRepo == "204" {
		statusCode := "204"
		return statusCode
	} else {
		statusCode := "404"
		return statusCode
	}
}