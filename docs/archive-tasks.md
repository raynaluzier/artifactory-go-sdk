# Archived Tasks Functions
The following functions are currently not used but archived in case a need for them arises.

## UploadArtifact
Takes in the Artifactory server's API address, Artifactory Identity token, source path of the new artifact, and target path within Artifactory where the new artifact should be uploaded to. The Global Variables `util.ServerApi` and `util.Token` are set by the function's inputs so these values can be used by the subsequent function calls without having to pass them in every time.

Once the variables are set, the file is uploaded from the provided source path (`c:\\lab\\artifact.ext` or `/lab/artifact.ext` to the target path (`/repo/folder/path`). If successful, the download URI of the returned from this operation. Next, the artifact URI is derived from the download URI. Both the download URI and artifact URI are returned

#### Inputs
| Name        | Description                                                                     | Type     | Required |
|-------------|---------------------------------------------------------------------------------|----------|:--------:|
| serverApi   | URL to the target Artifactory server; format: `server.com:8081/artifactory/api` | string   | TRUE     |
| token       | Identity Token for the Artifactory account executing the function calls         | string   | TRUE     |
| sourcePath  | Full file path where will be sourced from; **Needs proper escape chars          | string   | TRUE     |
| targetPath  | Target repo and folder destination of the artifact                              | string   | TRUE     |

#### Outputs
| Name         | Description                                                                   | Type     |
|--------------|-------------------------------------------------------------------------------|----------|
| downloadUri  | Download URI address used to retrieve/download the artifact from Artifactory  | string   |
| artifactUri  | Artifact URI address of the artifact and it's details within Artifactory      | string   |