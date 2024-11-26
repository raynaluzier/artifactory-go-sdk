# Common Functions

## AuthCreds
Uses a `.env` file to capture the target Artifactory server and Artifactory account Identity Token

#### Inputs
| Name        | Description                                                                              | Type   | Required |
|-------------|------------------------------------------------------------------------------------------|--------|:--------:|
| TOKEN       | Identity Token for the Artifactory account executing the function calls                  | string | TRUE     |
| ARTIFSERVER | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`          | string | TRUE     |
| OUTPUTDIR   | Desired directory to output file/artifact to, such as in `RetrieveArtifact` operations   | string | FALSE    |
|             | * If not specified, file will be dropped at the top-level directory of this module       |        |          |

#### Outputs
| Name        | Description                                                                      | Type   |
|-------------|----------------------------------------------------------------------------------|--------|
| artifServer | URL to the target Artifactory server; format: `server.com:8081/artifactory/api`  | string |
| bearer      | Forms bearer token to be passed with REST API Call to Artifactory                | string |