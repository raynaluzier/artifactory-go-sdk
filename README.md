# artifactory-go-sdk

## Summary
This SDK is a collection of Golang functions used to interact with JFrog Artifactory. While the functions can be called independently, the intent is to use them with a custom Packer plugin integration.

## Pre-Requisites
To run functions from this module, the following pre-requisites must be met:

1. An instance of Artifactory up and running, configured, and setup with at least a Pro license (which allows making REST API calls against it)

2. An account (such as a service account) on the Artifactory instance that can read and write to target repos, and that can read/write/delete artifacts from those target repos

3. An Identity Token created off the Artifactory account that will be running the functions/plugin

4. Go is installed on the system where the functions/plugin will be run: https://go.dev/doc/install

5. Ensure the `GOPATH` (the local directory to where the Go app is installed) is configured, as well as the environment variable to the path of the Go binary (on Windows, this would typically be `C:\Program Files\Go\bin`, for example)

6. Populating the Global Variables (see: `/util/util.go`): **ServerApi**, **Token**, and optionally, **Logging** and **OutputDir**. This can be done when calling a function in the `tasks` package (e.g. as part of a plugin operation), by configuring the `.env` file, or statically (only recommended for testing).

Every Artifactory API call requires passing the bearer token, and several require the Artifactory API server. Therefore, these MUST be set.

Using .env File: Configure the `.env` file with the Artifactory Identity Token and Artifactory server; the Output Directory is optional but helpful when downloading artifacts to ensure they're placed in the desired location. Alternatively, they will be placed in the user's HOME directory. Logging provides an option to change the level of logging to display.

    * `ARTIFACTORY_TOKEN` - Artifactory Identity Token --> Ex:  ARTIFACTORY_TOKEN=1234567890abcdefghijklmnopqrstuv
    * `ARTIFACTORY_SERVER` - Artifactory Server --> Ex:  ARTIFACTORY_SERVER=https://server.com:8081/artifactory/api
    * `ARTIFACTORY_OUTPUTDIR` - Output directory for downloading artifacts --> Ex: ARTIFACTORY_OUTPUTDIR=H:\output-dir\path\ or /output-dir/path/
    * `ARTIFACTORY_LOGGING` - Logging level (INFO, WARN, ERROR, DEBUG); defaults to 'INFO' --> Ex:  ARTIFACTORY_LOGGING=DEBUG

    Then use `os.Getenv` to set `util.ServerApi`, `util.Token`, `util.Logging`, and `util.OutputDir` respectively.

## About
This SDK is broken into several packages: `common`, `operations`, `search`, `tasks`, and `util` based on the underlying behavior of the functions. Some functions are specifically related to certain behaviors so they have been grouped together into packages as described below. 

What is returned from these functions (or what they can perform against artifacts) is entirely dependent on the permissions of the account running the functions. If the account can only see a single repo but 1,000 repos exist, then running the `ListRepos()` function is only going to return a single repo.

### Common
These functions perform small, generalized supporting tasks for the other behavior-specific modules. These functions can be found under the `common.go` file.

### Operations
These functions are related operational-type behaviors. 

**properties** - Functions related to operational actions involving properties can be found under the `properties.go` file. This would be functions such as GETTING specific property values of an artifact, GETTING all properties/values of a given artifact, FILTERING artifacts by properties/values, SETTING property values, and DELETING property values.

**general** - Functions related to more general operational actions can be found under the `general.go` file. This would be functions such as LISTING all repos, GETTING all child objects of an item, GETTING the path to an artifact, GETTING the download URI, GETTING the created date of an artifact, RETRIEVING (downloading) an artifact, UPLOADING a new artifact, and DELETING an artifact from Artifactory. 

### Search
These functions are related specifically to searching for one or many artifacts. There's multiple ways to do this and how that's done is dependent on the information provided. These functions can be found under the `search.go` file. Functions such as GETTING a list of artifacts by a certain property(ies), GETTING a list of artifacts by name, and FILTERING a list of artifacts by file type would be found here.

### Tasks
These functions are larger operations that first set the global variables, and then make a series of function calls to perform specific activities. While they can be called independently, they were created in support of a custom Packer plugin to streamline passing environment-specific variables, such as the Artifactory token, server, logging, and output directory. Rather than passing one or more of these to every function in the SDK (in addition to the required inputs), they are passed in ONCE to the desired function, the global variables are set, and then they are used automatically when calling each sub-function without having to pass them in over and over.

These larger tasks also group the targeted functions of a desired behavior into a single operation and keep the plugin code to a minimum and simplify performing that desired behavior. For example, finding an image/artifact and returning it's name, created date, and download URI involves six (6) different function calls and passing in specific information. Using the `GetImageDetails()` function is just a single call which handles those underlying function calls 'behind the scenes'.

### Utils
This is a list of the global variables used within this SDK. As with any Go package, they can be used by importing the `util` package path and then referencing them as `util.Token`, `util.ServerApi`, etc.

### Archive
Archive also exists as a package, but it's really just a place to hold potentially useful functions that were created but have no immediate use. Artifacts are split into archive files that match their associated behaviors; as of now, either the `archive-general.go` or `archive-search.go` files. 

The `archive-search.go` file contains functions that are more specific to artifacts that make use of Layouts, which may/may not be the case and would result in different behaviors or errors if used against artifacts that did not use Layouts. Therefore, more generalized operations and search capabilities were favored instead.

The `archive-general.go` file contains functions related to finding a specific artifact and then returning it's file path through recursive searches. Instead, finding the artifact by name, then filtering by file type, and optionally filtering by one or more specific properties/values was easier and more accurate. Therefore, the path-related functions were archived. 

## Function Reference
A reference outline of each function's behavior and any special notes can be found in the corresponding documents below.

- [Common](https://github.com/raynaluzier/artifactory-go-sdk/blob/main/docs/common.md)

- [Operations/General](https://github.com/raynaluzier/artifactory-go-sdk/blob/main/docs/ops-general.md)

- [Operations/Properties](https://github.com/raynaluzier/artifactory-go-sdk/blob/main/docs/ops-properties.md)

- [Search](https://github.com/raynaluzier/artifactory-go-sdk/blob/main/docs/search.md)

- [Tasks](https://github.com/raynaluzier/artifactory-go-sdk/blob/main/docs/tasks.md)


## How to Use

## Unit Testing