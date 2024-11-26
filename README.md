# go-artifactory

## Summary
This module is a collection of Golang functions used to interact with JFrog Artifactory. While the functions can be called independently, the intent is to use them in support of a custom Packer plugin integration.

## Pre-Requisites
To run functions from this module, the following pre-requisites must be met:
1. An instance of Artifactory up and running, configured, and setup with at least a Pro license (which allows making REST API calls against it)
2. An account (such as a service account) on the Artifactory instance that can read and write to target repos, and that can read/write/delete artifacts from those target repos
3. An Identity Token created off the Artifactory account that will be running the functions/plugin
4. Go is installed on the system where the functions/plugin will be run: https://go.dev/doc/install
5. Ensure the GOPATH (the local directory to where the Go app is installed) is configured, as well as the environment variable to the path of the Go binary (on Windows, this would typically be C:\Program Files\Go\bin, for example)
6. Configure the .env file with the Artifactory Identity Token and Artifactory server; the Output Directory is optional but helpful when downloading artifacts to ensure they're placed in the desired location. Alternatively, they will be placed in a folder under the root of this module.

## About
This module is broken into several packages: common, operations, and search based on the underlying behavior of the functions. Some functions are specifically related to certain behaviors so they have been grouped together into packages as described below. 

What is returned from these functions (or what they can perform against artifacts) is entirely dependent on the permissions of the account running the functions. If the account can only see a single repo but 1,000 repos exist, then running the `ListRepos` function is only going to return a single repo.

### Common
These functions perform small, generalized supporting tasks for the other behavior-specific modules. These functions can be found under the **common.go** file.

### Operations
These functions are related operational-type behaviors. 

`properties` - Functions related to operational actions involving properties can be found under the **properties.go** file. This would be functions such as GETTING specific property values of an artifact, GETTING all properties/values of a given artifact, FILTERING artifacts by properties/values, SETTING property values, and DELETING property values.

`general` - Functions related to more general operational actions can be found under the **general.go** file. This would be functions such as LISTING all repos, GETTING all child objects of an item, GETTING the path to an artifact, GETTING the download URI, GETTING the created date of an artifact, RETRIEVING (downloading) an artifact, UPLOADING a new artifact, and DELETING an artifact from Artifactory. 

### Search
These functions are related specifically to searching for one or many artifacts. There's multiple ways to do this and how that's done is dependent on the information provided. These functions can be found under the **search.go** file. Functions such as GETTING a list of artifacts by a certain property(ies), GETTING a list of artifacts by name, and FILTERING a list of artifacts by file type would be found here.

### Archive
Archive also exists as a package, but it's really just a place to hold potentially useful functions that were created but have no immediate use. These artifacts are found under the **archive-search.go** file as they are related to searching for artifacts. They are more specific to artifacts that make use of Layouts, which may/may not be the case and would result in different behaviors or errors if used against artifacts that did not use Layouts. Therefore, more generalized operations and search capabilities were favored instead.

## The Functions
The following outlines the behavior of each function and any special notes.

## common/AuthCreds
Uses a `.env` file to capture the target Artifactory server and Artifactory account Identity Token

### Inputs
| Name        | Description                                                                              | Type   | Required |
| TOKEN       | Identity Token for the Artifactory account executing the function calls                  | string | TRUE     |
| ARTIFSERVER | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`          | string | TRUE     |
| OUTPUTDIR   | Desired directory to output file/artifact to, such as in `RetrieveArtifact` operations   | string | FALSE    |
|             | * If not specified, file will be dropped at the top-level directory of this module       |        |          |

### Outputs
| Name        | Description                                                                      | Type   |
| artifServer | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`  | string |
| bearer      | Forms bearer token to be passed with REST API Call to Artifactory                | string |
