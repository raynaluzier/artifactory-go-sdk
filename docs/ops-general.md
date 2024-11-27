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


## GetItemChildren
Returns the children of the given item and whether that child object is a folder nor not (bool). The item can represent a repo name or a combo of repo/child_folder/subchild_folder/etc. If the item is the FULL path and filename to the artifact itself, no results will be returns as artifacts do not have children. However, artifacts can be children themselves.

Details of the child item, including it's child `Uri` (in this case '/folder' or '/file.ext') and `IsFolder` (true/false) values, are returned in a custom type called 'Contents'. This is used with the `GetArtifactPath` and `RecursiveSearch` functions to get and return the path to an artifact.

#### Inputs
| Name    | Description                                        | Type     | Required |
|---------|----------------------------------------------------|----------|:--------:|
| item    | Represent the repo name, folder, or file to check  | string   | TRUE     |

#### Outputs
| Name         | Description                                   | Type         |
|--------------|-----------------------------------------------|--------------|
| childDetails | Resulting string converted to lowercase       | []Contents   |
| err          | nil unless error; then returns error          | error        |


## GetArtifactPath
Takes in an artifact's name and searches Artifactory, returning the path to the artifact. These searches are CASE SENSITIVE. 
A path will be returned for every artifact FILE whose name includes the search string (e.g. paths for both 'win2022' and 'win2022-iis" would both be returned). Therefore, providing a partial name could result in multiple unintended paths being returned.

Additionally, multiple version files for a given artifact will result in the same path being added to the list multiple times. So we will search for and remove duplicates before returning the results.

** TODO: Check for both cases; add properties search as an option to further filter results

#### Inputs
| Name       | Description                              | Type     | Required |
|------------|------------------------------------------|----------|:--------:|
| artifName  | Name of the artifact being searched for  | string   | TRUE     |

#### Outputs
| Name         | Description                                                             | Type       |
|--------------|-------------------------------------------------------------------------|------------|
| listOfPaths  | List of found path(s) matching the full/partial name of the artifact    | []string   |
| err          | nil unless error; then returns error                                    | error      |


## RecursiveSearch
Recursively searches a list of child items for the specified artifact name. For each child item in the list, if the item isn't a folder, the process checks if the child item contains the desired artifact name. If so, the matching item's path will be added to the `foundPath` list. If not, the search path will be updated to check the next layer down, and the search will run again against the new search path.

The `Contains` function is case sensitive, so if the child item is NOT a folder, then provided artifact name will be converted to both upper and lowercase. Then a search will be done on the provided artifact name as it was provided originally to see if the child item is that artifact. If not, the check will be repeated with both the upper and lowercase versions of the name. If found, the path of that artifact will be returned.

Depending how artifacts are named (for ex: myfile-1.0.0.iso, myfile-1.1.1.iso, myfile-2.0.0.iso), searching for 'myfile' will result in multiple (and likely duplicate) paths to myfile-1.0.0.iso, myfile-1.1.1.iso, and myfile-2.0.0.iso being returned. We must capture these possiblities and do further filtering. Therefore, the paths are added to a list of strings called `foundPaths` to be processed later.

#### Inputs
| Name       | Description                                        | Type        | Required |
|------------|----------------------------------------------------|-------------|:--------:|
| list       | List of child items to search                      | []Contents  | TRUE     |
| artifName  | Name of the artifact to find                       | string      | TRUE     |
| searchPath | Path to child item being searched                  | string      | TRUE     |
| foundPaths | List of path(s) that contain target artiface name  | []string    | TRUE     |

#### Outputs
| Name       | Description                                        | Type       |
|------------|----------------------------------------------------|------------|
| foundPaths | List of path(s) that contain target artiface name  | []string   |


## GetDownloadUri
Requires full path to the artifact, including artifact name with extension. This function gets the artifact details and will return the download URI used to retrieve (download) the artifact.

#### Inputs
| Name          | Description                                            | Type    | Required |
|---------------|--------------------------------------------------------|---------|:--------:|
| artifPath     | Full repo/folder path to the target artifact           | string  | TRUE     |
| artifNameExt  | Full name of the artifact, including file extension    | string  | TRUE     |

#### Outputs
| Name        | Description                                        | Type     |
|-------------|----------------------------------------------------|----------|
| downloadUri | List of path(s) that contain target artiface name  | string   |
| err         | nil unless error; then returns error               | error    |


## GetCreateDate
Requires full path to the artifact, including artifact name with extension. This function gets the artifact details and will return the string date `created`.

#### Inputs
| Name          | Description                                              | Type    | Required |
|---------------|----------------------------------------------------------|---------|:--------:|
| artifactUri   | URI of the artifact itself (different from Download URI) | string  | TRUE     |

#### Outputs
| Name        | Description                           | Type     |
|-------------|---------------------------------------|----------|
| createdDate | Date/time the artifact was created    | string   |
| err         | nil unless error; then returns error  | error    |


## RetrieveArtifact
This function gets the artifact via the provided Download URI and copies it to the output directory specified in the environment variables file (.env). If no output directory path was provided, the artifact will be downloaded to the top-level directory of this code.

#### Inputs
| Name          | Description                                                   | Type    | Required |
|---------------|---------------------------------------------------------------|---------|:--------:|
| downloadUri   | URI of the artifact that allows the artifact to be downloaded | string  | TRUE     |

#### Outputs
| Name        | Description                                                 | Type     |
|-------------|-------------------------------------------------------------|----------|
| (msg) | String message indicating completion or failure of the download   | string   |
| err         | nil unless error; then returns error                        | error    |


## UploadFile
Uploads artifact to specified target path. 
`sourcePath` should be properly escaped and in the format of 'h:\\lab\\artifact.txt' or 
/lab/artifact.txt. 
`targetPath` should be in the format of '/repo-key/folder/path/'
The target filename will match the source file as it exists in the source directory.
`fileSuffix` is an optional placeholder for potential distinguishing values such as versions, etc. where a common artifact identifier (such as 'win2022') is used for every build and some other distinguishing value should be appended for uniquiness. If an empty string ("") is passed, then this will be ignored.

#### Inputs
| Name          | Description                                                           | Type    | Required |
|---------------|-----------------------------------------------------------------------|---------|:--------:|
| sourcePath  | Full file path where will be sourced from; **Needs proper escape chars  | string  | TRUE     |
| targetPath  | Target repo and folder destination of the artifact                      | string  | TRUE     |
| fileSuffix  | Placeholder for distinguishing values like dates, versions, etc         | string  | *FALSE   |
                *If not using a file suffix, "" (empty string) should be passed

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