# Tasks Functions
As described previously, these functions are intended to be used with a custom Packer plugin, but can be called independently if desired.

## GetImageDetails
Takes in the Artifactory server's API address, Artifactory Identity token, desired log level (if other than 'INFO'), the full or partial artifact name, file extension, and optionally one or more property key/values. The Global Variables `util.ServerApi`, `util.Token`, and `util.Logging` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

Once the variables are set, `GetArtifactsByName` takes in the artifact name provided and returns a list of one or more artifact URIs that match. Next, `FilterListByFileType` filters this list by the file extension input (defaults to .vmxt if blank). If the result is only a single artifact URI, this artifact will be returned. 

If the artifact list contains more than one artifact AND one or more property keys/values were provided, then the list will be filtered by artifacts with the matching property(ies) via `FilterListByProps`. As before, if only one artifact matches, this artifact is returned.

If the resulting list of artifacts still contains more than one artifact, then this list is parsed and the latest artifact is returned. If no properties were provided to filter by, the list is parsed for the latest artifact.

#### Inputs
| Name        | Description                                                                       | Type     | Required |
|-------------|-----------------------------------------------------------------------------------|----------|:--------:|
| serverApi   | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`   | string   | TRUE     |
| token       | Identity Token for the Artifactory account executing the function calls           | string   | TRUE     |
| logLevel    | Ouput logging level; INFO, WARN, ERROR, DEBUG; defaults to 'INFO'                 | string   | FALSE    |
| artifName   | Full or partial name of the artifact to search for                                |          | TRUE     |
| ext         | File extension of the artifact; defaults to .vmxt if left blank                   | string   | TRUE     |
| kvProps     | One or more property keys and values to filter by                                 | []string | FALSE    |
*Any inputs NOT required should pass in an empty variable to the function.*

#### Outputs
| Name         | Description                                                           | Type     |
|--------------|-----------------------------------------------------------------------|----------|
| artifactName | Full name of the artifact as it exists in Artifactory                 | string   |
| createDate   | Date the artifact was created within Artifactory                      | string   |
| downloadUri  | Download URI used to retrieve/download the artifact from Artifactory  | string   |