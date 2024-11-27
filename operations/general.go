package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/raynaluzier/go-artifactory/common"
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
	artifBase, bearer := common.AuthCreds()
	requestPath := artifBase + "/repositories"

	request, err = http.NewRequest("GET", requestPath, nil)
	request.Header.Add("Authorization", bearer)
	
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error on response.\n[ERROR] - ", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	//fmt.Println(string(body))

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
		fmt.Printf("Could not unmarshal %s\n", err)
	}

	if len(jsonData) != 0 {
		for _, k := range jsonData {
			listRepos = append(listRepos, k.Key)
		}
		return listRepos, nil
	} else {
		err := errors.New("No repos found")
		return nil, err
	}
}

func GetDownloadUri(artifUri string) (string, error) {
	_, bearer := common.AuthCreds()
	var downloadUri string

	if (artifUri != "") {
		request, err = http.NewRequest("GET", artifUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)

		var jsonData *artifJson
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			fmt.Printf("Could not unmarshal %s\n", err)
		}

		if jsonData.DownloadUri != "" {
			downloadUri = jsonData.DownloadUri
			return downloadUri, nil
		} else {
			err = errors.New("There is no download URI for the artifact.")
			return "", err
		}
	} else {
		err := errors.New("Unable to get artifact's download URI without the artifact's URI.")
		return "", err
	}
}

func GetCreateDate(artifactUri string) (string, error) {
	_, bearer := common.AuthCreds()
	var createdDate string

	if (artifactUri != "") {
		request, err = http.NewRequest("GET", artifactUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)

		var jsonData *artifJson
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			fmt.Printf("Could not unmarshal %s\n", err)
		}

		if jsonData.Created != "" {
			createdDate = jsonData.Created
			return createdDate, nil
		} else {
			err = errors.New("There is no create date for the artifact")
			return "", err
		}
	} else {
		err := errors.New("Unable to get artifact's created date without the artifact's URI.")
		return "", err
	}
}

func GetArtifactNameFromUri(artifUri string) (string) {
	fileName := path.Base(artifUri)
	ext := filepath.Ext(fileName)
	artifactName := strings.TrimSuffix(fileName, ext)
	return artifactName
}

func RetrieveArtifact(downloadUri string) (string, error) {
	// Gets the artifact via provided Download URI and copies it to the output directory specified in
	// the environment variables file
	var outputDir string
	_, bearer := common.AuthCreds()

	// If no output directory path was provided, the artifact file will be downloaded to the top-level
	// directory of this code
	if len(os.Getenv("OUTPUTDIR")) != 0 {
		OUTPUTDIR := os.Getenv("OUTPUTDIR")
		outputDir = common.EscapeSpecialChars(OUTPUTDIR)  // Ensure special characters are escaped
		outputDir = common.CheckAddSlashToPath(outputDir) // Ensure path ends with appropriate slash type
	} else {  // There's no OUTPUTDIR env var...
		fmt.Println("No output directory provided; output will be at top-level directory")
		outputDir = ""
	}

	if downloadUri != "" {
		request, err = http.NewRequest("GET", downloadUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		//fmt.Println(string(body))   // prints the contents of the file

		if response.StatusCode == 404 {
			err := errors.New("File not found.")
			return "File download failed.", err
		} else {
			// Create file name from download URI path of artifact
			fileUrl, err := url.Parse(downloadUri)
			if err != nil {
				log.Fatal(err)
			}

			// Get the file name from the path
			path := fileUrl.Path
			segments := strings.Split(path, "/")
			fileName := segments[len(segments)-1]

			// Creates the file at the defined path
			// Will overwrite the file if it already exists
			newFile, err := os.Create(outputDir + fileName)   
			if err != nil {
				log.Fatal(err)
				return "Error creating file at target location.", err
			}
			err = os.WriteFile(outputDir + fileName, body, 0644)
			if err != nil {
				log.Fatal(err)
				return "Error downloading file to target location.", err
			}
			defer newFile.Close()
		}
	} else {
		err := errors.New("No download URI was provided. Unable to download the artifact without the download URI.")
		return "File download failed.", err
	}

	return "Completed file download", nil
}

func UploadFile(sourcePath, targetPath, fileSuffix string) (string, error) {
	//** TO DO: Option to get previous 'version' and increment

	var downloadUri string
	var filePath string
	var fileName string
	var found bool
	artifBase, bearer := common.AuthCreds()
	separater := "-"										  // If adding a file suffix (like date, version, etc), use this separater between filename and suffix
	trimmedBase := artifBase[:len(artifBase)-4]               // Removing '/api' from base URI

	if len(sourcePath) != 0 && targetPath != "" { 
		// We need to ensure the provided source path/file are valid and exist
		if len(path.Ext(sourcePath)) != 0 {		                  // Ensures file with extension exists in source path
			sourcePath = common.EscapeSpecialChars(sourcePath)
			targetPath = common.EscapeSpecialChars(targetPath)
			targetPath = common.CheckAddSlashToPath(targetPath)
			
			// Determine source filename and source file path by platform type
			winPath := common.CheckPathType(sourcePath)
			if winPath == true {
				segments := strings.Split(sourcePath, "\\")	  	  // Split source path into segments
				fileName = segments[len(segments)-1]			  // Determine filename from path
				filePath = sourcePath[:len(sourcePath)-len(fileName)]  // Determine path without filename
			} else {   // Unix path
				segments := strings.Split(sourcePath, "/")	 	  // Split source path into segments
				fileName = segments[len(segments)-1]              // Determine filename from path
				filePath = sourcePath[:len(sourcePath)-len(fileName)]  // Determine path without filename				
			}
			
			// Get all files in the provided source directory
			filesInDirectory, err := os.ReadDir(filePath)
			if err != nil {
				return "", err
			}
			
			// For each file in the source directory, do a case insensitive file name comparison for a match
			// As Artifactory cares about case here, we want to make sure the filename supplied matches the case of the filename that actually exists in the source path
			for _, file := range filesInDirectory {
				isSameStr := common.StringCompare(fileName, file.Name())            // Filename from provided source path vs. filename pulled directly from source path
				if isSameStr == true {												// If true, we know files are the same
					found = true													// Mark that we found a matching file
					isExactStr, err := common.SearchForExactString(file.Name(), fileName)  // Now, checks if cases matches
					if err != nil {
						fmt.Println("Error searching for exact string: ", err)
					}
					
					if isExactStr == false {										// Files are the same, but provided and actual cases are different
						fileName = file.Name() 										// Set the provided filename to match the actual filename so we'll use to the correct case
					}
				}
			}

			// If we couldn't find a matching file at all, then we throw an error
			if found == false {
				err := errors.New("Unable to validate existance of source file. Source file doesn't exist.")
				return "", err
			}
			
			// We now have a validated source path and filename
			// Set target filename = filename + fileSuffix (if not blank)
			if len(fileSuffix) != 0 || fileSuffix != "" {							// If a file suffix (like version, date, etc) was provided...
				fileExt := path.Ext(fileName)										// Returns .[ext]
				justName := strings.Trim(fileName, fileExt)							// Trim off extension
				fileName = justName + separater + fileSuffix + fileExt
			}   // If blank, then the original filename will be used
			
			newArtifactPath := trimmedBase + targetPath + fileName                  // Forms: http://artifactory_base_api_url/repo-key/folder/artifact.txt
			data := strings.NewReader("@/" + sourcePath)                            // Formats the payload appropriately
			
			request, err = http.NewRequest("PUT", newArtifactPath, data)
			request.Header.Add("Authorization", bearer)
	
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				log.Println("Error on response.\n[ERROR] - ", err)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			fmt.Println(string(body))
	
			var jsonData *artifJson
			err = json.Unmarshal(body, &jsonData)
			if err != nil {
				fmt.Printf("Could not unmarshal %s\n", err)
			}
	
			if jsonData.DownloadUri != "" {
				downloadUri = jsonData.DownloadUri
				return downloadUri, nil
			} else {
				err = errors.New("There is no download URI for the artifact")
				return "", err
			}
		} else {
			err = errors.New("No file extension found in source path. Ensure source includes path and source file with extension.")
			return "", err
		}
	} else {
		message := ("Supplied source path: " + sourcePath + ", target path: " + targetPath)
		err := errors.New("Cannot upload file without source path/file, target path, and artifact file name")
		fmt.Println(message)
		return "", err
	}
}

func DeleteArtifact(artifUri string) (string, error) {
	_, bearer := common.AuthCreds()

	if artifUri != "" { 
		request, err = http.NewRequest("DELETE", artifUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		fmt.Println(string(body))
		
		if response.StatusCode == 204 {
			fmt.Println("Request completed successfully")
			statusCode = "204"
		} else {
			fmt.Println("Unable to complete request")
			statusCode = "404"
		}
	} else {
		message := ("Supplied artifact path is: " + artifUri)
		fmt.Println(message)
		err := errors.New("Unable to DELETE item without artifact URI.")
		return "", err
	}
	return statusCode, nil
}