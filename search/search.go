package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/raynaluzier/go-artifactory/common"
	"github.com/raynaluzier/go-artifactory/util"
)

var request *http.Request
var err error

func GetArtifactsByProps(listKvProps []string) ([]string, error) {
	// Takes in list of property key/values strings (ex: 'release=latest-stable', 'testing=passed')
	var strKvProps string
	listArtifUris := []string{}
	bearer := common.SetBearer(util.Token)
	requestPath := util.ServerApi + "/search/prop?"

	common.LogTxtHandler().Info(">>> Getting Artifacts by Property Names/Values...")

	if len(listKvProps) != 0 {
		if len(listKvProps) > 1 {
			// If there's more than one prop name/value supplied, adds the required '&' separater between them
			strKvProps = strings.Join(listKvProps, "&")
			request, err = http.NewRequest("GET", requestPath + strKvProps, nil)
			common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + requestPath + strKvProps)

		} else {
			request, err = http.NewRequest("GET", requestPath + listKvProps[0], nil)
			common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + requestPath + listKvProps[0])
		}

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
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Could not unmarshal response - " + strErr)
			}

			// As long as the results are not empty, parse thru the results and append the URI for each 
			// matching artifact to a list of strings
			if len(jsonData.Results) != 0 {
				for idx, r := range jsonData.Results {
					r = jsonData.Results[idx]
					listArtifUris = append(listArtifUris, r.Uri)
					common.LogTxtHandler().Info("FOUND ARTIFACT: " + r.Uri)
				}
				return listArtifUris, nil
			} else {
				err := errors.New("No artifacts returned.")
				common.LogTxtHandler().Warn("No artifacts returned.")
				return nil, err
			}
		}
	} else {
		// If no properties were supplied, we'll throw an error
		err := errors.New("Unable to search by Property without at least one Property Name and, optionally, Value")
		common.LogTxtHandler().Error("Supplied Property Name(s)/Value(s): " + strKvProps)
		common.LogTxtHandler().Error("Unable to search by Property without at least one Property Name and, optionally, Value")
		return nil, err
	}

	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to parse URL - " + strErr)
		return nil, err
	}
	return listArtifUris, nil
}

func GetArtifactsByName(artifName string) ([]string, error) {
	// Searches for artifacts by artifact name (can be partial)
	listArtifUris := []string{}
	bearer := common.SetBearer(util.Token)
	requestPath := util.ServerApi + "/search/artifact?name=" + artifName

	common.LogTxtHandler().Info(">>> Getting Artifacts by Name...")

	if artifName != "" {
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
				strErr := fmt.Sprintf("%v\n", err)
				common.LogTxtHandler().Error("Could not unmarshal response - " + strErr)
			}
			
			// As long as the results are not empty, parse thru the results and append the URI for each 
			// matching artifact to a list of strings
			if len(jsonData.Results) != 0 {
				for idx, r := range jsonData.Results {
					r = jsonData.Results[idx]
					listArtifUris = append(listArtifUris, r.Uri)
					common.LogTxtHandler().Info("FOUND ARTIFACT: " + r.Uri)
				}
				return listArtifUris, nil
			} else {
				err := errors.New("No results returned")
				common.LogTxtHandler().Warn("No results returned")
				return nil, err
			}
		}
	} else {
		// If at least a partial artifact name isn't supplied, we'll throw an error
		err := errors.New("Unable to search for Artifact without at least a partial Artifact name.")
		common.LogTxtHandler().Error("Supplied Artifact name is: " + artifName)
		common.LogTxtHandler().Error("Unable to search for Artifact without at least a partial Artifact name.")
		return nil, err
	}
}

func FilterListByFileType(ext string, listArtifacts []string) ([]string, error) {
	// Filters list of artifact URIs by file type
	// If no extension is provided, the default filter will be VMware Templates (.vmxt)
	var filteredList []string

	common.LogTxtHandler().Info(">>> Filtering Artifact URIs by File Extension...")
	common.LogTxtHandler().Info(">>>---> " + ext)

	if ext == "" {
		common.LogTxtHandler().Warn("*** No file type was specified. Using DEFAULT file type of '.vmxt' (VM Template).")
		common.LogTxtHandler().Warn("*** To change this, include a desired file type.")
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
				common.LogTxtHandler().Debug("FOUND MATCHING ARTIFACT WITH EXTENSTION " + ext + ": " + item)
			}
		}
	} else {
		err = errors.New("List of artifacts cannot be empty.")
		common.LogTxtHandler().Error("List of artifacts cannot be empty.")
		return nil, err
	}
	return filteredList, err
}