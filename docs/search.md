# Search Functions

## GetArtifactsByProps
Searches for an artifact by one or more property names and optionally values, if provided (e.g. 'release' or 'release=stable'), and will return all artifacts that meet the search criteria.

- Multiple properties with no values will be separated by a required ampersand ('&'), handled by the function.
- Multiple properties with values will also be separated by a required ampersand ('&'), handled by the function, and should be passed into the function in the following format of 'propKey=propValue'.
The result will be something like 'release&channel' or 'release=stable&channel=windows-prod-iis'.

**Property keys/values are CASE SENSITIVE.**

#### Inputs
| Name        | Description                                           | Type     | Required |
|-------------|-------------------------------------------------------|----------|:--------:|
| listKvProps | List of one or more property key/values to seach for  | []string | TRUE     |

#### Outputs
| Name          | Description                                         | Type     |
|---------------|-----------------------------------------------------|----------|
| listArtifUris | Resulting list of matching artifacts by their URIs  | []string |


## GetArtifactsByName
Searches for artifacts by full or partial artifact name. **The search is CASE INSENSITIVE.**

#### Inputs
| Name       | Description                                         | Type     | Required |
|------------|-----------------------------------------------------|----------|:--------:|
| artifName  | Full or partial name of the artifact to search for  | string   | TRUE     |

#### Outputs
| Name          | Description                                         | Type     |
|---------------|-----------------------------------------------------|----------|
| listArtifUris | Resulting list of matching artifacts by their URIs  | []string |
| err           | nil unless error; then returns error                | error    |


## FilterListByFileType
Filters a list of artifact URIs by desired file type. If no extension is provided, the default filter will be VMware Templates (.vmtx). If file extension provided doesn't include a leading '.', it will be added.

This function would primarily be used in conjunction with `GetArtifactsByName` as part of the artifact filtering process. 

#### Inputs
| Name          | Description                           | Type      | Required |
|---------------|---------------------------------------|-----------|:--------:|
| ext           | File extension to filter by           | string    | TRUE     |
| listArtifUris | List of artifact URIs to be filtered  | []string  | TRUE     |

#### Outputs
| Name          | Description                                         | Type     |
|---------------|-----------------------------------------------------|----------|
| filteredList  | Resulting list of matching artifacts by their URIs  | []string |
| err           | nil unless error; then returns error                | error    |