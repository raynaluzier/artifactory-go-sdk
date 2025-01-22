package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/raynaluzier/artifactory-go-sdk/common"
	"github.com/raynaluzier/artifactory-go-sdk/util"
)

type Contents struct {
	Child	 		string
	IsFolder		bool
}

type artifJson struct {
	Repo			string 	`json:"repo"`
	Path			string	`json:"path"`
	Created			string	`json:"created"`
	CreatedBy		string	`json:"createdBy"`
	LastModified	string	`json:"lastModified"`
	ModifiedBy		string	`json:"modifiedBy"`
	LastUpdated		string	`json:"lastUpdated"`
	DownloadUri 	string 	`json:"downloadUri"`
	MimeType 		string	`json:"mimeType"`
	Size			string	`json:"size"`
	Checksums	struct {
		Sha1		string	`json:"sha1"`
		Md5			string	`json:"md5"`
		Sha256		string	`json:"sha256"`
	}	`json:"checksums"`
	OriginalChecksums	struct {
		Sha1		string	`json:"sha1"`
		Md5			string	`json:"md5"`
		Sha256		string	`json:"sha256"`				
	}   `json:"originalChecksums"`
	Uri 			string	`json:"uri"`
}

var request *http.Request
var err error
var foundPaths []string


func ListRepos() ([]string, error) {
	var listRepos []string
	bearer := common.SetBearer(util.Token)
	requestPath := util.ServerApi + "/repositories"
	
	common.LogTxtHandler().Info(">>> Getting list of available repos...")
	common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + requestPath)

	request, err = http.NewRequest("GET", requestPath, nil)
	request.Header.Add("Authorization", bearer)
	
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error on API response from 'GET' " + requestPath + " - " + strErr)
	} else {
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		common.LogTxtHandler().Debug("REQUEST RESPONSE: " + string(body))

		// JSON return is an array of strings '[{"key":"repo_name1, "type":"LOCAL"...}, {"key":"repo_name2"}...]'
		type reposJson struct {
			Key 		string	`json:"key"`
			Description	string	`json:"description"`
			Type		string	`json:"type"`
			Url 		string	`json:"url"`
			PackageType	string	`json:"packageType"`
		}

		var jsonData []reposJson
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Could not unmarshal response - " + strErr)
		}

		if len(jsonData) != 0 {
			for _, k := range jsonData {
				listRepos = append(listRepos, k.Key)
				common.LogTxtHandler().Debug("FOUND REPO: " + k.Key)
			}
			return listRepos, nil
		} else {
			err := errors.New("No repos found")
			common.LogTxtHandler().Warn("No repos found")
			return nil, err
		}
	}
	return listRepos, nil
}

func GetDownloadUri(artifUri string) (string, error) {
	bearer := common.SetBearer(util.Token)
	var downloadUri string
	common.LogTxtHandler().Info(">>> Getting Download URI from Artifact URI: " + artifUri + "...")

	if (artifUri != "") {
		common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + artifUri)
		request, err = http.NewRequest("GET", artifUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error on response. " + strErr)
			return "", err
		} else {
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			common.LogTxtHandler().Debug("REQUEST RESPONSE: " + string(body))

			var jsonData *artifJson
			err = json.Unmarshal(body, &jsonData)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Could not unmarshal response - " + strErr)
			}

			if jsonData.DownloadUri != "" {
				downloadUri = jsonData.DownloadUri
				common.LogTxtHandler().Info("DOWNLOAD URI RETRIEVED: " + downloadUri)
				return downloadUri, nil
			} else {
				err = errors.New("There is no download URI for the artifact.")
				common.LogTxtHandler().Warn("There is no download URI for the artifact.")
				return "", err
			}
		}
	} else {
		err := errors.New("No artifact URI was provided.")
		common.LogTxtHandler().Error("Unable to get artifact's download URI without the artifact's URI.")
		return "", err
	}
}

func GetCreateDate(artifUri string) (string, error) {
	bearer := common.SetBearer(util.Token)
	var createdDate string
	common.LogTxtHandler().Info(">>> Getting Create Date for Artifact: " + artifUri + "...")

	if (artifUri != "") {
		common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + artifUri)
		request, err = http.NewRequest("GET", artifUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error on response. " + strErr)
			return "", err
		} else {
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			common.LogTxtHandler().Debug("REQUEST RESPONSE: " + string(body))

			var jsonData *artifJson
			err = json.Unmarshal(body, &jsonData)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Could not unmarshal response - " + strErr)
			}

			if jsonData.Created != "" {
				createdDate = jsonData.Created
				common.LogTxtHandler().Debug("CREATE DATE RETRIEVED: " + createdDate)
				return createdDate, nil
			} else {
				err = errors.New("There is no create date for the artifact.")
				common.LogTxtHandler().Warn("There is no create date for the artifact.")
				return "", err
			}
		}
	} else {
		err := errors.New("No artifact URI was provided.")
		common.LogTxtHandler().Error("Unable to get artifact's created date without the artifact's URI.")
		return "", err
	}
}

func GetArtifactNameFromUri(artifUri string) (string) {
	// Parses artifact name from artifact's URI
	common.LogTxtHandler().Info(">>> Getting Artifact Name from URI...")
	fileName := path.Base(artifUri)
	ext := filepath.Ext(fileName)
	artifactName := strings.TrimSuffix(fileName, ext)
	common.LogTxtHandler().Info("ARTIFACT NAME: " + artifactName)
	return artifactName
}

func RetrieveArtifact(downloadUri string) (string, error) {
	// Gets the artifact via provided Download URI and copies it to the output directory specified in
	// the environment variables file
	var outputDir string
	bearer := common.SetBearer(util.Token)

	common.LogTxtHandler().Info(">>> Retrieving Artifact by Download URI: " + downloadUri + "...")
	// If no output directory path was provided, the artifact file will be downloaded to the user's HOME directory
	if len(util.OutputDir) != 0 {
		outputDir = common.EscapeSpecialChars(util.OutputDir)  // Ensure special characters are escaped
		outputDir = common.CheckAddSlashToPath(outputDir) // Ensure path ends with appropriate slash type

	} else {  // No output directory specified...
		common.LogTxtHandler().Warn("*** No output directory provided; output will be user's home directory.")
		outputDir, err = os.UserHomeDir()
		
		if err != nil {
			common.LogTxtHandler().Error("Unable to get user's home directory.")
		} else {
			common.LogTxtHandler().Debug("User's home directory is: " + outputDir)
			outputDir = common.CheckAddSlashToPath(outputDir)
		}
	}

	if downloadUri != "" {
		common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + downloadUri)
		request, err = http.NewRequest("GET", downloadUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error on response. " + strErr)
			return "", err
		} else {
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			common.LogTxtHandler().Debug("REQUEST RESPONSE/ FILE CONTENTS: " + string(body))
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Error reading response body. " + strErr)
			}

			if response.StatusCode == 404 {
				err := errors.New("File not found.")
				common.LogTxtHandler().Error("File not found. File download failed.")
				return "File download failed.", err
			} else {
				// Create file name from download URI path of artifact
				fileUrl, err := url.Parse(downloadUri)
				if err != nil {
					strErr := fmt.Sprintf("%v\n", err)
					common.LogTxtHandler().Error("Unable to determine file path. " + strErr)
				}

				// Get the file name from the path
				path := fileUrl.Path
				segments := strings.Split(path, "/")
				fileName := segments[len(segments)-1]

				// Creates the file at the defined path
				// Will overwrite the file if it already exists
				newFile, err := os.Create(outputDir + fileName)   
				if err != nil {
					strErr := fmt.Sprintf("%v\n", err)
					common.LogTxtHandler().Error("Error creating file at target location. " + strErr)
					return "Error creating file at target location.", err
				}
				err = os.WriteFile(outputDir + fileName, body, 0644)
				if err != nil {
					strErr := fmt.Sprintf("%v\n", err)
					common.LogTxtHandler().Error("Error downloading file to target location. " + strErr)
					return "Error downloading file to target location.", err
				}
				defer newFile.Close()
			}
		}
	} else {
		err := errors.New("No download URI was provided. Unable to download the artifact without the download URI.")
		common.LogTxtHandler().Error("No download URI was provided. Unable to download the artifact without the download URI.")
		return "Error: File download failed.", err
	}

	return "Completed file download", nil
}

func UploadFile(sourcePath, targetPath, fileSuffix string) (string, error) {
	//** TO DO: Option to get previous 'version' and increment

	var downloadUri string
	var filePath string
	var fileName string
	var found bool
	bearer := common.SetBearer(util.Token)
	separater := "-"										  // If adding a file suffix (like date, version, etc), use this separater between filename and suffix
	trimmedBase := util.ServerApi[:len(util.ServerApi)-4]               // Removing '/api' from base URI

	common.LogTxtHandler().Info(">>> Uploading File From: " + sourcePath + "...")

	if len(sourcePath) != 0 && targetPath != "" { 
		// We need to ensure the provided source path/file are valid and exist
		if len(path.Ext(sourcePath)) != 0 {		// Ensures file with extension exists in source path
			common.LogTxtHandler().Debug("Escaping special characters in source/target paths.")
			sourcePath = common.EscapeSpecialChars(sourcePath)
			targetPath = common.EscapeSpecialChars(targetPath)
			common.LogTxtHandler().Debug("Checking end slash on target path and adding if necessary.")
			targetPath = common.CheckAddSlashToPath(targetPath)
			
			// Determine source filename and source file path by platform type
			common.LogTxtHandler().Debug("Checking source path type...")
			winPath := common.CheckPathType(sourcePath)
			if winPath == true {
				common.LogTxtHandler().Debug("Source path type identified as Windows-based.")
				segments := strings.Split(sourcePath, "\\")	  	  // Split source path into segments
				fileName = segments[len(segments)-1]			  // Determine filename from path
				filePath = sourcePath[:len(sourcePath)-len(fileName)]  // Determine path without filename
			} else {   // Unix path
				common.LogTxtHandler().Debug("Source path type identified as Unix-based.")
				segments := strings.Split(sourcePath, "/")	 	  // Split source path into segments
				fileName = segments[len(segments)-1]              // Determine filename from path
				filePath = sourcePath[:len(sourcePath)-len(fileName)]  // Determine path without filename				
			}
			
			// Get all files in the provided source directory
			common.LogTxtHandler().Debug("Reading all files in source directory...")
			filesInDirectory, err := os.ReadDir(filePath)
			if err != nil {
				return "", err
			}
			
			// For each file in the source directory, do a case insensitive file name comparison for a match
			// As Artifactory cares about case here, we want to make sure the filename supplied matches the case of the filename that actually exists in the source path
			common.LogTxtHandler().Debug("Performing case insensitive search for file...")
			found = false															// Initially set to false; then if found, turns true
			for _, file := range filesInDirectory {
				isSameStr := common.StringCompare(fileName, file.Name())            // Filename from provided source path vs. filename pulled directly from source path
				if isSameStr == true {												// If true, we know files are the same
					common.LogTxtHandler().Debug("FILE FOUND. Checking case. Will update to match case if necessary.")
					found = true													// Mark that we found a matching file
					isExactStr, err := common.SearchForExactString(file.Name(), fileName)  // Now, checks if cases matches
					if err != nil {
						strErr := fmt.Sprintf("%v\n", err)
						common.LogTxtHandler().Error("Error searching for exact string: " + strErr)
					}
					
					if isExactStr == false {										// Files are the same, but provided and actual cases are different
						fileName = file.Name() 										// Set the provided filename to match the actual filename so we'll use to the correct case
					}
					break
				}
			}

			// If we couldn't find a matching file at all, then we throw an error
			if found == false {
				err := errors.New("Unable to validate existance of source file. Source file doesn't exist.")
				common.LogTxtHandler().Error("Unable to validate existance of source file. Source file doesn't exist.")
				return "", err
			}
			
			// We now have a validated source path and filename
			// Set target filename = filename + fileSuffix (if not blank)
			if len(fileSuffix) != 0 || fileSuffix != "" {							// If a file suffix (like version, date, etc) was provided...
				fileExt := path.Ext(fileName)										// Returns .[ext]
				fmt.Println(fileName)
				justName := strings.TrimSuffix(fileName, fileExt)							// Trim off extension
				fileName = justName + separater + fileSuffix + fileExt
			}   // If blank, then the original filename will be used
			
			newArtifactPath := trimmedBase + targetPath + fileName                  // Forms: http://artifactory_base_api_url/repo-key/folder/artifact.txt
			data := strings.NewReader("@/" + sourcePath)                            // Formats the payload appropriately
			fmt.Println(data)
			common.LogTxtHandler().Debug("REQUEST: Sending 'PUT' request to: " + newArtifactPath)
			
			request, err = http.NewRequest("PUT", newArtifactPath, data)
			request.Header.Add("Authorization", bearer)
	
			client := &http.Client{}
			response, err := client.Do(request)

			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Error on response. " + strErr)
				return "", err
			} else {
				defer response.Body.Close()
				body, err := io.ReadAll(response.Body)
				common.LogTxtHandler().Debug("REQUEST RESPONSE: " + string(body))
		
				var jsonData *artifJson
				err = json.Unmarshal(body, &jsonData)
				if err != nil {
					strErr := fmt.Sprintf("%v\n", err)
					common.LogTxtHandler().Error("Could not unmarshal response - " + strErr)
				}
		
				if jsonData.DownloadUri != "" {
					downloadUri = jsonData.DownloadUri
					common.LogTxtHandler().Debug("DOWNLOAD URI RETRIEVED: " + downloadUri)
					return downloadUri, nil
				} else {
					err = errors.New("There is no download URI for the artifact")
					common.LogTxtHandler().Warn("There is no download URI for the artifact")
					return "", err
				}
			}
		} else {
			err = errors.New("No file extension found in source path. Ensure source includes path and source file with extension.")
			common.LogTxtHandler().Error("No file extension found in source path. Ensure source includes path and source file with extension.")
			return "", err
		}
	} else {
		err := errors.New("Cannot upload file without source path/file, target path, and artifact file name")
		common.LogTxtHandler().Error("Supplied source path: " + sourcePath + ", target path: " + targetPath)
		common.LogTxtHandler().Error("Cannot upload file without source path/file, target path, and artifact file name")
		return "", err
	}
}

func DeleteArtifact(artifUri string) (string, error) {
	bearer := common.SetBearer(util.Token)
	common.LogTxtHandler().Info(">>> Deleting Artifact: " + artifUri + "...")

	if artifUri != "" { 
		common.LogTxtHandler().Debug("REQUEST: Sending 'DELETE' request to: " + artifUri)
		request, err = http.NewRequest("DELETE", artifUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error on response. " + strErr)
			return "", err
		} else {
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Error getting response body. " + strErr)
			}
			common.LogTxtHandler().Info("REQUEST RESPONSE: " + string(body))
			
			if response.StatusCode == 204 {
				common.LogTxtHandler().Info("Request completed successfully")
				statusCode = "204"
			} else {
				common.LogTxtHandler().Info("Unable to complete request")
				statusCode = "404"
			}
		}
	} else {
		err := errors.New("Unable to DELETE item without artifact URI.")
		common.LogTxtHandler().Error("Supplied artifact path is: " + artifUri)
		common.LogTxtHandler().Error("Unable to DELETE item without artifact URI.")
		return "", err
	}
	return statusCode, nil
}

func GetLatestArtifactFromList(list []string) (string, error) {
	var latestItem string
	var dateMap []map[string]string

	for item := 0; item < len(list); item++ {
		addMap := make(map[string]string)
		created, err := GetCreateDate(list[item])
		if err != nil {
			common.LogTxtHandler().Error("Error getting created date.")
		}
		common.LogTxtHandler().Info("CREATED DATE RETRIEVED:" + created)
		addMap["artifact"] = list[item]
		addMap["created"] = created
		dateMap = append(dateMap, addMap)
	}

	sort.Slice(dateMap, func(i, j int) bool {
		return dateMap[i]["created"] < dateMap[j]["created"]
	})

	latest := len(dateMap) - 1
	latestItem = dateMap[latest]["artifact"]
	common.LogTxtHandler().Info("LATEST ITEM: " + latestItem)
	return latestItem, nil
}

func GetArtifact(downloadUri string) (string, error) {
	// Checks to see if artifact exists
	bearer := common.SetBearer(util.Token)
	common.LogTxtHandler().Info(">>> Getting artifact: " + downloadUri)

	if downloadUri != "" {
		request, err = http.NewRequest("GET", downloadUri, nil)
		common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + downloadUri)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error on response. " + strErr)
			return "", err
		} else {
			defer response.Body.Close()

			// If the request is successful, it will simply return a status code of 200
			if response.StatusCode == 200 {
				common.LogTxtHandler().Info("Request completed successfully")
				statusCode = "200"
			} else {
				// If the request fails, it will return a status code of 404
				common.LogTxtHandler().Info("Artifact not found.")
				statusCode = "404"
			}
		}
	} else {
		err := errors.New("Unable to get artifact without artifact's download URI.")
		common.LogTxtHandler().Error("No download URI provided. Unable to get artifact without artifact's download URI.")
		return "", err
	}

	if err != nil {
		common.LogTxtHandler().Error("Unable to parse URL")
		return "", err
	}

	return statusCode, nil
}

func CheckFileAndUpload(items []os.DirEntry, sourceDir, targetDir, fileName, imageName string) (string, error) {
	// sourceDir ex: c:\\lab\\ or /lab/ - assumes ending slash
	// targetDir ex: /repo-name/folder/ - assumes ending slash
	var sourcePath, targetPath string
	for _, item := range items {
		if item.Name() == fileName {
			sourcePath = sourceDir + fileName
			targetPath = targetDir + imageName + "/"
			suff := ""
			downloadUri, err := UploadFile(sourcePath, targetPath, suff)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Error uploading: " + fileName + " to: " + targetPath + " - " + strErr)
				return "Failed", err
			} else {
				common.LogTxtHandler().Info("File: " + fileName + " uploaded.")
				common.LogTxtHandler().Info("Download URI: " + downloadUri)
			}
		}
	}
	return "Success", nil
}

func CheckFileAndDownload(checkFile, downloadPath, task string) (string, error) {
	// checkFile - filename with extention
	// downloadPath - parsed Artifactory path to artifact without the artifact file name
	// task - what file check we are performing
	var resultMsg string
	statusCode, err := GetArtifact(downloadPath + checkFile)
	if statusCode == "200" {
		// If we found the artifact, download it...
		resultMsg, err = RetrieveArtifact(downloadPath + checkFile)
		if err != nil {
			common.LogTxtHandler().Error(resultMsg)
			return "Failed", err
		}
	} else {
		common.LogTxtHandler().Info("End of " + task)
	}
	return "Success", nil
}

func CheckFileLoopAndDownload(imageName, downloadPath, extString, task string) (string, error) {
	// imageName - name of image we'll use to construct the filename with
	// downloadPath - parsed Artifactory path to artifact without the artifact file name
	// extString - vSphere-based disk file extension - ex: ".vmdk", "-ctk.vmdk", "-flat.vmdk"
	// task - what file check we are performing
	var resultMsg, strI string
	for i := 1; i < 15; i++ {   // allowing possibility of up to 15 disk files
		strI = strconv.Itoa(i)
		checkFile := imageName + "_" + strI + extString
		statusCode, err := GetArtifact(downloadPath + checkFile)
		if statusCode == "200" {
			// If we found the artifact, download it...
			resultMsg, err = RetrieveArtifact(downloadPath + checkFile)
			if err != nil {
				common.LogTxtHandler().Error(resultMsg)
				return "Failed", err
			}
		} else {
			common.LogTxtHandler().Info("End of " + task)
			break
		}
	}
	return "Success", nil
}