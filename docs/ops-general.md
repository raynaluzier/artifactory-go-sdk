# General Operations Functions

## ListRepos
Gets the list of repositories the Artifactory account's Identity Token has access to.

#### Inputs
Takes no inputs

#### Outputs
| Name      | Description                           | Type     |
|-----------|---------------------------------------|----------|
| listRepos | List of available repos               | []string |
| err       | nil unless error; then returns error  | error    |


## GetDownloadUri
Requires full path to the artifact, including artifact name with extension. This function gets the artifact details and will return the download URI used to retrieve (download) the artifact.

#### Inputs
| Name       | Description                                              | Type    | Required |
|------------|----------------------------------------------------------|---------|:--------:|
| artifUri   | URI of the artifact itself (different from Download URI) | string  | TRUE     |

#### Outputs
| Name        | Description                                        | Type     |
|-------------|----------------------------------------------------|----------|
| downloadUri | List of path(s) that contain target artiface name  | string   |
| err         | nil unless error; then returns error               | error    |


## GetArtifactNameFromUri
Takes in an artifact URI, parses, and returns the name of the artifact.

#### Inputs
| Name       | Description                                              | Type    | Required |
|------------|----------------------------------------------------------|---------|:--------:|
| artifUri   | URI of the artifact itself (different from Download URI) | string  | TRUE     |

#### Outputs
| Name         | Description            | Type     |
|--------------|------------------------|----------|
| artifactName | Name of the artifact   | string   |


## GetCreateDate
Requires full path to the artifact, including artifact name with extension. This function gets the artifact details and will return the string date `created`.

#### Inputs
| Name       | Description                                              | Type    | Required |
|------------|----------------------------------------------------------|---------|:--------:|
| artifUri   | URI of the artifact itself (different from Download URI) | string  | TRUE     |

#### Outputs
| Name        | Description                           | Type     |
|-------------|---------------------------------------|----------|
| createdDate | Date/time the artifact was created    | string   |
| err         | nil unless error; then returns error  | error    |


## RetrieveArtifact
This function gets the artifact via the provided Download URI and copies it to the output directory specified in the environment variables file (.env). If no output directory path was provided, the artifact will be downloaded to the user's HOME directory.

#### Inputs
| Name                  | Description                                                   | Type    | Required |
|-----------------------|---------------------------------------------------------------|---------|:--------:|
| downloadUri           | URI of the artifact that allows the artifact to be downloaded | string  | TRUE     |
| ARTIFACTORY_OUTPUTDIR | Output directory where artifact downloads will be stored      | string  | FALSE    |

#### Outputs
| Name     | Description                                                       | Type     |
|----------|-------------------------------------------------------------------|----------|
| (msg)    | String message indicating completion or failure of the download   | string   |
| err      | nil unless error; then returns error                              | error    |


## UploadFile
Uploads artifact to specified target path. 
`sourcePath` should be properly escaped and in the format of 'h:\\lab\\artifact.txt' or 
/lab/artifact.txt. 
`targetPath` should be in the format of '/repo-key/folder/path/'
The target filename will match the source file as it exists in the source directory.
`fileSuffix` is an optional placeholder for potential distinguishing values such as versions, etc. where a common artifact identifier (such as 'win2022') is used for every build and some other distinguishing value should be appended for uniquiness with "-" as the separator ('win2022-1.1.1.iso'). If an empty string ("") is passed, then this will be ignored.

#### Inputs
| Name          | Description                                                           | Type    | Required |
|---------------|-----------------------------------------------------------------------|---------|:--------:|
| sourcePath  | Full file path where will be sourced from; **Needs proper escape chars  | string  | TRUE     |
| targetPath  | Target repo and folder destination of the artifact                      | string  | TRUE     |
| fileSuffix  | Placeholder for distinguishing values like dates, versions, etc         | string  | *FALSE   |
                *If not using a file suffix, an empty string ("") should be passed

#### Outputs
| Name    | Description                                                       | Type     |
|---------|-------------------------------------------------------------------|----------|
| (msg)   | String message indicating completion or failure of the download   | string   |
| err     | nil unless error; then returns error                              | error    |


## DeleteArtifact
Takes in an artifact's URI and executes a delete operation against it.

#### Inputs
| Name          | Description                                              | Type    | Required |
|---------------|----------------------------------------------------------|---------|:--------:|
| artifactUri   | URI of the artifact itself (different from Download URI) | string  | TRUE     |

#### Outputs
| Name        | Description                                                           | Type     |
|-------------|-----------------------------------------------------------------------|----------|
| statusCode  | Resulting status code of the delete operation (either "204" or "404") | string   |
| err         | nil unless error; then returns error                                  | error    |


## GetLatestArtifactFromList
Takes in list of artifact URIs, gets the created date for each of them, and returns the latest artifact.

#### Inputs
| Name   | Description              | Type      | Required |
|--------|--------------------------|-----------|:--------:|
| list   | List of artifact URIs    | []string  | TRUE     |

#### Outputs
| Name        | Description                           | Type     |
|-------------|---------------------------------------|----------|
| latestItem  | Artifact with latest 'create' date    | string   |
| err         | nil unless error; then returns error  | error    |