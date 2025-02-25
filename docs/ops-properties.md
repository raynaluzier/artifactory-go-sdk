# Property-based Operations Functions

## GetArtifactPropValues
Takes in the URI of the artifact, plus one or more property keys, and returns the values for only the properties included in the URI for the given artifact. Meaning, the artifact can have more properties assigned to it, but those values will not be returned unless they were part of the request.

**Searches are CASE SENSITIVE.**

The list of properties sent over in the URI path must be separated by commas (',') which we handle before making the REST API call.

The result is whatever properties we passed and their values, which can be anything, so the returned JSON data is considered unstructured. Therefore, we capture each key `Name` and `Value` in a custom `prop` struct called `properties` that will be returned from the function of type interface{}.

#### Inputs
| Name          | Description                                              | Type      | Required |
|---------------|----------------------------------------------------------|-----------|:--------:|
| artifactUri   | URI of the artifact itself (different from Download URI) | string    | TRUE     |
| listPropKeys  | List of property keys we want the values of              | []string  | TRUE     |

#### Outputs
| Name       | Description                                                             | Type        |
|------------|-------------------------------------------------------------------------|-------------|
| properties | Resulting list of requested property key(s) and value(s) of type []prop | interface{} |
| err        | nil unless error; then returns error                                    | error       |


## GetAllPropsForArtifact
Takes in the URI of a given artifact and pulls all of the properties and their values assigned to the artifact (versus just select properties, as above).

The result is whatever properties exist and their values, which can be anything, so the returned JSON data is considered unstructured. Therefore, we capture each key `Name` and `Value` in a custom `prop` struct called `properties` that will be returned from the function of type interface{}.

**Searches are CASE SENSITIVE.**

#### Inputs
| Name          | Description                                              | Type      | Required |
|---------------|----------------------------------------------------------|-----------|:--------:|
| artifactUri   | URI of the artifact itself (different from Download URI) | string    | TRUE     |

#### Outputs
| Name       | Description                                                             | Type        |
|------------|-------------------------------------------------------------------------|-------------|
| properties | Resulting list of requested property key(s) and value(s) of type []prop | interface{} |
| err        | nil unless error; then returns error                                    | error       |


## FilterListByProps
Takes in a list of artifact URIs, and for each URI, it pulls the artifact's list of properties. Then the function compares the list of one or more key/values pairs ('key=value') provided as inputs against the key/values assigned to the artifact. If there's a match, the artifact URI will be added to the `filteredList` list. 

If more than one property key/value pair was input as a filter ('release=stable', 'testing=passed', 'channel=win-prod-iis'), an instance of the matching artifact will be added to the `foundList` for each matched property. So the same URI may be added to the list multiple times for multiple property matches.

For example: If 3 property key/value pairs were input as filters, we would expect that any artifact that has ALL of those matching properties is most likely the artifact we're looking for. However, it's probable that multiple artifacts have at least some of those same property key/value pairs (like, 'release=stable', 'testing=passed') for a given artifact (say, a new 'win-22' image built over multiple days).

- If only one artifact is present in the foundList, we'll return this artifact.
- If multiple artifacts are returned in the foundList, we will count the instances of each artifact in the list. 
- If the number of duplicate artifacts found matches the number of input property key/value pairs, they will be added to the `filteredList`. 
- If there's only one artifact in the filteredList, this will be returned. 
- If multiple artifacts are still present in the filteredList, the created date for each artifact will be grabbed and the latest artifact will be returned.

**Artifact URIs and Property key/values are CASE SENSITIVE.**

#### Inputs
| Name          | Description                               | Type      | Required |
|---------------|-------------------------------------------|-----------|:--------:|
| listArtifUris | List of artifact URIs                     | []string  | TRUE     |
| listKvProps   | List of key/value pairs input as filters  | []string  | TRUE     |

#### Outputs
| Name       | Description                           | Type    |
|------------|---------------------------------------|---------|
| foundItem  | Resulting artifact URI                | string  |
| err        | nil unless error; then returns error  | error   |


## SetArtifactProps
Takes in the URI of a given artifact and one or more property key/value pairs and assigns them to the given artifact. If more than one property key/value is supplied, they must be separated by a semi-colon 
(';'), which is handled before making the REST API call.

**Inputs are CASE SENSITIVE.**

Special characters are disallowed: 	)( }{ ][ *+^$\/~`!@#%&<>;, and the SPACE character

#### Inputs
| Name          | Description                                              | Type      | Required |
|---------------|----------------------------------------------------------|-----------|:--------:|
| artifactUri   | URI of the artifact itself (different from Download URI) | string    | TRUE     |
| listKvProps   | List of key/value pairs input as filters                 | []string  | TRUE     |

#### Outputs
| Name        | Description                                                           | Type     |
|-------------|-----------------------------------------------------------------------|----------|
| statusCode  | Resulting status code of the delete operation (either "204" or "404") | string   |
| err         | nil unless error; then returns error                                  | error    |


## DeleteArtifactProps
Takes in the URI of a given artifact and one or more property keys and removes them from the given artifact. If more than one property key is supplied, they must be separated by a comma (','), which is handled before making the REST API call.

**Inputs are CASE SENSITIVE.**

If a property is provided that doesn't exist (which includes incorrectly cased properties), the API ignores this and will return a successful response.

#### Inputs
| Name          | Description                                              | Type      | Required |
|---------------|----------------------------------------------------------|-----------|:--------:|
| artifactUri   | URI of the artifact itself (different from Download URI) | string    | TRUE     |
| listProps     | List of property keys to delete from artifact            | []string  | TRUE     |

#### Outputs
| Name        | Description                                                           | Type     |
|-------------|-----------------------------------------------------------------------|----------|
| statusCode  | Resulting status code of the delete operation (either "204" or "404") | string   |
| err         | nil unless error; then returns error                                  | error    |