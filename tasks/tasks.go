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

func GetImageDetails(serverApi, token, artifName, ext string, kvProps []string) (string, string, string, string, error) {
	util.ServerApi = serverApi
	util.Token     = token
	var artifactUri string
	var strErr string

	common.LogTxtHandler().Debug(">>> GETTING IMAGE DETAILS...")
	common.LogTxtHandler().Debug("Getting artifacts by name...")
	listArtifacts, err := search.GetArtifactsByName(artifName)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error getting list of matching artifacts - " + strErr)
		return "", "", "", "", err
	}

	common.LogTxtHandler().Debug("Filtering list of artifacts by file type...")
	listByFileType, err := search.FilterListByFileType(ext, listArtifacts)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error filtering artifacts by file type - " + strErr)
		return "", "", "", "", err
	}

	if len(listByFileType) == 1 {
		// if just one artifact, we'll return it
		common.LogTxtHandler().Debug("List of artifacts contains one value...")
		artifactUri = listByFileType[0]

		common.LogTxtHandler().Debug("Artifact found: " + artifactUri)

	} else if len(listByFileType) > 1 && len(kvProps) != 0 {
		common.LogTxtHandler().Debug("List of artifacts contains multiple values...")
		common.LogTxtHandler().Debug("Filtering list of artifacts by properties...")
		artifactUri, err = operations.FilterListByProps(listByFileType, kvProps)
		if err != nil {
			strErr = fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error filtering artifacts by file type - " + strErr)
		}
	} else {
		// if no props passed, but more than one artif is in list, return latest
		common.LogTxtHandler().Debug("List of artifacts contains multiple values...")
		common.LogTxtHandler().Debug("Returning latest...")
		artifactUri, err = operations.GetLatestArtifactFromList(listByFileType)
		if err != nil {
			strErr = fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error getting latest artifact from list - " + strErr)
			return "", "", "", "", err
		}
		common.LogTxtHandler().Debug("Artifact found: " + artifactUri)
	}

	if artifactUri != "" {
		common.LogTxtHandler().Debug("Getting artifact name...")
		artifactName := operations.GetArtifactNameFromUri(artifactUri)
		common.LogTxtHandler().Debug("Artifact Name: " + artifactName)
		
		common.LogTxtHandler().Debug("Getting creation date for artifact...")
		createDate, err := operations.GetCreateDate(artifactUri)
		if err != nil {
			strErr = fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Unable to get create date of artifact - " + strErr)
			return "", "", "", "", err
		}
		common.LogTxtHandler().Debug("Creation Date is: " + createDate)
	
		common.LogTxtHandler().Debug("Getting download URI for artifact...")
		downloadUri, err := operations.GetDownloadUri(artifactUri)
		if err != nil {
			strErr = fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Unable to get download URI - " + strErr)
			return "", "", "", "", err
		}
		common.LogTxtHandler().Debug("Download URI: " + downloadUri)
		
		return artifactUri, artifactName, createDate, downloadUri, nil
	} else {
		return "", "", "", "", err
	}
}

// Must have Artifactory instance licensed at Pro or higher, access to create/remove repos and artifacts
func SetupTest(serverApi, token, testArtifactPath string, kvProps []string, uploadArtifact bool) (string, error) {
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
		downloadUri, err := operations.UploadFile(testArtifactPath, testRepoPath)
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

func UploadGeneralArtifact(serverApi, token, sourcePath, artifPath, fileName, folderName string) (string, error) {
	// Single file
	// 'folderName' may be the same as the image name, if wanting to place the artifact with a given image
	util.ServerApi = serverApi
	util.Token	   = token

	common.LogTxtHandler().Info(">>> Beginning validation and upload of "  + fileName)
	sourcePath = common.CheckAddSlashToPath(sourcePath)
	artifPath  = common.CheckAddSlashToPath(artifPath)

	common.LogTxtHandler().Debug("Source Path: " + sourcePath)
	common.LogTxtHandler().Debug("Artifact Path: " + artifPath)
	
	items, _ := os.ReadDir(sourcePath)            // Get list of files in Dir to check against our file
	result, err := operations.CheckFileAndUpload(items, sourcePath, artifPath, fileName, folderName)
	
	if result == "Success" {
		common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
		return result, nil
	} else if result == "Failed" && err == nil {
		common.LogTxtHandler().Info("File not found.")
		return result, nil
	} else {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
		return result, err
	}
}

func DownloadGeneralArtifact(serverApi, token, outputDir, artifPath, fileName, task string) (string, error) {
	util.ServerApi = serverApi
	util.Token	   = token
	util.OutputDir = outputDir
	common.LogTxtHandler().Info(">>> Beginning validation and download of "  + fileName)

	serverApi = common.FormatServerForDownloadUri(serverApi)
	serverApi = common.TrimEndSlashUrl(serverApi)
	downloadPath := serverApi + artifPath
	downloadPath = common.CheckAddSlashToPath(downloadPath)

	common.LogTxtHandler().Debug("Server API: " + serverApi)
	common.LogTxtHandler().Debug("Download Path: " + downloadPath)

	result, err := operations.CheckFileAndDownload(fileName, downloadPath, task)
	if result == "Failed" {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error downloading file: " + fileName + " - " + strErr)
		return result, err
	} else {
		common.LogTxtHandler().Debug("Successfully downloaded file: " + fileName)
		return result, nil
	}
}

func UploadArtifacts(serverApi, token, imageType, imageName, sourceDir, targetDir string) (string) {
	// Image files will placed in a folder named after the image, so no need to define a folder specifically for the image
	// targetDir --> /repo/ --> files will be in path: /repo/image1234/image1234.ova, for example

	// sourceDir ex: c:\\lab\\image_name or /lab/image_name - We'll check for/add ending slash if needed
	// targetDir ex: /repo-name/folder - We'll check for/add ending slash if needed
	util.ServerApi = serverApi
	util.Token	   = token
	var fileName string
	var err error
	var fileTypes, failedFiles, notFoundFiles []string
	imageType = strings.ToLower(imageType)

	if imageName != "" && sourceDir != "" && targetDir != "" {
		common.LogTxtHandler().Debug("UPLOADING NEW ARTIFACTS TO ARTIFACTORY...")
		newSourceDir := common.CheckAddSlashToPath(sourceDir)  // makes sure ending slash exists
		newTargetDir := common.CheckAddSlashToPath(targetDir)
		items, _ := os.ReadDir(sourceDir)

		if imageType == "ova" {
			fileName = imageName + ".ova"
			common.LogTxtHandler().Info("Searching for File Name: " + fileName)

			result, err := operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
			common.LogTxtHandler().Info(result)

			if result == "Failed" && err == nil {
				common.LogTxtHandler().Error("File: " + fileName + " not found.")
				return "File: " + fileName + " not found"
			} else if result == "Failed" && err != nil {
				strErr := fmt.Sprintf("%v", err)
				return "Error uploading file: " + fileName + " - " + strErr
			} else {
				return "End of upload process"
			}

		} else if imageType == "ovf" {
			fileTypes = []string{".ovf", ".mf"}
			for _, ft := range fileTypes {
				fileName = imageName + ft
				common.LogTxtHandler().Info("Searching for File Name: " + fileName)

				result, err := operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
				common.LogTxtHandler().Info(result)

				if result == "Failed" && err == nil {
					common.LogTxtHandler().Error("File: " + fileName + " not found.")
					notFoundFiles = append(notFoundFiles, fileName)
					break  // if one of the core files isn't found, we're stopping here
				} else if result == "Failed" && err != nil {
					strErr := fmt.Sprintf("%v", err)
					common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
					failedFiles = append(failedFiles, fileName)
					break  // if one of the core files isn't uploaded, we're stopping here
				} else {
					common.LogTxtHandler().Info("Successfully uploaded: " + fileName)
				}
			}

			// If there were no issues with the main files, we'll continue on
			// But we'll stop checking files if we run into issues
			if len(notFoundFiles) == 0 && len(failedFiles) == 0 {
				// Search and upload related OVF-based disk files
				for i := 1; i < 15; i++ {
					common.LogTxtHandler().Debug("Starting search for disk files. Up to 15 possible disks will be checked for.")
					strI := strconv.Itoa(i)
					fileName = imageName + "-disk" + strI + ".vmdk"					  // changing -disk-# to -disk#
					common.LogTxtHandler().Debug("Searching for File Name: " + fileName)

					result, err := operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
					common.LogTxtHandler().Info(result)

					if result == "Failed" && err == nil && i == 1 {  // if the first disk file isn't found, there's an issue
						common.LogTxtHandler().Error("Unable to locate first disk file: " + fileName)
						notFoundFiles = append(notFoundFiles, fileName)
						break
					} else if result == "Failed" && err == nil {     // if subsequent disk files aren't found, the machine probably doesn't have any more disks
						common.LogTxtHandler().Info("File: " + fileName + " not found.")
						break
					} else if result == "Failed" && err != nil {
						strErr := fmt.Sprintf("%v", err)
						common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
						failedFiles = append(failedFiles, fileName)
						break
					} else {
						common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
					}
				}	
			}

			if len(failedFiles) > 0 {
				return "Errors uploading one or more image files."
			} else if len(notFoundFiles) > 0 {
				return "One or more image files not found."
			} else {
				return "End of upload process"
			}
			
		} else if imageType == "vmtx" {
			var result string
			fileTypes = []string{".vmtx", ".nvram", ".vmsd", ".vmxf"}
			for _, ft := range fileTypes {
				fileName = imageName + ft
				common.LogTxtHandler().Info("Searching for File Name: " + fileName)

				result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
				common.LogTxtHandler().Info(result)

				if result == "Failed" && err == nil {
					common.LogTxtHandler().Error("File: " + fileName + " not found.")
					notFoundFiles = append(notFoundFiles, fileName)
					break  // if one of the core files isn't found, we're stopping here

				} else if result == "Failed" && err != nil {
					strErr := fmt.Sprintf("%v", err)
					common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
					failedFiles = append(failedFiles, fileName)
					break  // if one of the core files isn't uploaded, we're stopping here

				} else {
					common.LogTxtHandler().Info("Successfully uploaded: " + fileName)
				}
			}

			// If there were no issues with the main files, we'll continue on
			// But we'll stop checking files if we run into issues
			// Disk files can start their numbering at different spots/formats, so we have to be careful how we define an error
			if len(notFoundFiles) == 0 && len(failedFiles) == 0 {
				// Search and upload non-numbered virtual disk file
				common.LogTxtHandler().Debug("Starting search for disk files...")
				fileName = imageName + ".vmdk"
				common.LogTxtHandler().Debug("Searching for File Name: " + fileName)

				result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
				common.LogTxtHandler().Info(result)

				if err != nil {
					strErr := fmt.Sprintf("%v", err)
					common.LogTxtHandler().Error(result + " - " + strErr)
				}

				// Search and upload numbered disk files
				for i := 1; i < 15; i++ {
					common.LogTxtHandler().Debug("Starting search for numbered disk files. Up to 15 possible disks will be checked for.")
					strI := strconv.Itoa(i)
					fileName = imageName + "_" + strI + ".vmdk"
					common.LogTxtHandler().Debug("Searching for File Name: " + fileName)

					result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
					common.LogTxtHandler().Info(result)

					if result == "Failed" && err == nil { 
						common.LogTxtHandler().Debug("File: " + fileName + " not found.")
					} else if result == "Failed" && err != nil {
						strErr := fmt.Sprintf("%v", err)
						common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
						failedFiles = append(failedFiles, fileName)
						break
					} else {
						common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
					}
				}

				// Search and upload -ctk disk files ----------------------------------->
				// Search and upload non-numbered virtual disk file
				fileName = imageName + "-ctk.vmdk"
				common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
				
				result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
				common.LogTxtHandler().Info(result)

				if result == "Failed" && err == nil { 
					common.LogTxtHandler().Debug("File: " + fileName + " not found.")
				} else if result == "Failed" && err != nil {
					strErr := fmt.Sprintf("%v", err)
					common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
					failedFiles = append(failedFiles, fileName)
				} else {
					common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
				}

				// Search and upload numbered -ctk disk files
				for i := 1; i < 15; i++ {
					common.LogTxtHandler().Debug("Starting search for numbered CTK disk files. Up to 15 possible disks will be checked for.")
					strI := strconv.Itoa(i)
					fileName = imageName + "_" + strI + "-ctk.vmdk"
					common.LogTxtHandler().Debug("Searching for File Name: " + fileName)

					result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
					common.LogTxtHandler().Info(result)

					if result == "Failed" && err == nil { 
						common.LogTxtHandler().Debug("File: " + fileName + " not found.")
					} else if result == "Failed" && err != nil {
						strErr := fmt.Sprintf("%v", err)
						common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
						failedFiles = append(failedFiles, fileName)
						break
					} else {
						common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
					}
				}
				
				// Search and upload -flat disk files ----------------------------------->
				// Search and upload non-numbered virtual disk file
				fileName = imageName + "-flat.vmdk"
				common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
				
				result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
				common.LogTxtHandler().Info(result)

				if result == "Failed" && err == nil { 
					common.LogTxtHandler().Debug("File: " + fileName + " not found.")
				} else if result == "Failed" && err != nil {
					strErr := fmt.Sprintf("%v", err)
					common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
					failedFiles = append(failedFiles, fileName)
				} else {
					common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
				}

				// Search and upload numbered -flat disk files
				for i := 1; i < 15; i++ {
					common.LogTxtHandler().Debug("Starting search for numbered FLAT disk files. Up to 15 possible disks will be checked for.")
					strI := strconv.Itoa(i)
					fileName = imageName + "_" + strI + "-flat.vmdk"
					common.LogTxtHandler().Debug("Searching for File Name: " + fileName)

					result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
					common.LogTxtHandler().Info(result)

					if result == "Failed" && err == nil { 
						common.LogTxtHandler().Debug("File: " + fileName + " not found.")
					} else if result == "Failed" && err != nil {
						strErr := fmt.Sprintf("%v", err)
						common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
						failedFiles = append(failedFiles, fileName)
						break
					} else {
						common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
					}
				}
				
				//------- Just covering our bases on other possible disk files that may be present ---------
				// Search [image]-00000#.vmdk files
				for i := 1; i < 15; i++ {
					common.LogTxtHandler().Debug("Starting search for -00000# disk files. Up to 15 possible disks will be checked for.")
					strI := strconv.Itoa(i)
					if i >= 1 && i < 10 {
						fileName = imageName + "-00000" + strI + ".vmdk"
						common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
					} else {
						fileName = imageName + "-0000" + strI + ".vmdk"
						common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
					}

					result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
					common.LogTxtHandler().Info(result)

					if result == "Failed" && err == nil { 
						common.LogTxtHandler().Debug("File: " + fileName + " not found.")
					} else if result == "Failed" && err != nil {
						strErr := fmt.Sprintf("%v", err)
						common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
						failedFiles = append(failedFiles, fileName)
						break
					} else {
						common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
					}
				}

				// Search [image]-00000#-ctk.vmdk files
				for i := 1; i < 15; i++ {
					common.LogTxtHandler().Debug("Starting search for -00000#-ctk disk files. Up to 15 possible disks will be checked.")
					strI := strconv.Itoa(i)
					if i >= 1 && i < 10 {
						fileName = imageName + "-00000" + strI + "-ctk.vmdk"
						common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
					} else {
						fileName = imageName + "-0000" + strI + "-ctk.vmdk"
						common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
					}

					result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
					common.LogTxtHandler().Info(result)

					if result == "Failed" && err == nil { 
						common.LogTxtHandler().Debug("File: " + fileName + " not found.")
					} else if result == "Failed" && err != nil {
						strErr := fmt.Sprintf("%v", err)
						common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
						failedFiles = append(failedFiles, fileName)
						break
					} else {
						common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
					}
				}

				// Search [image]-00000#-delta.vmdk files // should only exist if there's a snapshot, including just in case
				for i := 1; i < 15; i++ {
					common.LogTxtHandler().Debug("Starting search for -00000#-delta disk files. Up to 15 possible disks will be checked for.")
					strI := strconv.Itoa(i)
					if i >= 1 && i < 10 {
						fileName = imageName + "-00000" + strI + "-delta.vmdk"
						common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
					} else {
						fileName = imageName + "-0000" + strI + "-delta.vmdk"
						common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
					}

					result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
					common.LogTxtHandler().Info(result)

					if result == "Failed" && err == nil { 
						common.LogTxtHandler().Debug("File: " + fileName + " not found.")
					} else if result == "Failed" && err != nil {
						strErr := fmt.Sprintf("%v", err)
						common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
						failedFiles = append(failedFiles, fileName)
						break
					} else {
						common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
					}
				}

				// Search [image]-00000#-flat.vmdk files
				for i := 1; i < 15; i++ {
					common.LogTxtHandler().Debug("Starting search for -00000#-flat disk files. Up to 15 possible disks will be checked for.")
					strI := strconv.Itoa(i)
					if i >= 1 && i < 10 {
						fileName = imageName + "-00000" + strI + "-flat.vmdk"
						common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
					} else {
						fileName = imageName + "-0000" + strI + "-flat.vmdk"
						common.LogTxtHandler().Debug("Searching for File Name: " + fileName)
					}

					result, err = operations.CheckFileAndUpload(items, newSourceDir, newTargetDir, fileName, imageName)
					common.LogTxtHandler().Info(result)

					if result == "Failed" && err == nil { 
						common.LogTxtHandler().Debug("File: " + fileName + " not found.")
					} else if result == "Failed" && err != nil {
						strErr := fmt.Sprintf("%v", err)
						common.LogTxtHandler().Error("Error uploading file: " + fileName + " - " + strErr)
						failedFiles = append(failedFiles, fileName)
						break
					} else {
						common.LogTxtHandler().Info("Successfully uploaded file: " + fileName)
					}
				}

				if len(failedFiles) > 0 {
					return "Errors uploading one or more image files."
				} else if len(notFoundFiles) > 0 {
					return "One or more image files not found."
				} else {
					return "End of upload process"
				}
			} else {
				return "Unable to find and/or upload one or more files."
			}
			
		} else {
			common.LogTxtHandler().Error("Unsupported or blank image type. Supported image types are OVA, OVF, and VMTX.")
			if imageType != "" {
				return "Unsupported image type"
			} else {
				return "Image type is blank"
			}
		}
	} else {
		common.LogTxtHandler().Error("One or more required inputs have not been provided.")
		common.LogTxtHandler().Error("IMAGE TYPE: " + imageType)
		common.LogTxtHandler().Error("IMAGE NAME: " + imageName)
		common.LogTxtHandler().Error("SOURCE DIR: " + sourceDir)
		common.LogTxtHandler().Error("TARGET DIR:" + targetDir)
		return "Missing required inputs"
	}
}

func SetProps(serverApi, token, artifUri string, kvProps []string) (string, error) {
	util.ServerApi = serverApi
	util.Token	   = token

	common.LogTxtHandler().Debug("UPDATING PROPERTIES OF ARTIFACT...")

	statusCode, err := operations.SetArtifactProps(artifUri, kvProps)
	common.LogTxtHandler().Debug("Status code of Set Artifact Properties task: " + statusCode)

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
	// Takes in download URI that corresponds to OVA, OVF, or VMTX file in Artifactory; 
	// Will then determine other expected associated artifacts and download those as well
		// Appends an incrementing numeric value (string; up to 15) to disk type and checks for existance of disk file
		// At first occurrance of the file not being found, the check breaks and moves on
	// ** If planning to import image file into vCenter, make the output directory the destination datastore

	util.ServerApi = serverApi
	util.Token     = token
	
	common.LogTxtHandler().Info("DOWNLOADING ARTIFACT(S) FROM ARTIFACTORY...")
	
	var artifactPath string
	var downloadList []string
	var resultMsg string
	var err error
	var strI string
	var extString string
	var task string

	if downloadUri != "" && outputDir != "" {
		downloadUri   = strings.ToLower(downloadUri)
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
			err = os.MkdirAll(newOutputDir, 0755)   // Will create any directories in the given path if doesn't exist
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Error creating directory: " + newOutputDir + " - " + strErr)
			} else {
				common.LogTxtHandler().Info("Successfully created directory: " + newOutputDir)
			}
		}
		
		if ext == ".ova" {
			common.LogTxtHandler().Info("Image type identified as OVA. Downloading OVA file...")
			resultMsg, err = operations.RetrieveArtifact(downloadUri)
			if err != nil {
				common.LogTxtHandler().Error(resultMsg) // Will contain "Error" with additional info
			}

		} else if ext == ".ovf" {
			// Download OVF and assoc files
			common.LogTxtHandler().Info("Image type identified as OVF. Downloading OVF files...")
			downloadList = []string{".ovf", ".mf"}
			for _, item := range downloadList {
				artifactPath = downloadPath + imageName + item // builds URI path for each expected file type
				common.LogTxtHandler().Info("Downloading: " + artifactPath)
				resultMsg, err = operations.RetrieveArtifact(artifactPath)
				if err != nil {
					common.LogTxtHandler().Error("Error downloading " + artifactPath)
					common.LogTxtHandler().Error(resultMsg)
					break
				}
			}
			if err != nil {
				common.LogTxtHandler().Error("Download OVF files: " + resultMsg) // Will contain "Error" with additional info
			} else {
				common.LogTxtHandler().Info("Download OVF files: " + resultMsg)  // "Completed file download"
			}

			// We want the downloadList to fully complete before moving on, but if there were any errors,
			// we won't try to download the disk files.
			if err == nil {
				// Download Disk File(s)
				for i := 1; i < 15; i++ {   // allowing possibility of up to 15 disk files
					strI = strconv.Itoa(i)
					checkFile := imageName + "-disk" + strI + ".vmdk"					// changing from -disk-# to -disk#
					common.LogTxtHandler().Info("Checking for existance of disk file: " + checkFile)

					statusCode, err := operations.GetArtifact(downloadPath + checkFile)
					if statusCode == "200" {
						// If we found the artifact, download it...
						common.LogTxtHandler().Info("Disk file FOUND. Downloading...")
						resultMsg, err = operations.RetrieveArtifact(downloadPath + checkFile)
						if err != nil {
							common.LogTxtHandler().Error(resultMsg) // Will contain "Error" with additional info
						}
					} else {
						common.LogTxtHandler().Info("Disk file doesn't exist. Reached end of disk files.")
						common.LogTxtHandler().Info("End of OVF disk file checks.")
						break
					}
				}
			} else {
				common.LogTxtHandler().Error("Errors encountered. The remainder of the file download process will terminate.")
			}
		} else if ext == ".vmtx" {
			// Download known, static VMTX files
			common.LogTxtHandler().Info("Image type identified as VMTX. Downloading VMTX files...")
			downloadList = []string{".nvram", ".vmsd", ".vmtx", ".vmxf"}
			for _, item := range downloadList {
				artifactPath = downloadPath + imageName + item // builds URI path for each expected file type
				common.LogTxtHandler().Info("Downloading: " + artifactPath)

				resultMsg, err = operations.RetrieveArtifact(artifactPath)
				if err != nil {
					common.LogTxtHandler().Error("Error downloading " + artifactPath)
					common.LogTxtHandler().Error(resultMsg)
					break
				}
			}
			if err != nil {
				common.LogTxtHandler().Error("Download VMTX files: " + resultMsg) // Will contain "Error" with additional info
			} else {
				common.LogTxtHandler().Info("Download VMTX files: " + resultMsg)  // "Completed file download"
			}

			// We want the downloadList to fully complete before moving on, but if there were any errors,
			// we won't try to download the disk files.
			if err == nil {
				// Download Disk File(s)
				checkFile := imageName + ".vmdk"
				common.LogTxtHandler().Info("Checking for existance of disk file: " + checkFile)
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
				common.LogTxtHandler().Info("Checking for existance of disk file: " + checkFile)
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
				common.LogTxtHandler().Info("Checking for existance of disk file: " + checkFile)
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
				common.LogTxtHandler().Info("Checking for existance of file: " + checkFile)
				task = "vmware.log file check"
				resultMsg, err = operations.CheckFileAndDownload(checkFile, downloadPath, task)
				if err != nil {
					common.LogTxtHandler().Error(resultMsg)
				}
			} else {
				common.LogTxtHandler().Error("Errors encountered. The remainder of the file download process will terminate.")
				return "File download failed"
			}
			// We are ignoring any potential .scoreboard and .hlog files that may exist
			// They are not necessary for the imaging process.
		}
		return "End of download process"
	} else {
		common.LogTxtHandler().Error("One or more required inputs have not been provided.")
		common.LogTxtHandler().Error("DOWNLOAD URI: " + downloadUri)
		common.LogTxtHandler().Error("OUTPUT DIRECTORY: " + outputDir)
		return "Missing required inputs"
	}
}