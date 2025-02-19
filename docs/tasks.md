# Tasks Functions
As described previously, these functions are intended to be used with a custom Packer plugin, but can be called independently if desired.

## GetImageDetails
Takes in the Artifactory server's API address, Artifactory Identity token, desired log level (if other than 'INFO'), the full or partial artifact name, file extension, and optionally one or more property key/values. The Global Variables `util.ServerApi`, `util.Token`, and `util.Logging` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

Once the variables are set, `GetArtifactsByName` takes in the artifact name provided and returns a list of one or more artifact URIs that match. Next, `FilterListByFileType` filters this list by the file extension input (defaults to .vmtx if blank). If the result is only a single artifact URI, this artifact will be returned. 

If the artifact list contains more than one artifact AND one or more property keys/values were provided, then the list will be filtered by artifacts with the matching property(ies) via `FilterListByProps`. As before, if only one artifact matches, this artifact is returned.

If the resulting list of artifacts still contains more than one artifact, then this list is parsed and the latest artifact is returned. If no properties were provided to filter by, the list is parsed for the latest artifact.

#### Inputs
| Name        | Description                                                                       | Type     | Required |
|-------------|-----------------------------------------------------------------------------------|----------|:--------:|
| serverApi   | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`   | string   | TRUE     |
| token       | Identity Token for the Artifactory account executing the function calls           | string   | TRUE     |
| logLevel    | Ouput logging level; INFO, WARN, ERROR, DEBUG; defaults to 'INFO'                 | string   | FALSE    |
| artifName   | Full or partial name of the artifact to search for                                |          | TRUE     |
| ext         | File extension of the artifact; defaults to .vmtx if left blank                   | string   | TRUE     |
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


## UploadGeneralArtifact
Takes in the Artifactory server's API address, Artifactory Identity token, source path of the artifact, target path within Artifactory where the artifact should be uploaded to, file name of the artifact, and a folder name where the artifact should be placed. Folder name may be the same as the image name if placing the artifact with associated image files, or it could be a separate folder for flexibility. Otherwise, leave `folderName` blank and the file will be place directly in the `artifPath`.

The Global Variables `util.ServerApi` and `util.Token` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

Once the variables are set, the source file is verified that it exists in the directory and if so, uploaded from the provided source path to the target Artifactory path (`/repo/folder/path`). The result string of "Success" or "Failed" is returned.

#### Inputs
| Name        | Description                                                                     | Type     | Required |
|-------------|---------------------------------------------------------------------------------|----------|:--------:|
| serverApi   | URL to the target Artifactory server; format: `server.com:8081/artifactory/api` | string   | TRUE     |
| token       | Identity Token for the Artifactory account executing the function calls         | string   | TRUE     |
| sourcePath  | Directory where file will be sourced from; **Needs proper escape chars          | string   | TRUE     |
| artifPath   | Target repo and folder destination of the artifact                              | string   | TRUE     |
| fileName    | Name of the file with extenion being checked and uploaded                       | string   | TRUE     |
| folderName  | Optional folder name to place the file into                                     | string   | FALSE    |

#### Outputs
| Name   | Description                                                 | Type     |
|--------|-------------------------------------------------------------|----------|
| result | Resulting string of "Success" or "Failed" for the operation | string   |
| err    | Error message returned if upload failed                     | string   |


## DownloadGeneralArtifact
Takes in the Artifactory server's API address, Artifactory Identity token, desired output directory, Artifactory path within Artifactory where the artifact should be download from, file name of the artifact, and task string used for logging.

The Global Variables `util.ServerApi`, `util.Token`, and `util.OutputDir` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

Once the variables are set, the desired file is verified that it exists in the Artifactory path (ex: /repo/opt-folder/), and if so, it's downloaded to the output directory. The result string of "Success" or "Failed" is returned.

#### Inputs
| Name        | Description                                                                     | Type     | Required |
|-------------|---------------------------------------------------------------------------------|----------|:--------:|
| serverApi   | URL to the target Artifactory server; format: `server.com:8081/artifactory/api` | string   | TRUE     |
| token       | Identity Token for the Artifactory account executing the function calls         | string   | TRUE     |
| outputDir   | Directory where file will be downloaded to; **Needs proper escape chars         | string   | TRUE     |
| artifPath   | Target repo and folder destination of the artifact; ex: `/repo/opt-folder/`     | string   | TRUE     |
| fileName    | Name of the file with extenion being checked and uploaded                       | string   | TRUE     |
| task        | String that logs each file being downloaded for info purposes                   | string   | TRUE     |

#### Outputs
| Name   | Description                                                 | Type     |
|--------|-------------------------------------------------------------|----------|
| result | Resulting string of "Success" or "Failed" for the operation | string   |
| err    | Error message returned if upload failed                     | string   |


## UploadArtifacts
Takes in the Artifactory server's API address, Artifactory Identity token, image type (OVA, OVF, or VMTX), image name, source path of the new artifact (ex: c:\\lab or /lab), target path within Artifactory where the new artifact should be uploaded to (ex: /repo/opt-folder/), and optionally a file suffix if using the same artifact base name and needing to make it unique (ex: version, date, etc separated by '-'). 

The Global Variables `util.ServerApi` and `util.Token` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

Once the variables are set, the image type and image name are evaluated to determined the expected files that should exist. Disk files are evaluated for up to 15 disks. 

The files are validated against the source directory and if they exist, they are uploaded from the provided source path (`c:\\lab` or `/lab` to the target path (`/repo/folder/path`) into a folder based on the image name (so /repo/opt-folder/image1234/image1234.ova, etc.). As each file is successfully uploaded, the download URI is output in the logs. Upon completion, a string-based status of the operation is returned.

#### Inputs
| Name        | Description                                                                                                      | Type     | Required |
|-------------|------------------------------------------------------------------------------------------------------------------|----------|:--------:|
| serverApi   | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`                                  | string   | TRUE     |
| token       | Identity Token for the Artifactory account executing the function calls                                          | string   | TRUE     |
| logLevel    | Ouput logging level; INFO, WARN, ERROR, DEBUG; defaults to 'INFO'                 | string   | FALSE    |
| imageType   | Type of image to be uploaded (OVA, OVF, or VMTX are supported)                                                   | string   | TRUE     |
| imageName   | Base name of the image (ex: win2022)                                                                             | string   | TRUE     |
| sourceDir   | Directory path (without any filename) where the image will be sourced from; **Needs proper escape chars          | string   | TRUE     |
| targetDir   | Target repo/path of destination for image (files will automatically be placed in their own image-based folder)   | string   | TRUE     |
| fileSuffix  | Placeholder for distinguishing values like dates, versions, etc                                                  | string   | FALSE    |

#### Outputs
| Name      | Description                               | Type     |
|-----------|-------------------------------------------|----------|
| (result)  | Resulting status string of the operation  | string   |


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


## DownloadArtifacts
Takes in the Artifactory server's API address, Artifactory Identity token, download URI for the primary image file (OVA, OVF, or VMTX), and a desired output directory. If the image is going to be imported into a vCenter instance as part of the build process, then the output directory should be a datastore path available to the system where Packer is running, such as through a share. 

The Global Variables `util.ServerApi`, `util.Token`, and `util.Output` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

Once the variables are set, the download URI is parsed to determine the primary image's file name, extension, and image name. A folder will be created on the output directory named based on the image name. Next, the image type and image name are evaluated to determined the expected files that should exist. Disk files are evaluated for up to 15 disks. Each expected file is checked against Artifactory to ensure it exists, and if so, will be downloaded to the image's folder in the specified output directory. Shoud any of the files not be found, the process will exit with an error.

#### Inputs
| Name        | Description                                                                     | Type     | Required |
|-------------|---------------------------------------------------------------------------------|----------|:--------:|
| serverApi   | URL to the target Artifactory server; format: `server.com:8081/artifactory/api` | string   | TRUE     |
| token       | Identity Token for the Artifactory account executing the function calls         | string   | TRUE     |
| downloadUri | Download URI address of the primary image artifact (type: OVA, OVF, or VMTX)    | string   | TRUE     |
| outputDir   | Target directory where the downloaded files should be placed*                   | string   | TRUE     |
**If the image is going to be imported into a vCenter instance as part of the build process, then the output directory should be a datastore path available to the system where Packer is running, such as through a share.** 

#### Outputs
| Name      | Description                               | Type     |
|-----------|-------------------------------------------|----------|
| (result)  | Resulting status string of the operation  | string   |