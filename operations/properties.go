package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/raynaluzier/go-artifactory/common"
)

var statusCode string

type prop struct {
	Name string
	Value string
}

func GetArtifactPropVals(artifUri string, listPropKeys []string) (interface{}, error){
	// Returns the values for only the properties included in the URI for the given artifact
	// Search is CASE SENSTIVE
	_, bearer := common.AuthCreds()
	var properties []prop

	common.LogTxtHandler().Info(">>> Getting Values for Specified Artifact Property(ies): " + artifUri)
	
	if artifUri != "" {
		if len(listPropKeys) > 1 {
			// If there's more than one property name supplied, adds the required ',' separater between them
			strProps := strings.Join(listPropKeys, ",")
			request, err = http.NewRequest("GET", artifUri + "?properties=" + strProps, nil)
			common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + artifUri + "?properties=" + strProps)
		} else if len(listPropKeys) == 1 && listPropKeys[0] != "" {
			request, err = http.NewRequest("GET", artifUri + "?properties=" + listPropKeys[0], nil)
			common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + artifUri + "?properties=" + listPropKeys[0])
		} else {
			err := errors.New("Unable to search for Artifact properties without one or more property names")
			common.LogTxtHandler().Error("Unable to search for Artifact properties without one or more property names")
			return nil, err
		}
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

			if err != nil || response.StatusCode == 404 {
				err := errors.New("No matching property(ies) could be found.")
				common.LogTxtHandler().Error("No matching property(ies) could be found.")
				return nil, err
			} else {
				// Declares a map whose key type is a string with any value type
				// This is used because the returned JSON data is unstructured; 'properties' contains one or more key/values that
				// correspond to a property name and property value that can be anything
				var result map[string]any

				// Unmarshal the JSON return
				err = json.Unmarshal(body, &result)
				if err != nil {
					strErr := fmt.Sprintf("%v\n", err)
					common.LogTxtHandler().Error("Could not unmarshal response - " + strErr)
				}

				// The property keys are returned as a string, but the values must be converted to string first
				// and the surrounding [ ] brackets are trimmed off
				// Each key/value pair are stored in a struct of type 'prop' and returned, allowing for easier parsing later
				var strValue string

				parseProps := result["properties"].(map[string]any)
				if len(parseProps) != 0 {
					for k, v := range parseProps {
						strValue = fmt.Sprintf("%v", v)
						strValue = strings.Trim(strValue, "]")
						strValue = strings.Trim(strValue, "[")
						properties = append(properties, prop{Name: k, Value: strValue})
						common.LogTxtHandler().Debug("FOUND PROPERTY: " + k + " with VALUE: " + strValue)
					}
					return properties, nil
					/*for idx := 0; idx < len(properties); idx++ {
						fmt.Println(properties[idx].Name, properties[idx].Value)
					}*/
				} else {
					err := errors.New("No results returned.")
					common.LogTxtHandler().Warn("No results returned.")
					return nil, err
				}
			}
		}
	} else {
		if len(listPropKeys) != 0 && listPropKeys[0] != "" {
			err := errors.New("Unable to search for Artifact properties without the artifact's URI.")
			common.LogTxtHandler().Error("Unable to search for Artifact properties without the artifact's URI.")
			return nil, err
		} else {
			err := errors.New("Unable to search for Artifact properties without the artifact's URI and one or more property names.")
			common.LogTxtHandler().Error("Unable to search for Artifact properties without the artifact's URI and one or more property names.")
			return nil, err
		}
	}
}

func GetAllPropsForArtifact(artifUri string) (interface{}, error) {
	_, bearer := common.AuthCreds()
	var properties [] prop

	common.LogTxtHandler().Info(">>> Getting All Properties for Artifact: " + artifUri + "...")

	if artifUri != "" {
		common.LogTxtHandler().Debug("REQUEST: Sending 'GET' request to: " + artifUri)
		request, err = http.NewRequest("GET", artifUri + "?properties", nil)
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

			if err != nil || response.StatusCode == 404 {
				err := errors.New("No property(ies) found.")
				common.LogTxtHandler().Debug("No property(ies) found.")
				return nil, err
			} else {
				// Declares a map whose key type is a string with any value type
				// This is used because the returned JSON data is unstructured; 'properties' contains one or more key/values that
				// correspond to a property name and property value that can be anything
				var result map[string]any

				// Unmarshal the JSON return
				err = json.Unmarshal(body, &result)
				if err != nil {
					strErr := fmt.Sprintf("%v\n", err)
					common.LogTxtHandler().Error("Could not unmarshal response " + strErr)
				}

				// The property keys are returned as a string, but the values must be converted to string as well,
				// and the surrounding [ ] brackets are trimmed off
				// Each key/value pair are stored in a struct of type 'prop' and returned, allowing for easier parsing later
				var strValue string

				parseProps := result["properties"].(map[string]any)
				if len(parseProps) != 0 {
					for k, v := range parseProps {
						strValue = fmt.Sprintf("%v", v)
						strValue = strings.Trim(strValue, "]")
						strValue = strings.Trim(strValue, "[")
						properties = append(properties, prop{Name: k, Value: strValue})
						common.LogTxtHandler().Debug("FOUND PROPERTY: " + k + " with VALUE: " + strValue)
					}
					return properties, nil
					
				} else {
					err := errors.New("No results returned.")
					common.LogTxtHandler().Warn("No results returned.")
					return nil, err
				}
			}
		}
	} else {
		err := errors.New("Unable to retrieve properties of the artifact without the Artifact's URI.")
		common.LogTxtHandler().Error("Unable to retrieve properties of the artifact without the Artifact's URI.")
		return nil, err
	}
}

func FilterListByProps(listArtifUris, listKvProps []string) (string, error) {
	var foundList []string
	var filteredList []string
	var structData []map[string]interface{}
	numProps := len(listKvProps)
	var foundItem string

	common.LogTxtHandler().Info(">>> Filtering Artifact URIs by Property Keys/Values...")
	for p := 0; p < len(listKvProps); p++ {
		common.LogTxtHandler().Info(">>>---> " + listKvProps[p])
	}

	if len(listArtifUris) != 0 && len(listKvProps) != 0 {
		for a := 0; a < len(listArtifUris); a++ {
			// For each artifact URI in list, get it's properties/values; there can be one or more properties/values assigned
			artifProps, err := GetAllPropsForArtifact(listArtifUris[a])  // ex return: [{release stable} {testing passed}]
			if err != nil {
				common.LogTxtHandler().Debug("No properties returned for artifact: " + listArtifUris[a])
			} else {
				// Convert custom data type 'prop' object passed out as interface{} into JSON format
				jsonBytes, err := json.Marshal(artifProps)
				if err != nil {
					strErr := fmt.Sprintf("%v\n", err)
					common.LogTxtHandler().Error("Error on response. " + strErr)
				}
				// Convert the JSON data into a map of arbitrary values to support any type (in this case, our custom 'prop' type)
				err = json.Unmarshal([]byte(jsonBytes), &structData)
				if err != nil {
					strErr := fmt.Sprintf("%v\n", err)
					common.LogTxtHandler().Error("Could not unmarshal response - " + strErr)
				}

				// For each returned key/value property assigned to the artifact...
				for idx := 0; idx < len(structData); idx++ {
					// Convert each pair to a string and format to match the listKvProps input ('key=value')
					propName := fmt.Sprintf("%v", structData[idx]["Name"]) 
					propVal := fmt.Sprintf("%v", structData[idx]["Value"])
					propCompare := propName + "=" + propVal

					for k := 0; k < len(listKvProps); k++ {
						if propCompare == listKvProps[k] {
							foundList = append(foundList, listArtifUris[a])
							common.LogTxtHandler().Debug("Property found: " + listArtifUris[a])
						}
					}
				}
			}
		}

		if len(foundList) > 1 {
			// Count the occurance of duplicate artifacts and return a map of the artifact and duplicate count
			countMap := common.ReturnWithDupCounts(foundList)
			strCountMap := fmt.Sprintf("%v", countMap)
			common.LogTxtHandler().Debug("Count of duplicate artifacts and duplication count: " + strCountMap)

			for str, count := range countMap {
				// If the number of duplicate artifacts found matches the number of input property key/value pairs, add them to a filter list
				if count == numProps {
					filteredList = append(filteredList, str)
					common.LogTxtHandler().Debug("ARTIFACT FOUND WITH MATCHED PROPERTIES: " + str)
				}
			}
			// If only one item resulted in the filtered list, we will return it
			if len(filteredList) == 1 {
				foundItem = filteredList[0]
				common.LogTxtHandler().Info("FOUND ITEM: " + filteredList[0])
				return foundItem, nil
			} else {
				// For each artifact in the filter list, we grab it's 'created' date and add that artifact and date to an array of maps
				common.LogTxtHandler().Warn("More than one artifact with matching properties was found.")
				common.LogTxtHandler().Warn("Getting latest artifact...")

				foundItem, err = GetLatestArtifactFromList(filteredList)
				if err != nil {
					common.LogTxtHandler().Error("Error getting latest created date.")
				}
				return foundItem, nil
			}
		} else if len(foundList) == 1 {
			if numProps == len(foundList) {
				foundItem = foundList[0]
				common.LogTxtHandler().Info("FOUND ARTIFACT: " + foundItem)
				return foundItem, nil
			} else {
				err := errors.New("Artifacts found with at least one matching property. But no artifact was found with all properties.")
				common.LogTxtHandler().Error("Artifacts found with at least one matching property. But no artifact was found with all properties.")
				return "", err
			}
			
		} else {
			err := errors.New("No matching artifacts were found.")
			common.LogTxtHandler().Error("No matching artifacts were found.")
			return "", err
		}
	}
	return foundItem, nil
}

func SetArtifactProps(artifUri string, listKvProps []string) (string, error) {
	// Inputs are CASE SENSITIVE
	_, bearer := common.AuthCreds()
	requestPath := artifUri + "?properties="
	common.LogTxtHandler().Info(">>> Setting Specified Property(ies) for: " + artifUri)

	if common.ContainsSpecialChars(listKvProps) == true {
		err := errors.New("Properties cannot contain special characters --> )( }{ ][ *+^$\\/~`!@#%&<>;, and SPACE")
		common.LogTxtHandler().Error("Special character found.")
		common.LogTxtHandler().Error("Properties cannot contain special characters --> )( }{ ][ *+^$\\/~`!@#%&<>;, and SPACE")
		return "", err
	} else {
		if artifUri != "" && len(listKvProps) != 0 {
			// Determines whether we will format a list of property keys/values first, or pass a single property key/value pair
			// before making the API call
			if len(listKvProps) > 1 {
				// If there's more than one property keys/values supplied, adds the required ';' separater between them
				strProps := strings.Join(listKvProps, ";")
				common.LogTxtHandler().Debug("PROPERTIES TO BE PASSED: " + strProps)

				request, err = http.NewRequest("PUT", requestPath + strProps, nil)
				common.LogTxtHandler().Debug("REQUEST: Sending 'PUT' request to: " + requestPath + strProps)
			} else if len(listKvProps) == 1 && listKvProps[0] != "" {
				request, err = http.NewRequest("PUT", requestPath + listKvProps[0], nil)
				common.LogTxtHandler().Debug("REQUEST: Sending 'PUT' request to: " + requestPath + listKvProps[0])
			} else {
				err := errors.New("Unable to set Artifact properties without one or more property names and values.")
				common.LogTxtHandler().Error("Unable to set Artifact properties without one or more property names and values.")
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

				// If the request is successful, it will simply return a status code of 204
				if response.StatusCode == 204 {
					common.LogTxtHandler().Info("Request completed successfully")
					statusCode = "204"
				} else {
					// If the request fails, it will return a status code of 400
					common.LogTxtHandler().Info("Unable to complete request")
					statusCode = "400"
				}
			}
		} else {
			numProps := len(listKvProps)
			if numProps != 0 {
				err := errors.New("Unable to set Artifact properties without artifact's URI.")
				common.LogTxtHandler().Error("No artifact URI provided. Unable to set Artifact properties without artifact's URI.")
				return "", err
			} else {
				err := errors.New("Unable to set Artifact properties without artifact's URI and one or more property names/values.")
				common.LogTxtHandler().Error("No property names/values provided. Unable to set Artifact properties without artifact's URI and one or more property names/values.")
				return "", err
			}
		}
	}

	if err != nil {
		common.LogTxtHandler().Error("Unable to parse URL")
		return "", err
	}

	return statusCode, nil
}

func DeleteArtifactProps(artifUri string, listProps []string) (string, error) {
	// Inputs are CASE SENSITIVE
	// If a property is provided that doesn't exist (which includes incorrectly cased properties), the API ignores this and will return a successful response
	_, bearer := common.AuthCreds()
	requestPath := artifUri + "?properties="
	common.LogTxtHandler().Info(">>> Deleting Specified Property(ies) for Artifact: " + artifUri)

	if artifUri != "" && len(listProps) != 0 {
		// Determines whether we will format a list of property keys first, or pass a single property key
		// before making the API call
		if len(listProps) > 1 {
			// If there's more than one property keys supplied, adds the required ',' separater between them
			strProps := strings.Join(listProps, ",")
			common.LogTxtHandler().Debug("PROPERTIES TO BE PASSED: " + strProps)
			common.LogTxtHandler().Debug("REQUEST: Sending 'DELETE' request to: " + requestPath + strProps)

			request, err = http.NewRequest("DELETE", requestPath + strProps, nil)
		} else if len(listProps) == 1 && listProps[0] != "" {
			request, err = http.NewRequest("DELETE", requestPath + listProps[0], nil)
			common.LogTxtHandler().Debug("REQUEST: Sending 'DELETE' request to: " + requestPath + listProps[0])
		} else {
			err := errors.New("Unable to delete Artifact properties without one or more property names.")
			common.LogTxtHandler().Error("Unable to delete Artifact properties without one or more property names.")
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

			// If the request is successful, it will simply return a status code of 204
			if response.StatusCode == 204 {
				common.LogTxtHandler().Info("Request completed successfully")
				statusCode = "204"
			} else {
				// If the request fails, it will return a status code of 400
				common.LogTxtHandler().Info("Unable to complete request")
				statusCode = "400"
			}
		}
	} else {
		numProps := len(listProps)
		if numProps != 0 {
			err := errors.New("Unable to delete Artifact properties without artifact URI.")
			common.LogTxtHandler().Error("No artifact URI provided. Unable to delete Artifact properties without artifact URI.")
			return "", err
		} else {
			err := errors.New("Unable to delete Artifact properties without artifact URI and one or more property names.")
			common.LogTxtHandler().Error("No artifact properties provided. Unable to delete Artifact properties without artifact URI and one or more property names.")
			return "", err
		}
	}

	if err != nil {
		common.LogTxtHandler().Error("Unable to parse URL")
		return "", err
	}

	return statusCode, nil
}