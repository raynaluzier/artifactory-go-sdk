package tasks

import (
	"fmt"
	"log"

	"github.com/raynaluzier/artifactory-go-sdk/common"
	"github.com/raynaluzier/artifactory-go-sdk/operations"
	"github.com/raynaluzier/artifactory-go-sdk/search"
	"github.com/raynaluzier/artifactory-go-sdk/util"
)

func GetImageDetails(serverApi, token, logLevel, artifName, ext string, kvProps []string) (string, string, string, string) {
	util.ServerApi = serverApi
	util.Token     = token
	util.Logging   = logLevel
	var artifactUri string
	var strErr string

	listArtifacts, err := search.GetArtifactsByName(artifName)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error getting list of matching artifacts - " + strErr)
	}

	listByFileType, err := search.FilterListByFileType(ext, listArtifacts)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error filtering artifacts by file type - " + strErr)
	}

	if len(listByFileType) == 1 {
		// if just one artifact, we'll return it
		artifactUri = listByFileType[0]
	} else if len(listByFileType) > 1 && len(kvProps) != 0 {
		artifactUri, err = operations.FilterListByProps(listByFileType, kvProps)
		if err != nil {
			strErr = fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error filtering artifacts by file type - " + strErr)
		}
	} else {
		// if no props passed, but more than one artif is in list, return latest
		artifactUri, err = operations.GetLatestArtifactFromList(listByFileType)
		if err != nil {
			strErr = fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error getting latest artifact from list - " + strErr)
		}
	}

	artifactName := operations.GetArtifactNameFromUri(artifactUri)

	createDate, err := operations.GetCreateDate(artifactUri)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to get create date of artifact - " + strErr)
	}

	downloadUri, err := operations.GetDownloadUri(artifactUri)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to get download URI - " + strErr)
	}
	
	return artifactUri, artifactName, createDate, downloadUri
}

// Must have Artifactory instance licensed at Pro or higher, access to create/remove repos and artifacts
func SetupTest(serverApi, token, testRepoName, configFilePath, testArtifact, fileSuffix string, kvProps []string) (string, error) {
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

	statusCode, err := operations.SetArtifactProps(artifactUri, kvProps)
	if statusCode != "204" {
		log.Fatal(err)
	}

	return artifactUri, nil
}

func TeardownTest(serverApi, token, testRepoName, artifactUri string) (string) {
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