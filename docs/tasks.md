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

## SetupTest
As part of the Artifactory plugin acceptance test, this function takes in the Artifactory server's API address, Artifactory Identity token, the full path to the test artifact that gets created (which is created from the plugin - ex: test-artifact.txt in the HOME directory of the user running the acceptance test), artifact suffix (optional, such as '1.0.0' or a date string, etc), and key/value pair of test properties (ex: release=latest-stable). The gloabl variables `util.ServerApi` and `util.Token` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

In addition, the function takes in a boolean value for whether or not to upload the test artifact, which allows for more flexibility when setting up the test environment, depending on the acceptance test. For example, the data source test requires the test artifact to be uploaded as part of the setup prep, while the post-processor for artifact uploads only needs the test repo to exist first.

Once the test directory is created along with the test artifact, and variables are set, a test repo called `/test-packer-plugin` is created within the Artifactory instance. The test artifact `$HOME/test-artifact.txt` is then uploaded, where applicable, to the test repo that was just created. If successful, the download URI is made available. 

From there, the artifact URI is derived from the download URI and the key/value properties are set on the test artifact. At this point, the test environment is ready for acceptance testing. If only the test repo needed to be created, the string 'Completed' is returned instead.

#### Inputs
| Name             | Description                                                                                               | Type     | Required |
|------------------|-----------------------------------------------------------------------------------------------------------|----------|:--------:|
| serverApi        | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`                           | string   | TRUE     |
| token            | Identity Token for the Artifactory account executing the function calls                                   | string   | TRUE     |
| testArtifactPath | Path to the test artifact created at start of acceptance testing; $HOME/test-directory/test-artifact.txt  | string   | TRUE     |
| artifactSuffix   | Full or partial name of the artifact to search for                                                        | string   | FALSE    |
| kvProps          | One or more property keys and values to filter by                                                         | []string | TRUE     |
| uploadArtifact   | true/false; whether or not to upload the test artifact as part of the test environment setup              | bool     | TRUE     |

#### Outputs
| Name          | Description                                                | Type     |
|---------------|------------------------------------------------------------|----------|
| *artifactUri  | Artifact URI address of the newly created test artifact    | string   |
| *(status)     | 'Complete'/'Incomplete'; Status of repo setup              | string   |
*If test repo was created AND test artifact was uploaded, then the Artifact URI will be returned. If only the test repo was created, then string value 'Completed' will be returned (or 'Incomplete' if there was an error).

## TeardownTest
As part of the Artifactory plugin acceptance test, this function takes in the Artifactory server's API address and Artifactory Identity token and deletes the test repo `/test-packer-plugin` along with the test artifact. From the plugin, the `test-directory/test-artifact.txt` created in the user's HOME directory at the start of the acceptance test is also deleted.

#### Inputs
| Name             | Description                                                                     | Type     | Required |
|------------------|---------------------------------------------------------------------------------|----------|:--------:|
| serverApi        | URL to the target Artifactory server; format: `server.com:8081/artifactory/api` | string   | TRUE     |
| token            | Identity Token for the Artifactory account executing the function calls         | string   | TRUE     |

#### Outputs
| Name       | Description                                                           | Type     |
|------------|-----------------------------------------------------------------------|----------|
| statusCode | Returns status code "200" if teardown is successful, or "400" if not  | string   |


## UploadArtifact
Takes in the Artifactory server's API address, Artifactory Identity token, source path of the new artifact, target path within Artifactory where the new artifact should be uploaded to, and optionally a file suffix if using the same artifact base name and needing to make it unique (ex: version, date, etc separated by '-'). The Global Variables `util.ServerApi` and `util.Token` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

Once the variables are set, the file is uploaded from the provided source path (`c:\\lab\artifact.ext` or `/lab/artifact.ext` to the target path (`/repo/folder/path`). If successful, the download URI of the returned from this operation. Next, the artifact URI is derived from the download URI. Both the download URI and artifact URI are returned

#### Inputs
| Name        | Description                                                                     | Type     | Required |
|-------------|---------------------------------------------------------------------------------|----------|:--------:|
| serverApi   | URL to the target Artifactory server; format: `server.com:8081/artifactory/api` | string   | TRUE     |
| token       | Identity Token for the Artifactory account executing the function calls         | string   | TRUE     |
| sourcePath  | Full file path where will be sourced from; **Needs proper escape chars          | string   | TRUE     |
| targetPath  | Target repo and folder destination of the artifact                              | string   | TRUE     |
| fileSuffix  | Placeholder for distinguishing values like dates, versions, etc                 | string   | FALSE    |

#### Outputs
| Name         | Description                                                                   | Type     |
|--------------|-------------------------------------------------------------------------------|----------|
| downloadUri  | Download URI address used to retrieve/download the artifact from Artifactory  | string   |
| artifactUri  | Artifact URI address of the artifact and it's details within Artifactory      | string   |


## SetProps
Takes in the Artifactory server's API address, Artifactory Identity token, artifact URI address, and one or more key/value property pairs. The Global Variables `util.ServerApi` and `util.Token` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

Once the variables are set, the artifact is assigned the new properties and a status code of "200" or "400" is returned depending on success or failure of the operation.

#### Inputs
| Name        | Description                                                                     | Type     | Required |
|-------------|---------------------------------------------------------------------------------|----------|:--------:|
| serverApi   | URL to the target Artifactory server; format: `server.com:8081/artifactory/api` | string   | TRUE     |
| token       | Identity Token for the Artifactory account executing the function calls         | string   | TRUE     |
| artifactUri | Artifact URI address of the newly created test artifact                         | string   | TRUE     |
| kvProps     | One or more property keys and values to assign to the artifact                  | []string | TRUE     |

#### Outputs
| Name       | Description                                                           | Type     |
|------------|-----------------------------------------------------------------------|----------|
| statusCode | Returns status code "200" if teardown is successful, or "400" if not  | string   |