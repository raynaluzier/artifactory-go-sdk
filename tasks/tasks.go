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
func SetupTest(serverApi, token, testArtifactPath, artifactSuffix string, kvProps []string, uploadArtifact bool) (string, error) {
	// testArtifactPath is the full path to the artifact -> ex - c:\lab\test-artifact.txt
	// testRepoPath is the target path to put the artifact in -> /test-packer-plugin

	util.ServerApi = serverApi
	util.Token	   = token

	// Setup test repo
	testRepoPath, err  := common.CreateTestRepo()   //-->  /test-packer-plugin
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to create repo: " + strErr)
		return "Incomplete", err
	}

	if uploadArtifact == true {
		// Upload test artifact to test repo
		// Checks for ending slash on target repo path as part of this
		downloadUri, err := operations.UploadFile(testArtifactPath, testRepoPath, artifactSuffix)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Unable to get download URI: " + strErr)
			return "", err
		}

		artifactUri := common.SetArtifUriFromDownloadUri(downloadUri)

		// Set properties on the test artifact
		statusCode, err := operations.SetArtifactProps(artifactUri, kvProps)
		if statusCode != "204" {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error setting artifact properties: " + strErr)
		}
		return artifactUri, nil
	} else {
		return "Complete", nil
	}
	
}

func TeardownTest(serverApi, token string) (string) {
	util.ServerApi = serverApi
	util.Token	   = token
	common.LogTxtHandler().Debug("DELETING TEST REPO AND ARTIFACT...")

	// Deletes test repo; also deletes test artifact with it
	statusCode, err := common.DeleteTestRepo()
	if statusCode == "200" {
		common.LogTxtHandler().Info("Deletion of test repo with test artifact completed successfully.")
	} else {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to delete test repo and artifact - " + strErr)
	}

	if statusCode == "200" {
		return statusCode
	} else {
		return statusCode
	}
}

func UploadArtifact(serverApi, token, sourcePath, targetPath, fileSuffix string) (string, string, error) {
	util.ServerApi = serverApi
	util.Token	   = token
	var artifactUri string
	common.LogTxtHandler().Debug("UPLOADING NEW ARTIFACT TO ARTIFACTORY...")

	downloadUri, err := operations.UploadFile(sourcePath, targetPath, fileSuffix)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to upload artifact - " + strErr)
		return "", "", err
	} else {
		artifactUri = common.SetArtifUriFromDownloadUri(downloadUri)
	}

	return downloadUri, artifactUri, nil
}

func SetProps(serverApi, token, artifUri string, kvProps []string) (string, error) {
	util.ServerApi = serverApi
	util.Token	   = token

	common.LogTxtHandler().Debug("UPDATING PROPERTIES OF ARTIFACT...")

	statusCode, err := operations.SetArtifactProps(artifUri, kvProps)
	if statusCode == "204" {
		props, err := operations.GetAllPropsForArtifact(artifUri)

		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Unable to get artifact properties - " + strErr)
			return "", err
		}
		fmt.Println(props)
		return statusCode, nil

	} else {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to set artifact properties - " + strErr)
		return "", err
	}
}