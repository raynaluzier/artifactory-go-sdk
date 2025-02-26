package archive

import (
	"fmt"

	"github.com/raynaluzier/artifactory-go-sdk/common"
	"github.com/raynaluzier/artifactory-go-sdk/operations"
	"github.com/raynaluzier/artifactory-go-sdk/util"
)

func UploadArtifact(serverApi, token, sourcePath, targetPath string) (string, string, error) {
	// Single file
	util.ServerApi = serverApi
	util.Token	   = token
	var artifactUri string
	common.LogTxtHandler().Debug("UPLOADING NEW ARTIFACT TO ARTIFACTORY...")

	downloadUri, err := operations.UploadFile(sourcePath, targetPath)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to upload artifact - " + strErr)
		return "", "", err
	} else {
		artifactUri = common.SetArtifUriFromDownloadUri(downloadUri)
	}

	return downloadUri, artifactUri, nil
}