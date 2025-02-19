# Common Functions

## SetBearer
Takes in the Artifactory account Identity Token and forms the bearer token to be used in subsequent REST API calls.

#### Inputs
| Name       | Description                                                                              | Type     | Required |
|------------|------------------------------------------------------------------------------------------|----------|:--------:|
| token      | Identity Token for the Artifactory account executing the function calls                  | string   | TRUE     |

#### Outputs
| Name        | Description                                                                              | Type     |
|-------------|------------------------------------------------------------------------------------------|----------|
| bearer      | Forms bearer token to be passed with REST API Call to Artifactory                        | string   |


## CheckOsPlatform
Detects the operating system this program is running on and will return `windows`, `linux`, or `darwin` (MAC).

#### Inputs
Takes no inputs.

#### Outputs
| Name   | Description              | Type     |
|--------|--------------------------|----------|
| os     | Resulting OS platform    | string   |


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


## ParseArtifUriForPath
Takes in an Artifact URI, trims off the filename and server API / storage path and returns the '/repo/folder/path/'. Within the Artifactory post-processor plugin, this can be used in place of a target path value where the location of an existing artifact is used as the target location for a new artifact.

#### Inputs
| Name        | Description                                                           | Type   | Required |
|-------------|-----------------------------------------------------------------------|--------|:--------:|
| serverApi   | Server API address of the Artifactory server instance                 | string | TRUE     |
| artifactUri | Artifact URI address of an existing image/artifact within Artifactory | string | TRUE     |

#### Outputs
| Name         | Description                                                            | Type   |
|--------------|------------------------------------------------------------------------|--------|
| artifactPath | Returns the '/repo/folder/path/' parsed from the provided artifact URI | string |


## ParseArtifUriForFilename
Takes in an Artifact URI and determines the associated filename.

#### Inputs
| Name        | Description                                                           | Type   | Required |
|-------------|-----------------------------------------------------------------------|--------|:--------:|
| artifactUri | Artifact URI address of an existing image/artifact within Artifactory | string | TRUE     |

#### Outputs
| Name      | Description                                                            | Type   |
|-----------|------------------------------------------------------------------------|--------|
| fileName  | Returns the 'filename.ext' parsed from the provided artifact URI       | string |


## ParseFilenameForImageName
Takes in the artifact's filename (filename.ext) and trims off the extension and returns just the name of the image without the file extension.

#### Inputs
| Name      | Description                                       | Type   | Required |
|-----------|---------------------------------------------------|--------|:--------:|
| fileName  | File name of the artifact with file extension     | string | TRUE     |

#### Outputs
| Name      | Description                    | Type   |
|-----------|--------------------------------|--------|
| imageName | Returns the name of the image  | string |


## SetLoggingLevel
Uses the Global Variable `util.Logging` to set the desired logging level and returns the slog.Level equivalent value to be used by the desired logging handlers (LogTxtHandler or LogJsonHandler). If not specified, logging level defaults to INFO.

#### Inputs
Takes no inputs

#### Outputs
| Name      | Description                                                                                | Type       |
|-----------|--------------------------------------------------------------------------------------------|------------|
| logLevel  | Will be slog.LevelInfo, slog.LevelWarn, slog.LevelError, or slog.LevelDebug based on input | slog.Level |


## LogTxtHandler
Takes in the appropriate logging level type from `SetLoggingLevel()` and sets the level in the handler options. Then a new Text handler interface is created with the specified logging format and defines where they are written to, in this case, Stdout.

Output example:  `time=2024-12-02T10:35:41.267-07:00 level=INFO msg="This is your info message."`

#### Inputs
| Name        | Description                                             | Type     | Required |
|-------------|---------------------------------------------------------|----------|:--------:|
| LOGGING     | Desired log level; Accepts: INFO, WARN, ERROR, DEBUG    | string   | FALSE    |

#### Outputs
| Name      | Description                                                                                | Type       |
|-----------|--------------------------------------------------------------------------------------------|------------|
| logLevel  | Will be slog.LevelInfo, slog.LevelWarn, slog.LevelError, or slog.LevelDebug based on input | slog.Level |

#### Usage
Example: `someLogLevel := common.SetLoggingLevel()`
         `common.LogTxtHandler(someLogLevel).Info("Info stuff. All is well!")`
         `common.LogTxthandler(someLogLevel).Debug("Found object: test-artifact.txt")`

If 'INFO' is set in .env, then only the .Info, .Warn, and .Error logs will be output.
If 'WARN' is set, then only .Warn, and .Error logs will be output.
If 'ERROR' is set, then only .Error logs will be output.
If 'DEBUG' is set, then all logs - .Info, .Warn, .Error, and .Debug - will be output.


## LogJsonHandler
Takes in the appropriate logging level type from `SetLoggingLevel()` and sets the level in the handler options. Then a new JSON handler interface is created with the specified logging format and defines where they are written to, in this case, Stdout. The JSON handler is useful for parsing and performing other actions based on the output, or writing to an external logging system.

Output example:  `{"time":"2024-12-02T10:13:31.252815-07:00","level":"INFO","msg":"Some JSON Info message."}`

#### Inputs
| Name        | Description                                             | Type     | Required |
|-------------|---------------------------------------------------------|----------|:--------:|
| LOGGING     | Desired log level; Accepts: INFO, WARN, ERROR, DEBUG    | string   | FALSE    |

#### Outputs
| Name      | Description                                                                                | Type       |
|-----------|--------------------------------------------------------------------------------------------|------------|
| logLevel  | Will be slog.LevelInfo, slog.LevelWarn, slog.LevelError, or slog.LevelDebug based on input | slog.Level |

#### Usage
Example: `someLogLevel := common.SetLoggingLevel()`
         `common.LogJsonHandler(someLogLevel).Info("Info stuff. All is well!")`
         `common.LogJsonhandler(someLogLevel).Debug("Found object: test-artifact.txt")`

If 'INFO' is set in .env, then only the .Info, .Warn, and .Error logs will be output.
If 'WARN' is set, then only .Warn, and .Error logs will be output.
If 'ERROR' is set, then only .Error logs will be output.
If 'DEBUG' is set, then all logs - .Info, .Warn, .Error, and .Debug - will be output.


## TrimEndSlashUrl
Takes in a URL string and trims off the ending slash if it exists. If the URL doesn't end with a slash, then the original URL is returned.

#### Inputs
| Name | Description | Type     | Required |
|------|-------------|----------|:--------:|
| url  | URL to trim | string   | TRUE     |

#### Outputs
| Name       | Description               | Type    |
|------------|---------------------------|---------|
| trimmedUrl | URL with no ending slash  | string  |


## CreateTestDirectory
Used as part of the Artifactory plugin acceptance testing, creates a test directory (the plugin will name it "test-directory") in the user's home directory. As part of this, it will update the directory path to include an ending slash (forward or back, depending on the OS platform) and sets the directory permissions to 0755.

**This will be removed once acceptance testing within the plugin is complete.**

#### Inputs
| Name     | Description                              | Type     | Required |
|----------|------------------------------------------|----------|:--------:|
| dirName  | Name of the test directory to be created | string   | TRUE     |

#### Outputs
| Name       | Description                                                    | Type    |
|------------|----------------------------------------------------------------|---------|
| newDirPath | Returns the directory name created with ending slash appended  | string  |


## RenameFile
Used as part of the Artifactory plugin acceptance testing, takes in the full path to the test file and renames it to the new file path.

#### Inputs
| Name         | Description                                   | Type     | Required |
|--------------|-----------------------------------------------|----------|:--------:|
| oldFilePath  | Full path to the file that will be renamed    | string   | TRUE     |
| newFilePath  | Full path to the resulting file after rename  | string   | TRUE     |

#### Outputs
| Name     | Description                          | Type    |
|----------|--------------------------------------|---------|
| (result) | Returns either "Failed" or "Success" | string  |


## CreateTestFile
Used as part of the Artifactory plugin acceptance testing, this takes in the test directory path (for plugin acceptance testing, this will be "$HOME_DIR/test-directory/", formatted per OS-platform), name of the file, and file contents to create a new test file called "test-artifact.txt" (passed in by the plugin), write brief file contents, and then will using the `RenameTestFile` function to rename it to "test-artifact.ova" to follow the supported image types of the plugin.

**This will be removed once acceptance testing within the plugin is complete.**

#### Inputs
| Name         | Description                                   | Type     | Required |
|--------------|-----------------------------------------------|----------|:--------:|
| dirPath      | Directory name that will house the test file  | string   | TRUE     |
| fileName     | Name of the test file that will be created    | string   | TRUE     |
| fileContents | Short amount of text to write to the file     | string   | TRUE     |

#### Outputs
| Name        | Description                                                                     | Type    |
|-------------|---------------------------------------------------------------------------------|---------|
| newFilePath | The full path to the test file (ex: $HOME_DIR/test-directory/test-artifact.ova) | string  |


## CreateTestRepo
Used as part of the Artifactory plugin acceptance testing, this creates a test repo called "test-packer-plugin" within the defined Artifactory environment.

**This will be removed once acceptance testing within the plugin is complete.**

#### Inputs
None

#### Outputs
| Name         | Description                                      | Type    |
|--------------|--------------------------------------------------|---------|
| testRepoPath | Path to the test repo (i.e. /test-packer-plugin) | string  |


## DeleteTestRepo
Used as part of the Artifactory plugin acceptance testing, this deletes test repo called "test-packer-plugin" within the defined Artifactory environment once acceptance testing is complete. Any test artifacts that are in this test repo are automatically removed as well.

#### Inputs
None

#### Outputs
| Name         | Description                                         | Type    |
|--------------|-----------------------------------------------------|---------|
| statusCode | Returns "200" if deletion is successful, "400" if not | string  |


## DeleteTestFile
Used as part of the Artifactory plugin acceptance testing, this deletes test file called "test-artifact.ova" from the previously created test directory $HOME_DIR/test-directory (formatted per OS-platform) once acceptance testing is complete.

#### Inputs
| Name         | Description                               | Type     | Required |
|--------------|-------------------------------------------|----------|:--------:|
| dirPath      | Directory name that houses the test file  | string   | TRUE     |

#### Outputs
None


## DeleteTestDirectory
Used as part of the Artifactory plugin acceptance testing, this deletes test directory called "test-directory" from the user's $HOME_DIR (formatted per OS-platform) once acceptance testing is complete.

#### Inputs
| Name         | Description                               | Type     | Required |
|--------------|-------------------------------------------|----------|:--------:|
| dirPath      | Directory name that housed the test file  | string   | TRUE     |

#### Outputs
None