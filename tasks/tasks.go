package tasks

import (
	"fmt"

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
func SetupTest(serverApi, token string) (string, string, string, error) {
	// testArtifactPath is the full path to the artifact - ex - c:\lab\test-artifact.txt
	// testRepo is the target path to put the artifact in - /repo/folder

	util.ServerApi = serverApi
	util.Token	   = token

	testDirName      := "test-directory"
    testArtifactName := "test-artifact.txt"
    artifactSuffix   := ""
    artifactContents := "Just some test content."
	var kvProps []string
	
	// Prep test artifact
	testDirPath  := common.CreateTestDirectory(testDirName)
	testArtifactPath := common.CreateTestFile(testDirPath, testArtifactName, artifactContents)
	kvProps = append(kvProps,"release=latest-stable")

	// Setup test repo
	testRepoPath, err  := common.CreateTestRepo()   //-->  /test-packer-plugin
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to create repo: " + strErr)
		return "", "", "", err
	}

	// Upload test artifact to test repo
	// Checks for ending slash on target repo path as part of this
	downloadUri, err := operations.UploadFile(testArtifactPath, testRepoPath, artifactSuffix)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to get download URI: " + strErr)
		return "", "", "", err
	}

	artifactUri := common.SetArtifUriFromDownloadUri(downloadUri)

	// Set properties on the test artifact
	statusCode, err := operations.SetArtifactProps(artifactUri, kvProps)
	if statusCode != "204" {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error setting artifact properties: " + strErr)
	}

	return artifactUri, testDirPath, testArtifactPath, nil
}

func TeardownTest(serverApi, token, artifactUri, testDirPath, testArtifactPath string) (string) {
	util.ServerApi = serverApi
	util.Token	   = token
	common.LogTxtHandler().Debug("TEST ARTIFACT URI: " + artifactUri)

	// Delete local test file
	common.DeleteTestFile(testArtifactPath)

	// Delete locat test directory
	common.DeleteTestDirectory(testDirPath)

	// Delete test repo; will delete test artifact with it
	statusCode, err := common.DeleteTestRepo()
	if statusCode == "200" {
		common.LogTxtHandler().Info("Deletion of test repo completed successfully.")
	} else {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to delete test repo - " + strErr)
	}

	if statusCode == "200" {
		return statusCode
	} else {
		return statusCode
	}
}