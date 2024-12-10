package archive

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/raynaluzier/artifactory-go-sdk/common"
	"github.com/raynaluzier/artifactory-go-sdk/operations"
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

func GetItemChildren(item string) ([]Contents, error) {
	// Item can represent a repo name or a combo of repo/child_folder/subchild_folder/etc
	bearer := common.SetBearer(util.Token)
	requestPath := util.ServerApi + "/storage/" + item
	common.LogTxtHandler().Info(">>> Getting Item Children for Item" + item + "...")

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
		common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + requestPath)
		request, err = http.NewRequest("GET", requestPath, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error on response. " + strErr)
			return nil, err

		} else {
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			common.LogTxtHandler().Debug("REQUEST RESPONSE: " + string(body))

			var jsonData *itemResults
			err = json.Unmarshal(body, &jsonData)
			if err != nil {
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Could not unmarshal response - " + strErr)
			}

			// If the item has children, parse the data and return the abbreviated
			// URI ('/folder', '/folder/artifact.ext', etc) and whether the child item is a folder or not (bool)
			if len(jsonData.Children) != 0 {
				common.LogTxtHandler().Debug("CHILD OBJECTS FOUND FOR: " + item)
				for idx, c := range jsonData.Children {
					c = jsonData.Children[idx]
					childDetails = append(childDetails, Contents{Child: c.Uri, IsFolder: c.Folder})

					strIsFolder := fmt.Sprintf("%v\n", c.Folder)
					common.LogTxtHandler().Debug("CHILD: " + c.Uri + " - IS FOLDER: " + strIsFolder)
				}
				return childDetails, nil

			} else {
				// If no children found, we return empty contents; this isn't an error condition
				common.LogTxtHandler().Warn("No child objects found for " + item)
				return childDetails, nil
			}

			/*
			for idx := 0; idx < len(childDetails); idx++ {
				fmt.Println(childDetails[idx].Child, childDetails[idx].IsFolder)
			}*/
			// ex: [{/test-artifact-1.1.txt false} {/test-artifact-1.2.txt false} {/test-artifact-1.3.txt false}] <nil>
			}
	} else {
		err := errors.New("No item or path provided. Unable to get child items without parent item/path.")
		common.LogTxtHandler().Error("No item or path provided. Unable to get child items without parent item/path.")
		return nil, err
	}
}

func GetArtifactPath(artifName string) ([]string, error) {
	// Takes in an artifact's name and searches Artifactory, returning the path to the artifact
	// Searches are CASE SENSITIVE; so RecursiveSearch does upper and lowercase checks as well
	var childList []Contents
	var listOfPaths []string
	foundPaths = nil
	common.LogTxtHandler().Info(">>> Getting Artifact Path for Artifact " + artifName + "...")

	if artifName != "" {
		listRepos, err := operations.ListRepos()
		if err != nil || listRepos[0] == "" || len(listRepos) == 0 {
			err := errors.New("No repos found.")
			common.LogTxtHandler().Warn("No repos found.")
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
				common.LogTxtHandler().Info("Removing duplicate paths...")
				listOfPaths = common.RemoveDuplicateStrings(listOfPaths)
				if len(listOfPaths) > 1 {
					common.LogTxtHandler().Info("More than one possible artifact path found.")
				}
				return listOfPaths, nil

			} else if len(listOfPaths) == 1 && listOfPaths[0] != "" {
				return listOfPaths, nil

			} else if len(listOfPaths) == 0 || listOfPaths[0] == "" {
				err := errors.New("Unable to find path to artifact.")
				common.LogTxtHandler().Error("Unable to find path to artifact.")
				return nil, err
			}
		} else {
			err := errors.New("List of repos to check is empty. Either there are no repos or you do not have sufficient permissions to the repo(s).")
			common.LogTxtHandler().Error("List of repos to check is empty. Either there are no repos or you do not have sufficient permissions to the repo(s).")
			return nil, err
		}
	} else {
		err := errors.New("Unable to determine path to artifact without the artifact name.")
		common.LogTxtHandler().Error("Unable to determine path to artifact without the artifact name.")
		return nil, err
	}

	// Now we have a list of paths to check through... can do another search for props... ** FINISH THIS
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
					common.LogTxtHandler().Debug("Found possible path: " + searchPath)

				} else if strings.Contains(list[item].Child, lowStr) {
					foundPaths = append(foundPaths, searchPath)
					common.LogTxtHandler().Debug("Found possible path: " + searchPath)

				} else if strings.Contains(list[item].Child, upStr) {
					foundPaths = append(foundPaths, searchPath)
					common.LogTxtHandler().Debug("Found possible path: " + searchPath)
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