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
	Child    string
	IsFolder bool
}

type artifJson struct {
	Repo         string `json:"repo"`
	Path         string `json:"path"`
	Created      string `json:"created"`
	CreatedBy    string `json:"createdBy"`
	LastModified string `json:"lastModified"`
	ModifiedBy   string `json:"modifiedBy"`
	LastUpdated  string `json:"lastUpdated"`
	DownloadUri  string `json:"downloadUri"`
	MimeType     string `json:"mimeType"`
	Size         string `json:"size"`
	Checksums    struct {
		Sha1   string `json:"sha1"`
		Md5    string `json:"md5"`
		Sha256 string `json:"sha256"`
	} `json:"checksums"`
	OriginalChecksums struct {
		Sha1   string `json:"sha1"`
		Md5    string `json:"md5"`
		Sha256 string `json:"sha256"`
	} `json:"originalChecksums"`
	Uri string `json:"uri"`
}

var foundPaths []string

func ListRepos() ([]string, error) {
	var listRepos []string
	bearer := common.SetBearer(util.Token)
	requestPath := util.ServerApi + "/repositories"

	common.LogTxtHandler().Info(">>> Getting list of available repos...")
	common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + requestPath)

	request, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", bearer)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error on API response from 'GET' " + requestPath + " - " + strErr)
		return nil, err
	} else {
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		common.LogTxtHandler().Debug("REQUEST RESPONSE: " + string(body))

		// JSON return is an array of strings '[{"key":"repo_name1, "type":"LOCAL"...}, {"key":"repo_name2"}...]'
		type reposJson struct {
			Key         string `json:"key"`
			Description string `json:"description"`
			Type        string `json:"type"`
			Url         string `json:"url"`
			PackageType string `json:"packageType"`
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

	if artifUri != "" {
		common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + artifUri)
		request, err := http.NewRequest("GET", artifUri, nil)
		if err != nil {
			return "", err
		}
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

	if artifUri != "" {
		common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + artifUri)
		request, err := http.NewRequest("GET", artifUri, nil)
		if err != nil {
			return "", err
		}
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

func GetArtifactNameFromUri(artifUri string) string {
	// Parses artifact name from artifact's URI
	common.LogTxtHandler().Info(">>> Getting Artifact Name from URI...")
	fileName := path.Base(artifUri)
	ext := filepath.Ext(fileName)
	artifactName := strings.TrimSuffix(fileName, ext)
	common.LogTxtHandler().Info("ARTIFACT NAME: " + artifactName)
	return artifactName
}

func RetrieveArtifact(downloadUri string) (string, error) {
	if downloadUri == "" {
		return "File download failed.", errors.New("no download URI was provided")
	}

	fileURL, err := url.Parse(downloadUri)
	if err != nil {
		return "File download failed.", fmt.Errorf("parse download URI: %w", err)
	}
	fileName := path.Base(fileURL.Path)
	if fileName == "." || fileName == "/" || fileName == "" {
		return "File download failed.", errors.New("download URI does not include a file name")
	}

	outputDir := util.OutputDir
	if outputDir == "" {
		outputDir, err = os.UserHomeDir()
		if err != nil {
			return "File download failed.", fmt.Errorf("get user's home directory: %w", err)
		}
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "File download failed.", fmt.Errorf("create output directory: %w", err)
	}

	request, err := http.NewRequest(http.MethodGet, downloadUri, nil)
	if err != nil {
		return "File download failed.", err
	}
	request.Header.Set("Authorization", common.SetBearer(util.Token))
	response, err := (&http.Client{}).Do(request)
	if err != nil {
		return "File download failed.", err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "File download failed.", fmt.Errorf("download failed with status %d", response.StatusCode)
	}
	file, err := os.Create(filepath.Join(outputDir, fileName))
	if err != nil {
		return "File download failed.", fmt.Errorf("create destination file: %w", err)
	}
	if _, err := io.Copy(file, response.Body); err != nil {
		file.Close()
		return "File download failed.", fmt.Errorf("write downloaded artifact: %w", err)
	}
	if err := file.Close(); err != nil {
		return "File download failed.", fmt.Errorf("close downloaded artifact: %w", err)
	}
	return "Completed file download", nil
}

func UploadFile(sourcePath, targetPath string) (string, error) {
	if sourcePath == "" || targetPath == "" {
		return "", errors.New("cannot upload without a source file and target path")
	}
	if filepath.Ext(sourcePath) == "" {
		return "", errors.New("source path must include a filename with an extension")
	}

	filePath := filepath.Dir(sourcePath)
	fileName := filepath.Base(sourcePath)
	filesInDirectory, err := os.ReadDir(filePath)
	if err != nil {
		return "", err
	}
	for _, file := range filesInDirectory {
		if strings.EqualFold(fileName, file.Name()) {
			fileName = file.Name()
			break
		}
	}
	actualSourcePath := filepath.Join(filePath, fileName)
	file, err := os.Open(actualSourcePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			common.LogTxtHandler().Warn("File doesn't exist.")
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	baseURL := strings.TrimSuffix(strings.TrimRight(util.ServerApi, "/"), "/api")
	artifactURL := strings.TrimRight(baseURL, "/") + "/" + strings.Trim(targetPath, "/") + "/" + url.PathEscape(fileName)
	request, err := http.NewRequest(http.MethodPut, artifactURL, file)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", common.SetBearer(util.Token))

	common.LogTxtHandler().Debug("REQUEST: Sending 'PUT' request to: " + artifactURL)
	response, err := (&http.Client{}).Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("Artifactory upload failed with status %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}

	var result artifJson
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("decode Artifactory upload response: %w", err)
	}
	if result.DownloadUri == "" {
		return "", errors.New("Artifactory upload response did not include a download URI")
	}
	return result.DownloadUri, nil
}

func DeleteArtifact(artifUri string) (string, error) {
	bearer := common.SetBearer(util.Token)
	var statusCode string
	common.LogTxtHandler().Info(">>> Deleting Artifact: " + artifUri + "...")

	if artifUri != "" {
		common.LogTxtHandler().Debug("REQUEST: Sending 'DELETE' request to: " + artifUri)
		request, err := http.NewRequest("DELETE", artifUri, nil)
		if err != nil {
			return "", err
		}
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
	var statusCode string
	var request *http.Request
	var err error
	common.LogTxtHandler().Info(">>> Getting artifact: " + downloadUri)

	if downloadUri != "" {
		request, err = http.NewRequest("GET", downloadUri, nil)
		if err != nil {
			return "", err
		}
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

func CheckFileAndUpload_OLD(items []os.DirEntry, sourceDir, targetDir, fileName, imageName string) (string, error) {
	// sourceDir ex: c:\\lab\\ or /lab/ - assumes ending slash
	// targetDir ex: /repo-name/folder/ - assumes ending slash
	var sourcePath, targetPath, downloadUri string
	var err error

	for _, item := range items {
		common.LogTxtHandler().Debug("Checking " + fileName + " against " + item.Name() + "...")
		// File comparison is case-sensitive; it will fail if they don't match.
		// So we'll convert them both to lowercase first...
		itemNameLower := strings.ToLower(item.Name())
		fileNameLower := strings.ToLower(fileName)
		imageNameLower := strings.ToLower(imageName) // as the filename will be lowercase, making folder name match case otherwise will be in a separate folder if existing image folder already exists in lowercase

		if itemNameLower == fileNameLower { //item.Name() == fileName {
			common.LogTxtHandler().Debug("File matches...preparing to upload.")
			sourcePath = sourceDir + fileName

			if imageName != "" {
				targetPath = targetDir + imageNameLower + "/"
			} else {
				targetPath = common.CheckAddSlashToPath(targetDir)
			}

			common.LogTxtHandler().Debug("Source Path: " + sourcePath)
			common.LogTxtHandler().Debug("Target Path: " + targetPath)

			downloadUri, err = UploadFile(sourcePath, targetPath)

			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Error uploading: " + fileName + " to: " + targetPath + " - " + strErr)

			} else if downloadUri != "" {
				common.LogTxtHandler().Info("File: " + fileName + " uploaded.")
				common.LogTxtHandler().Info("Download URI: " + downloadUri)
			}
			break
		} else {
			common.LogTxtHandler().Debug("Checking next file...")
		}
	}

	if downloadUri != "" {
		return "Success", nil
	} else if downloadUri == "" && err == nil {
		return "Failed", nil
	} else {
		return "Failed", err
	}
}

func CheckFileAndUpload(sourceDir, targetDir, fileName, imageName string) (string, error) {
	// sourceDir ex: c:\\lab\\ or /lab/ - assumes ending slash
	// targetDir ex: /repo-name/folder/ - assumes ending slash
	var sourcePath, targetPath, downloadUri string
	var err error

	sourcePath = sourceDir + fileName

	if imageName != "" {
		targetPath = targetDir + imageName + "/"
	} else {
		targetPath = common.CheckAddSlashToPath(targetDir)
	}

	common.LogTxtHandler().Debug("Source Path: " + sourcePath)
	common.LogTxtHandler().Debug("Target Path: " + targetPath)

	// 'UploadFile' validates file before upload
	downloadUri, err = UploadFile(sourcePath, targetPath)

	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error uploading: " + fileName + " to: " + targetPath + " - " + strErr)
	} else if downloadUri != "" {
		common.LogTxtHandler().Info("File: " + fileName + " uploaded.")
		common.LogTxtHandler().Info("Download URI: " + downloadUri)
	} else {
		common.LogTxtHandler().Debug("Checking next file...")
	}

	if downloadUri != "" {
		return "Success", nil
	} else if downloadUri == "" && err == nil {
		return "Failed", nil
	} else {
		return "Failed", err
	}
}

func CheckFileAndDownload(checkFile, downloadPath, task string) (string, error) {
	// checkFile - filename with extension
	// downloadPath - parsed Artifactory path to artifact without the artifact file name
	// task - what file check we are performing
	var resultMsg string
	common.LogTxtHandler().Debug("Download Path: " + downloadPath)
	statusCode, err := GetArtifact(downloadPath + checkFile)
	common.LogTxtHandler().Debug("Status code of GetArtifact: " + statusCode)

	if statusCode == "200" {
		// If we found the artifact, download it...
		resultMsg, err = RetrieveArtifact(downloadPath + checkFile)
		if err != nil {
			common.LogTxtHandler().Error(resultMsg)
			return "Failed", err
		}
		common.LogTxtHandler().Info("End of " + task)
		return "Success", nil
	} else {
		common.LogTxtHandler().Error("Artifact not found. Check the server path or artifact name.")
		err := errors.New("Artifact not found. Check the server path or artifact name.")
		return "Failed", err
	}
}

func CheckFileLoopAndDownload(imageName, downloadPath, extString, task string) (string, error) {
	// imageName - name of image we'll use to construct the filename with
	// downloadPath - parsed Artifactory path to artifact without the artifact file name
	// extString - vSphere-based disk file extension - ex: ".vmdk", "-ctk.vmdk", "-flat.vmdk"
	// task - what file check we are performing
	var resultMsg, strI string
	for i := 1; i < 15; i++ { // allowing possibility of up to 15 disk files
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
