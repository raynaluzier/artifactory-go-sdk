package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/raynaluzier/go-artifactory/common"
)

var request *http.Request
var err error

func GetArtifactsByProps(listKvProps []string) ([]string, error) {
	// Searches for an artifact by one or more property names and optionally values if provided (e.g. release or release=stable)
	var strKvProps string
	listArtifUris := []string{}
	artifBase, bearer := common.AuthCreds()
	requestPath := artifBase + "/search/prop?"

	if len(listKvProps) != 0 {
		if len(listKvProps) > 1{
			// If there's more than one prop name/value supplied, adds the required '&' separater between them
			strKvProps = strings.Join(listKvProps, "&")
			request, err = http.NewRequest("GET", requestPath + strKvProps, nil)
		} else {
			request, err = http.NewRequest("GET", requestPath + listKvProps[0], nil)
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
			err := errors.New("No artifacts returned")
			return nil, err
		}

	} else {
		message := ("Supplied Property Name(s)/Value(s): " + strKvProps)
		fmt.Println(message)
		err := errors.New("Unable to search by Property without at least one Property Name and, optionally, Value")
		return nil, err
	}
}

func GetArtifactsByName(artifName string) ([]string, error) {
	// Searches for artifacts by artifact name (can be partial); search is CASE INSENSITIVE
	listArtifUris := []string{}
	artifBase, bearer := common.AuthCreds()
	requestPath := artifBase + "/search/artifact?name=" + artifName

	if artifName != "" {
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
		
		// JSON return is results with an array of one or more URI strings
		type resultsJson struct {
			Results []struct{
				Uri string `json:"uri"`
			} `json:"results"`
		}
		
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
		message := ("Supplied Artifact name is: " + artifName)
		fmt.Println(message)
		err := errors.New("Unable to search for Artifact without at least a partial Artifact name")
		return nil, err
	}
}

func FilterListByFileType(ext string, listArtifacts []string) ([]string, error) {
	// If no extension is provided, the default filter will be VMware Templates (.vmxt)
	var filteredList []string

	if ext == "" {
		ext = ".vmxt"
	}

	if len(listArtifacts) != 0 {
		if strings.Contains(ext, ".") {			// If the file extension already contains '.', don't do anything
		} else {
			ext = "." + ext						// Otherwise, add leading '.'
		}

		for _, item := range listArtifacts {
			if path.Ext(item) == ext {
				filteredList = append(filteredList, item)
			}
		}
	} else {
		err = errors.New("List of artifacts cannot be empty.")
		return nil, err
	}
	return filteredList, err
}
