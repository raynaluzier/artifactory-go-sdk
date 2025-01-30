package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	// Single file
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

func UploadArtifacts(serverApi, token, imageType, imageName, sourceDir, targetDir, fileSuffix string) (string) {
	// sourceDir ex: c:\\lab or /lab - We'll check for/add ending slash if needed
	// targetDir ex: /repo-name/folder - We'll check for/add ending slash if needed
	util.ServerApi = serverApi
	util.Token	   = token
	var fileName string
	var err error
	var fileTypes []string
	imageType = strings.ToLower(imageType)

	common.LogTxtHandler().Debug("UPLOADING NEW ARTIFACT TO ARTIFACTORY...")
	newSourceDir := common.CheckAddSlashToPath(sourceDir)  // makes sure ending slash exists
	newTargetDir := common.CheckAddSlashToPath(targetDir)
	items, _ := os.ReadDir(sourceDir)

	if imageType == "ova" {
		if fileSuffix != "" {
			fileName = imageName + "-" + fileSuffix + ".ova"
		} else {
			fileName = imageName + ".ova"
		}

		result, err := operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error(result + " - " + strErr)
		}


	} else if imageType == "ovf" {
	    fileTypes = []string{".ovf", ".mf"}
		for _, ft := range fileTypes {
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + ft
			} else {
				fileName = imageName + ft
			}

			result, err := operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error(result + " - " + strErr)
			}
		}

		// Search and upload related OVF-based disk files
	 	for i := 1; i < 15; i++ {
			strI := strconv.Itoa(i)
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + "-disk-" + strI + ".vmdk"
			} else {
				fileName = imageName + "-disk-" + strI + ".vmdk"
			}

			result, err := operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error(result + " - " + strErr)
			}
		}
	} else if imageType == "vmtx" {
		var result string
		fileTypes = []string{".nvram", ".vmsd", ".vmtx", ".vmxf"}
		for _, ft := range fileTypes {
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + ft
			} else {
				fileName = imageName + ft
			}

			result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error(result + " - " + strErr)
			}
		}

		// Search and upload non-numbered virtual disk file
		if fileSuffix != "" {
			fileName = imageName + "-" + fileSuffix + ".vmdk"
		} else {
			fileName = imageName + ".vmdk"
		}
		result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error(result + " - " + strErr)
		}

		// Search and upload numbered disk files
		for i := 1; i < 15; i++ {
			strI := strconv.Itoa(i)
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + "_" + strI + ".vmdk"
			} else {
				fileName = imageName + "_" + strI + ".vmdk"
			}

			result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error(result + " - " + strErr)
			}
		}

		// Search and upload -ctk disk files ----------------------------------->
		// Search and upload non-numbered virtual disk file
		if fileSuffix != "" {
			fileName = imageName + "-" + fileSuffix + "-ctk.vmdk"
		} else {
			fileName = imageName + "-ctk.vmdk"
		}
		result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error(result + " - " + strErr)
		}

		// Search and upload numbered -ctk disk files
		for i := 1; i < 15; i++ {
			strI := strconv.Itoa(i)
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + "_" + strI + "-ctk.vmdk"
			} else {
				fileName = imageName + "_" + strI + "-ctk.vmdk"
			}

			result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error(result + " - " + strErr)
			}
		}
		
		// Search and upload -flat disk files ----------------------------------->
		// Search and upload non-numbered virtual disk file
		if fileSuffix != "" {
			fileName = imageName + "-" + fileSuffix + "-flat.vmdk"
		} else {
			fileName = imageName + "-flat.vmdk"
		}
		result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error(result + " - " + strErr)
		}

		// Search and upload numbered -flat disk files
		for i := 1; i < 15; i++ {
			strI := strconv.Itoa(i)
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + "_" + strI + "-flat.vmdk"
			} else {
				fileName = imageName + "_" + strI + "-flat.vmdk"
			}

			result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error(result + " - " + strErr)
			}
		}
	} else {
		common.LogTxtHandler().Error("Unsupported or blank image type. Supported image types are OVA, OVF, and VMTX.")
		return "Unsupported image type"
	}
	return "End of upload process"
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

func DownloadArtifacts(serverApi, token, downloadUri, outputDir string) string {
	// Takes in download URI that corresponds to OVA, OVF, or VMTX file; 
	// Will then determine other expected associated artifacts and download those as well
	util.ServerApi = serverApi
	util.Token = token
	
	common.LogTxtHandler().Debug("DOWNLOADING ARTIFACT FROM ARTIFACTORY...")
	
	var artifactPath string
	var downloadList []string
	var resultMsg string
	var err error
	var strI string
	var extString string
	var task string
	downloadUri = strings.ToLower(downloadUri)

	fileName 	 := common.ParseUriForFilename(downloadUri)
	downloadPath := strings.TrimSuffix(downloadUri, fileName) // still has slash
	ext 		 := filepath.Ext(fileName)
	imageName    := common.ParseFilenameForImageName(fileName)

	common.LogTxtHandler().Debug("File Name: " + fileName)
	common.LogTxtHandler().Debug("Download Path: " + downloadPath)
	common.LogTxtHandler().Debug("Extension of File: " + ext)
	common.LogTxtHandler().Debug("Image Name: " + imageName)

	// Create imageName-based folder under Output Dir to house file downloads
	outputDir      = common.CheckAddSlashToPath(outputDir)
	newOutputDir  := outputDir + imageName
	common.LogTxtHandler().Debug("Original Output Directory: " + outputDir)
	common.LogTxtHandler().Debug("New Output Directory: " + newOutputDir)

	util.OutputDir = newOutputDir          // Setting subdir as the new output directory
	// Check for output directory and create if it doesn't exist
	_, err = os.Stat(newOutputDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(newOutputDir, 0755)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error creating directory: " + newOutputDir + " - " + strErr)
		} else {
			common.LogTxtHandler().Info("Successfully created directory: " + newOutputDir)
		}
	}
	
	if ext == ".ova" {
		resultMsg, err = operations.RetrieveArtifact(downloadUri)

		if err != nil {
			common.LogTxtHandler().Error(resultMsg) // Will contain "Error" with additional info
		}
	} else if ext == ".ovf" {
		// Download OVF and assoc files
		downloadList = []string{".ovf", ".mf"}
		for _, item := range downloadList {
			artifactPath = downloadPath + imageName + item // builds URI path for each expected file type
			resultMsg, err = operations.RetrieveArtifact(artifactPath)
		}
		if err != nil {
			common.LogTxtHandler().Error("Download OVF files: " + resultMsg) // Will contain "Error" with additional info
		} else {
			common.LogTxtHandler().Info("Download OVF files: " + resultMsg)  // "Completed file download"
		}

		// Download Disk File(s)
		for i := 1; i < 15; i++ {   // allowing possibility of up to 15 disk files
			strI = strconv.Itoa(i)
			checkFile := imageName + "-disk-" + strI + ".vmdk"
			statusCode, err := operations.GetArtifact(downloadPath + checkFile)
			if statusCode == "200" {
				// If we found the artifact, download it...
				resultMsg, err = operations.RetrieveArtifact(downloadPath + checkFile)
				if err != nil {
					common.LogTxtHandler().Error(resultMsg) // Will contain "Error" with additional info
				}
			} else {
				common.LogTxtHandler().Info("End of OVF disk file checks.")
				break
			}
		}
	} else if ext == ".vmtx" {
		// Download known, static VMTX files
		downloadList = []string{".nvram", ".vmsd", ".vmtx", ".vmxf"}
		for _, item := range downloadList {
			artifactPath = downloadPath + imageName + item // builds URI path for each expected file type
			resultMsg, err = operations.RetrieveArtifact(artifactPath)
		}
		if err != nil {
			common.LogTxtHandler().Error("Download VMTX files: " + resultMsg) // Will contain "Error" with additional info
		} else {
			common.LogTxtHandler().Info("Download VMTX files: " + resultMsg)  // "Completed file download"
		}

		// Download Disk File(s)
		checkFile := imageName + ".vmdk"
		task = "Unnumbered virtual disk file check"
		resultMsg, err = operations.CheckFileAndDownload(checkFile, downloadPath, task)
		if err != nil {
			common.LogTxtHandler().Error(resultMsg)
		}
		
		// Loop for virtual disk files ----------------------------->
		extString = ".vmdk"
		task      = "Numbered virtual disk file check"
		resultMsg, err = operations.CheckFileLoopAndDownload(imageName, downloadPath, extString, task)
		if err != nil {
			common.LogTxtHandler().Error(resultMsg)
		}
		
		// Loop for disk -ctk files ----------------------------->
		checkFile = imageName + "-ctk.vmdk"
		task = "Unnumbered virtual ctk disk file check"
		resultMsg, err = operations.CheckFileAndDownload(checkFile, downloadPath, task)
		if err != nil {
			common.LogTxtHandler().Error(resultMsg)
		}
		
		extString = "-ctk.vmdk"
		task      = "Numbered virtual ctk disk file check"
		resultMsg, err = operations.CheckFileLoopAndDownload(imageName, downloadPath, extString, task)
		if err != nil {
			common.LogTxtHandler().Error(resultMsg)
		}
		
		// Loop for VM data disk (-flat) files ----------------------------->
		checkFile = imageName + "-flat.vmdk"
		task = "Unnumbered virtual data disk file check"
		resultMsg, err = operations.CheckFileAndDownload(checkFile, downloadPath, task)
		if err != nil {
			common.LogTxtHandler().Error(resultMsg)
		}
		
		extString = "-flat.vmdk"
		task      = "Numbered virtual data disk file check"
		resultMsg, err = operations.CheckFileLoopAndDownload(imageName, downloadPath, extString, task)
		if err != nil {
			common.LogTxtHandler().Error(resultMsg)
		}
		
		// Download associated vmware.log, if it exists  -------------------->
		checkFile = "vmware.log"
		task = "vmware.log file check"
		resultMsg, err = operations.CheckFileAndDownload(checkFile, downloadPath, task)
		if err != nil {
			common.LogTxtHandler().Error(resultMsg)
		}
	}

	// We are ignoring any potential .scoreboard and .hlog files that may exist
	// They are not necessary for our activities. 
	
	return "Completed"
	// Maybe future option to move here?
	// should convert AFTER copy - files are smaller, fewer
}