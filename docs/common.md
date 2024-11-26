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
| Name        | Description                                                                              | Type     | Required |
|-------------|------------------------------------------------------------------------------------------|----------|:--------:|
| inputStr    | Input string that should be converted to lowercase                                       | string   | TRUE     |

#### Outputs
| Name        | Description                                                                              | Type     |
|-------------|------------------------------------------------------------------------------------------|----------|
| lowerStr    | Resulting string converted to lowercase                                                  | string   |


## ConvertToUppercase
Converts input string to uppercase and returns converted value.

#### Inputs
| Name        | Description                                                                              | Type     | Required |
|-------------|------------------------------------------------------------------------------------------|----------|:--------:|
| inputStr    | Input string that should be converted to uppercase                                       | string   | TRUE     |

#### Outputs
| Name        | Description                                                                              | Type     |
|-------------|------------------------------------------------------------------------------------------|----------|
| upperStr    | Resulting string converted to uppercase                                                  | string   |


## RemoveDuplicateStrings
Searches list of provided strings and removes any duplicates. 

#### Inputs
| Name           | Description                                                                           | Type     | Required |
|----------------|---------------------------------------------------------------------------------------|----------|:--------:|
| listOfStrings  | List of strings to check for duplicates                                               | []string | TRUE     |

#### Outputs
| Name        | Description                                                                              | Type     |
|-------------|------------------------------------------------------------------------------------------|----------|
| list        | Resulting list of strings will all duplicates removed                                    | []string |