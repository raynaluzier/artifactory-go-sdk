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

func GetItemChildren(item string) ([]Contents, error) {
	// Item can represent a repo name or a combo of repo/child_folder/subchild_folder/etc
	artifBase, bearer := common.AuthCreds()
	requestPath := artifBase + "/storage/" + item

	type itemResults struct {
		Repo			string		`json:"repo"`
		Path			string		`json:"path"`
		Created			string		`json:"created"`
		LastModified 	string		`json:"lastModified"`
		LastUpdated		string		`json:"lastUpdated"`
		Children []struct {
			Uri		string		`json:"uri"`
			Folder	bool		`json:"folder"`
		}
		Uri				string		`json:"uri"`
		CreatedBy		string		`json:"createdBy"`	// Not exposed at repo level
		ModifiedBy		string		`json:"modifiedBy"`	// Not exposed at repo level
	}

	var childDetails []Contents

	if (item != "") {
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

		var jsonData *itemResults
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			fmt.Printf("Could not unmarshal %s\n", err)
		}

		// If the item has children, parse the data and return the abbreviated
		// URI ('/folder', '/folder/artifact.ext', etc) and whether the child item is a folder or not (bool)
		if len(jsonData.Children) != 0 {
			for idx, c := range jsonData.Children {
				c = jsonData.Children[idx]
				childDetails = append(childDetails, Contents{Child: c.Uri, IsFolder: c.Folder})
			}
			return childDetails, nil
		} else {
			// If no children found, we return empty contents; this isn't an error condition
			fmt.Println("No child objects found for " + item)
			return childDetails, nil
		}

		/*
		for idx := 0; idx < len(childDetails); idx++ {
			fmt.Println(childDetails[idx].Child, childDetails[idx].IsFolder)
		}*/
		// ex: [{/test-artifact-1.1.txt false} {/test-artifact-1.2.txt false} {/test-artifact-1.3.txt false}] <nil>
	} else {
		err := errors.New("No item or path provided. Unable to get child items without parent item/path.")
		return nil, err
	}
}


func GetArtifactPath(artifName string) ([]string, error) {
	// Takes in an artifact's name and searches Artifactory, returning the path to the artifact
	// Searches are CASE SENSITIVE
	var childList []Contents
	var listOfPaths []string
	foundPaths = nil

	if artifName != "" {
		listRepos, err := ListRepos()
		if err != nil || listRepos[0] == "" || len(listRepos) == 0 {
			err := errors.New("No repos found")
			return nil, err
		}

		if len(listRepos) != 0 {
			for idx := 0; idx < len(listRepos); idx++ {
				childList, err = GetItemChildren(listRepos[idx])
				if len(childList) != 0 {
					listOfPaths = RecursiveSearch(childList, artifName, listRepos[idx], foundPaths)
				}
			}

			if len(listOfPaths) > 1 {
				// We'll search the list for duplicates and remove them
				listOfPaths = common.RemoveDuplicateStrings(listOfPaths)
				if len(listOfPaths) > 1 {
					fmt.Println("More than one possible artifact path found")
				}
				return listOfPaths, nil
			} else if len(listOfPaths) == 1 && listOfPaths[0] != "" {
				return listOfPaths, nil
			} else if len(listOfPaths) == 0 || listOfPaths[0] == "" {
				err := errors.New("Unable to find path to artifact")
				return nil, err
			}
		} else {
			err := errors.New("List of repos to check is empty. Either there are no repos or you do not have sufficient permissions to the repo(s).")
			return nil, err
		}
	} else {
		err := errors.New("Unable to determine path to artifact without the artifact name")
		return nil, err
	}

	// Now we have a list of repos to check through... can do another search for props... ** FINISH THIS
	return listOfPaths, err
}

func RecursiveSearch(list []Contents, artifName, searchPath string, foundPaths []string) ([]string) {
	// Recursively searches a list of child items for the specificied artifact name 
	var nextList []Contents
	var currentPath string
	currentPath = searchPath

	if len(list) != 0 {
		for item := 0; item < len(list); item++ {					// For each item in list...
			if list[item].IsFolder == false {						// If not a folder, does artifact match?
				// 'Contains' search is case sensitive; so we'll convert the input artifact name and convert to both cases and recheck
				lowStr := common.ConvertToLowercase(artifName)
				upStr := common.ConvertToUppercase(artifName)

				if strings.Contains(list[item].Child, artifName) {      // If we don't find it initially, we'll check with cases converted
					foundPaths = append(foundPaths, searchPath)         // If found, item's path appended to found list
				} else if strings.Contains(list[item].Child, lowStr) {
					foundPaths = append(foundPaths, searchPath)
				} else if strings.Contains(list[item].Child, upStr) {
					foundPaths = append(foundPaths, searchPath)
				}
			} else {  // IsFolder == true; so we get its children and repeat the search
				searchPath = currentPath + list[item].Child		   // 1st "/repo" + "/folder", 2nd "/repo/folder" + "/folder", etc
				nextList, err = GetItemChildren(searchPath)
				if len(nextList) != 0 {
					foundPaths = RecursiveSearch(nextList, artifName, searchPath, foundPaths)
				}
			}
		}
	}
	return foundPaths
}

func GetDownloadUri(artifPath, artifNameExt string) (string, error) {
	// Requires full path to the artifact, include full artifact name with extention
	artifBase, bearer := common.AuthCreds()
	requestPath := artifBase + "/storage" + artifPath + "/" + artifNameExt
	var downloadUri string

	if (artifPath != "") && (artifNameExt != "") {
		request, err = http.NewRequest("GET", requestPath, nil)
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
			err = errors.New("There is no download URI for the artifact")
			return "", err
		}
	} else {
		message := ("Supplied artifact path: " + artifPath + " and full artifact name: " + artifNameExt)
		fmt.Println(message)
		err := errors.New("Unable to get artifact details without full path to the artifact")
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
		message := ("Supplied artifact path: " + artifactUri)
		fmt.Println(message)
		err := errors.New("Unable to get artifact details without full path to the artifact")
		return "", err
	}
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
