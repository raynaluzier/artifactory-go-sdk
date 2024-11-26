# Common Functions

## AuthCreds
Uses a `.env` file to capture the target Artifactory server and Artifactory account Identity Token, and then returns the variable for the Artifactory server and bearer token to be used in subsequent REST API calls.

#### Inputs
| Name        | Description                                                                              | Type     | Required |
|-------------|------------------------------------------------------------------------------------------|----------|:--------:|
| TOKEN       | Identity Token for the Artifactory account executing the function calls                  | string   | TRUE     |
| ARTIFSERVER | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`          | string   | TRUE     |
| OUTPUTDIR   | Desired directory to output file/artifact to, such as in `RetrieveArtifact` operations   | string   | FALSE    |
|             | * If not specified, file will be dropped at the top-level directory of this module       |          |          |

#### Outputs
| Name        | Description                                                                              | Type     |
|-------------|------------------------------------------------------------------------------------------|----------|
| artifServer | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`          | string   |
| bearer      | Forms bearer token to be passed with REST API Call to Artifactory                        | string   |


## ConvertToLowercase
Converts input string to lowercase and returns converted value.

#### Inputs
| Name        | Description                                           | Type     | Required |
|-------------|-------------------------------------------------------|----------|:--------:|
| inputStr    | Input string that should be converted to lowercase    | string   | TRUE     |

#### Outputs
| Name        | Description                                           | Type     |
|-------------|-------------------------------------------------------|----------|
| lowerStr    | Resulting string converted to lowercase               | string   |


## ConvertToUppercase
Converts input string to uppercase and returns converted value.

#### Inputs
| Name        | Description                                           | Type     | Required |
|-------------|-------------------------------------------------------|----------|:--------:|
| inputStr    | Input string that should be converted to uppercase    | string   | TRUE     |

#### Outputs
| Name        | Description                                | Type     |
|-------------|--------------------------------------------|----------|
| upperStr    | Resulting string converted to uppercase    | string   |


## RemoveDuplicateStrings
Searches list of provided strings and removes any duplicates. 

#### Inputs
| Name           | Description                              | Type     | Required |
|----------------|------------------------------------------|----------|:--------:|
| listOfStrings  | List of strings to check for duplicates  | []string | TRUE     |

#### Outputs
| Name        | Description                                            | Type     |
|-------------|--------------------------------------------------------|----------|
| list        | Resulting list of strings will all duplicates removed  | []string |


## ReturnWithDupCounts
Counts occurances of each string and returns a map of strings and a count of the number of times a duplicate instance of that string was found.

#### Inputs
| Name           | Description                              | Type     | Required |
|----------------|------------------------------------------|----------|:--------:|
| listOfStrings  | List of strings to check for duplicates  | []string | TRUE     |

#### Outputs
| Name        | Description                                        | Type           |
|-------------|----------------------------------------------------|----------------|
| countMap    | Resulting map of strings and duplication counts    | map[string]int |


## ReturnDuplicates
Takes in the `countMap` map of strings and number of occurances (ex: map[str1:1, str2:5, str3:1]). For any strings with more than one occurance, the string is added to the 'duplicates' list and returned

#### Inputs
| Name      | Description                           | Type           | Required |
|-----------|---------------------------------------|----------------|:--------:|
| countMap  | Map of string and duplication counts  | map[string]int | TRUE     |

#### Outputs
| Name        | Description                                                  | Type     |
|-------------|--------------------------------------------------------------|----------|
| duplicates  | Resulting list of strings that have more than one occurance  | []string |


## SetArtifUriFromDownloadUri
Some artifact operations can take either the Artifact's URI or it's download URI, which are slightly different URI strings. However, some operations cannot use the resulting download URI, so this function allows us a quick way to get the artifact's URI in instances where all we have is the download URI. 

#### Inputs
| Name        | Description                   | Type   | Required |
|-------------|-------------------------------|--------|:--------:|
| downloadUri | Download URI of the artifact  | string | TRUE     |

#### Outputs
| Name        | Description                  | Type     |
|-------------|------------------------------|----------|
| artifUri  | Resulting URI of the artifact  | string   |


## SearchForExactString
Searches an input string for an exact search team; for example: Search term "win2022" will return TRUE if the input string is "win2022" and FALSE if "win2022-iis".

#### Inputs
| Name        | Description                             | Type   | Required |
|-------------|-----------------------------------------|--------|:--------:|
| searchTerm  | The string name of an actual object     | string | TRUE     |
| inputStr    | The string provided through user input  | string | TRUE     |

#### Outputs
| Name    | Description                                    | Type  |
|---------|------------------------------------------------|-------|
| result  | True/false whether the strings matched exactly | bool  |
| err     | nil unless error; then returns error           | error |


## EscapeSpecialChars
Takes the input string (such as a directory path provided by an environment variable) and adds escape characters.
For example:  F:\mypath\ becomes F:\\mypath.

#### Inputs
| Name   | Description                                  | Type   | Required |
|--------|----------------------------------------------|--------|:--------:|
| input  | Input string to check; likely directory path | string | TRUE     |

#### Outputs
| Name    | Description                                    | Type   |
|---------|------------------------------------------------|--------|
| input   | Resulting string that's been properly escaped  | string |


## createJsonString
Used with  `EscapeSpecialChars` to create a JSONString from the input as part of the process to properly format string, like directories, that may include "\".

#### Inputs
| Name   | Description                                  | Type   | Required |
|--------|----------------------------------------------|--------|:--------:|
| input  | Input string to check; likely directory path | string | TRUE     |

#### Outputs
| Name       | Description                            | Type   |
|------------|----------------------------------------|--------|
| jsonString | Resulting JSON string of the input     | string |


## CheckPathType
Checks the provided path to see if it's Unix-based (has '/') or Windows-based (has '\'). This is used in combination with `CheckAddSlashToPath` to add the appropriate ending slash type to given path if needed.

#### Inputs
| Name  | Description                                          | Type   | Required |
|-------|------------------------------------------------------|--------|:--------:|
| path  | Path to provided directory; such as Output Directory | string | TRUE     |

#### Outputs
| Name       | Description                                        | Type |
|------------|----------------------------------------------------|------|
| isWinPath | Returns true of the provided path is Windows-based  | bool |


## StringCompare
Performs case INSENSITIVE comparison of strings (like file names) and returns TRUE if they match. This comparison does NOT do partial string comparisons, so 'win' and 'win-2022' would be false.

#### Inputs
| Name       | Description                                              | Type   | Required |
|------------|----------------------------------------------------------|--------|:--------:|
| inputStr   | String that was provided through some kind of user input | string | TRUE     |
| actualStr  | String pulled from actual object name                    | string | TRUE     |

#### Outputs
| Name       | Description                                               | Type |
|------------|-----------------------------------------------------------|------|
| ture/false | Returns true of compared string match, regardless of case | bool |


## CheckAddSlashToPath
Used with `CheckPathType`; based on path type (Windows vs. Unix), checks the provided path to see if it ends with the appropriate back or forward slashes. If not present, the function will add a slash as appropriate to the platform type. This ensures the output directory path provided is formatted as required.

#### Inputs
| Name       | Description                                              | Type   | Required |
|------------|----------------------------------------------------------|--------|:--------:|
| inputStr   | String that was provided through some kind of user input | string | TRUE     |
| actualStr  | String pulled from actual object name                    | string | TRUE     |

#### Outputs
| Name       | Description                                               | Type |
|------------|-----------------------------------------------------------|------|
| ture/false | Returns true if compared string match, regardless of case | bool |


## ContainsSpecialChars
Checks for the special characters that are disallowed by Artifactory in Properties. This function returns TRUE if any of the characters are found.

#### Inputs
| Name       | Description                                     | Type     | Required |
|------------|-------------------------------------------------|----------|:--------:|
| strings    | List of strings to check for special characters | []string | TRUE     |

#### Outputs
| Name       | Description                                                    | Type |
|------------|----------------------------------------------------------------|------|
| ture/false | Returns true if any of the strings contains special characters | bool |