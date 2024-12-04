# Archived Search Functions
The following functions are currently not used but archived in case a need for them arises.

## GetArtifactsByNameRepo
Searches for artifacts by either full or partial artifact name and optionally the full or partial repo name, if known. Retruns a list of one or more artifacts that match the search criteria.

*If not passing in a repo name, then an empty string ("") should be passed in.

#### Inputs
| Name       | Description                                  | Type     | Required |
|------------|----------------------------------------------|----------|:--------:|
| artifName  | Full or partial artifact name to search for  | string   | TRUE     |
| repo       | Full or partial repo name to search for      | string   | *FALSE   |

#### Outputs
| Name          | Description                             | Type      |
|---------------|-----------------------------------------|-----------|
| listArtifUris | Resulting list of artifact URIs         | []string  |
| err           | nil unless error; then returns error    | error     | 


## GetArtifactVersions
Only available if the folder structure was setup with a 'Layout' (artifacts will have a value for Module ID present). Requires at least the group ID (top level folder; must be the FULL name), artifact name (must be the FULL name), and optionally the repo name.

This search function returns a list of versions for artifacts matching the search terms and only takes in the information for a single artifact. Search terms are CASE SENSITIVE.

*If not passing in a repo name, then an empty string ("") should be passed in.

#### Inputs
| Name       | Description                                                 | Type     | Required |
|------------|-------------------------------------------------------------|----------|:--------:|
| groupId    | Top-level folder in the artifact's path; must be FULL name  | string   | TRUE     |
| artifName  | Artifact name to search for; must be FULL name              | string   | TRUE     |
| repo       | Repo name to search within                                  | string   | *FALSE   |

#### Outputs
| Name          | Description                             | Type      |
|---------------|-----------------------------------------|-----------|
| listVersions  | Resulting list of artifact versions     | []string  |
| err           | nil unless error; then returns error    | error     | 


## GetArtifactLatestVersions
Only available if the folder structure was setup with a 'Layout' (artifacts will have a value for Module ID present). Requires at least the group ID (top level folder; must be the FULL name), artifact name (must be the FULL name), and optionally the repo name.

This search function returns the latest versions for an artifact matching the search terms and only takes in the information for a single artifact. Search terms are CASE SENSITIVE.

*If not passing in a repo name, then an empty string ("") should be passed in.

#### Inputs
| Name       | Description                                                 | Type     | Required |
|------------|-------------------------------------------------------------|----------|:--------:|
| groupId    | Top-level folder in the artifact's path; must be FULL name  | string   | TRUE     |
| artifName  | Artifact name to search for; must be FULL name              | string   | TRUE     |
| repo       | Repo name to search within                                  | string   | *FALSE   |

#### Outputs
| Name          | Description                                    | Type      |
|---------------|------------------------------------------------|-----------|
| latestVersion | Latest version available for a given artifact  | []string  |
| err           | nil unless error; then returns error           | error     | 