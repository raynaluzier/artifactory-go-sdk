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
This function gets the artifact via the provided Download URI and copies it to the output directory path specified in the `util.OutputDir` Global Variable. If no output directory path was provided, the artifact will be downloaded to the user's HOME directory.

#### Inputs
| Name                  | Description                                                   | Type    | Required |
|-----------------------|---------------------------------------------------------------|---------|:--------:|
| downloadUri           | URI of the artifact that allows the artifact to be downloaded | string  | TRUE     |

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


## GetArtifact
Takes in the download URI of an artifact and makes a 'GET' REST API call against that URI. A status code of "200" is returned if it exists or "404" if it doesn't.

#### Inputs
| Name        | Description                                           | Type    | Required |
|-------------|-------------------------------------------------------|---------|:--------:|
| downloadUri | The artifact's download address within Artifactory    | string  | TRUE     |

#### Outputs
| Name        | Description                                           | Type     |
|-------------|-------------------------------------------------------|----------|
| statusCode  | Result of the GET call to the download URI address    | string   |


## CheckFileAndUpload
Takes in a list of files pulled from the source directory (collected by the parent function, `UploadArtifacts`); for each item name in the list, checks it against the target filename. If the file exists, the filename is appended to the source path. The target path is updated to include the image name as the target folder. Then the file is uploaded to Artifactory. If successful, the artifact's download URI is output and the string-based result of the operation is returned as "Success".

This supports and is called by the `UploadArtifacts` function.

#### Inputs
| Name      | Description                                                  | Type           | Required |
|-----------|--------------------------------------------------------------|----------------|:--------:|
| items     | List of directory files read from source                     | []os.DirEntry  | TRUE     |
| sourceDir | Source directory that contains the files to upload           | string         | TRUE     |
| targetDir | Target Artifactory path where files will be uploaded to      | string         | TRUE     |
| fileName  | File to be checked that it exists and uploaded               | string         | TRUE     |
| imageName | Name of the image; used as a folder to place the image files | string         | TRUE     |

#### Outputs
| Name     | Description                                                            | Type     |
|----------|------------------------------------------------------------------------|----------|
| (result) | String result of either "Failed" or "Success" at the upload operation  | string   |


## CheckFileAndDownload
Takes the Artifactory artifact path (ex: /repo/folder) and image filename with extension, then does a 'GET' REST API call to Artifactory for the item. If the file is found, a status code of "200" is returned and the artifact is downloaded (it uses the output directory set as a global variable). Then the function returns a string result of "Success" back to the calling function if the download completed without error. Otherwise, "Failed" is returned.

This supports and is called by the `DownloadArtifacts` function.

#### Inputs
| Name         | Description                                                           | Type    | Required |
|--------------|-----------------------------------------------------------------------|---------|:--------:|
| checkFile    | The filename with extension we are checking                           | string  | TRUE     |
| downloadPath | Parsed Artifactory path to the artifact without the artifact filename | string  | TRUE     |
| task         | What vSphere disk file check we are performing                        | string  | TRUE     |

#### Outputs
| Name     | Description                                                            | Type     |
|----------|------------------------------------------------------------------------|----------|
| (result) | String result of either "Failed" or "Success" at the upload operation  | string   |


## CheckFileLoopAndDownload
There are several vSphere disk file types that may exist as part of a VM Template, such as .vmdk, -ctk.vmdk, -flat.vmdk. There are numbered and unnumbered versions of all of these files that may exist as well, depending on the number of disks that are attached to the template. The number of disks can vary from image to image. This function will form an expected disk file name to check for with a disk number that increments with each pass to a max of 15 (i.e. we're accounting for the possibility of a max of 15 attached disks to the template).

The disk file name is appended to the download path to form a URI in Artifactory and a 'GET' REST API call is done against the formed URI to see if the file exists in Artifactory. If the file is found, a status code of "200" is returned and the artifact is downloaded (it uses the output directory set as a global variable). Then the function returns a string result of "Success" back to the calling function if the download completed without error. Otherwise, "Failed" is returned.

If, as the disk number increments and is checked, that numbered disk file isn't found, we assume that we've reached the end of the number of disks that are attached to the VM Template and the process breaks.

This supports and is called by the `DownloadArtifacts` function. 

#### Inputs
| Name         | Description                                                           | Type    | Required |
|--------------|-----------------------------------------------------------------------|---------|:--------:|
| imageName    | Name of the image we'll use to construct the filename with            | string  | TRUE     |
| downloadPath | Parsed Artifactory path to the artifact without the artifact filename | string  | TRUE     |
| extString    | vSphere-based disk file extenison; ex: .vmdk, -ctk.vmdk, -flat.vmdk   | string  | TRUE     |
| task         | What vSphere disk file check we are performing                        | string  | TRUE     |

#### Outputs
| Name        | Description                                           | Type     |
|-------------|-------------------------------------------------------|----------|
| statusCode  | Result of the GET call to the download URI address    | string   |