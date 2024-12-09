package archive

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/raynaluzier/go-artifactory/common"
	"github.com/raynaluzier/go-artifactory/util"
)

func GetArtifactsByNameRepo(artifName, repo string) ([]string, error) {
	// Searches for artifacts by artifact name (can be partial) and optionally Artifactory repo (can be partial)
	listArtifUris := []string{}
	bearer := common.SetBearer(util.Token)
	requestPath := util.ServerApi + "/search/artifact?name=" + artifName

	if artifName != "" && repo != "" {
		// Determines whether a repo was also supplied and if so, includes it in the API call
		if repo != "" {
			request, err = http.NewRequest("GET", requestPath + "&repos=" + repo, nil)
		} else {
			request, err = http.NewRequest("GET", requestPath, nil)
		}
	
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		//fmt.Println(string(body))
		
		// JSON return is results with an array of one or more URI strings
		type resultsJson struct {
			Results []struct{
				Uri string `json:"uri"`
			} `json:"results"`
		}
		
		// Unmarshal the JSON return
		var jsonData *resultsJson
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			fmt.Printf("Could not unmarshal %s\n", err)
		}
		
		// As long as the results are not empty, parse thru the results and append the URI for each 
		// matching artifact to a list of strings
		if len(jsonData.Results) != 0 {
			for idx, r := range jsonData.Results {
				r = jsonData.Results[idx]
				listArtifUris = append(listArtifUris, r.Uri)
			}
			return listArtifUris, nil
		} else {
			err := errors.New("No results returned")
			return nil, err
		}
	} else {
		// If at least a partial artifact name AND repo aren't supplied, we'll throw an error
		message := ("Supplied Artifact name is: " + artifName + ", Repo is: " + repo)
		fmt.Println(message)
		err := errors.New("Unable to search for Artifact without at least a partial Artifact name and parent Repository name")
		return nil, err
	}
}

func GetArtifactVersions(groupId, artifName, repo string) ([]string, error) {
	// Requires at least the Group ID (top level folder, must be FULL name) and Artifact Name (must be FULL name); optionally repo
	// Only available if folder structure was setup with a Layout (artifacts will have a value for Module ID present)
    // Search terms are CASE SENSITIVE
	listVersions := []string{}
	bearer := common.SetBearer(util.Token)
	requestPath := util.ServerApi + "/search/versions?g=" + groupId + "&a=" + artifName

	// Ensures the group ID and artifact name fields are not empty; these are required for this search type
	if (groupId != "") && (artifName != "") {
		if repo != "" {
			// If a repo was supplied, it will be added to the API call
			request, err = http.NewRequest("GET", requestPath + "&repos" + repo, nil)
		} else {
			request, err = http.NewRequest("GET", requestPath, nil)
		}
			
		request.Header.Add("Authorization", bearer)
		
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		//fmt.Println(string(body))

		// JSON return is results with an array of one or more available version strings 
		type resultsJson struct {
			Results []struct{
				Version string `json:"version"`
			} `json:"results"`
		}
		
		// Unmarshal the JSON return
		var jsonData *resultsJson
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			fmt.Printf("Could not unmarshal %s\n", err)
		}
		
		// As long as the results are not empty, parse thru the results and append the available versions for 
		// the given artifact to a list of strings
		if len(jsonData.Results) != 0 {
			for idx, r := range jsonData.Results {
				r = jsonData.Results[idx]
				listVersions = append(listVersions, r.Version)
			}
			return listVersions, nil
		} else {
			err := errors.New("No version results returned")
			return nil, err
		}

	} else {
		// If the group ID AND artifact name are not supplied, we'll throw an error
		message := ("Supplied group ID is: " + groupId + " and artifact name is: " + artifName)
		fmt.Println(message)
		err := errors.New("Group ID and Artifact Name values can't be empty")
		return nil, err
	}
}

func GetArtifactLatestVersion(groupId, artifName, repo string) (string, error) {
	// Requires at least the Group ID (top level folder, must be FULL name) and Artifact Name (must be FULL name); optionally repo
	// Only available if folder structure was setup with a Layout (artifacts will have a value for Module ID present)
	// Search is CASE SENSITIVE
	var latestVersion string
	bearer := common.SetBearer(util.Token)
	requestPath := util.ServerApi + "/search/latestVersion?g=" + groupId + "&a=" + artifName

	// Ensures the group ID and artifact name fields are not empty; these are required for this search type
	if (groupId != "") && (artifName != "") {
		if repo != "" {
			// If a repo was supplied, it will be added to the API call
			request, err = http.NewRequest("GET", requestPath + "&repos" + repo, nil)
		} else {
			request, err = http.NewRequest("GET", requestPath, nil)
		}
			
		request.Header.Add("Authorization", bearer)
		
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		//fmt.Println(string(body))        

		// If the status is OK, the latest version (which is simply a string) will be returned
		if response.StatusCode == http.StatusOK {
			latestVersion = string(body)
		} else {
			err := errors.New("No version results returned")
			return "", err
		}

	} else {
		// If the group ID AND artifact name are not supplied, we'll throw an error
		message := ("Supplied group ID is: " + groupId + " and artifact name is: " + artifName)
		fmt.Println(message)
		err := errors.New("Group ID and Artifact Name values can't be empty")
		return "", err
	}

	if err != nil {
		fmt.Println("Unable to parse URL")
		return "", err
	}

	return latestVersion, nil
}