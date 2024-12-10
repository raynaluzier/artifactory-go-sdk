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