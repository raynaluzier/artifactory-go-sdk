package common

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/raynaluzier/go-artifactory/util"
)

var logLevel slog.Level

func SetBearer(token string) string {
	bearer := "Bearer " + token
	return bearer
}

func CheckOsPlatform() string {
	// Detects the operating system this program is running on
	os := runtime.GOOS
	return os
}

func ConvertToLowercase(inputStr string) string {
	lowerStr := strings.ToLower(inputStr)
	return lowerStr
}

func ConvertToUppercase(inputStr string) string {
	upperStr := strings.ToUpper(inputStr)
	return upperStr
}

func RemoveDuplicateStrings(listOfStrings []string) ([]string) {
	allStrings := make(map[string]bool)

	list := []string{}
	for _, item := range listOfStrings {
		if _, value := allStrings[item]; !value {
			allStrings[item] = true
			list = append(list, item)
		}
	}
	return list
}

func ReturnWithDupCounts(listOfStrings []string) (map[string]int) {
	countMap := make(map[string]int)
	
	for _, str := range listOfStrings {
		countMap[str]++
	}
	return countMap
}

func ReturnDuplicates(countMap map[string]int) []string {
	// Takes in a count map of strings and their number of duplicate occurances (map[str1:1, str2:5, str3:1])
	// For any strings with more than one occurance, the string is added to the duplicates list and returned
	duplicates := []string{}

	for str, count := range countMap {
		if count > 1 {
			duplicates = append(duplicates, str)
		}
	}
	return duplicates
}

func SetArtifUriFromDownloadUri(downloadUri string) string {
	downloadUri = strings.Replace(downloadUri, "8082", "8081", 1)     // Modify the server port from 8082 to 8081
	artifSuffix := strings.TrimPrefix(downloadUri, util.ServerApi)    // /repo-key/folder/path/artifact.ext
	artifUri := util.ServerApi + "/storage" + artifSuffix             // http://server.com:8081/artifactory/api/storage/repo-key/folder/path/artifact.ext
	
	return artifUri
}

func SearchForExactString(searchTerm, inputStr string) (bool, error) {
	// For example: "win2022" will return true if input string is "win2022", false if "win2022-iis"
	result, err := regexp.MatchString("(?sm)^" + searchTerm + "$", inputStr)
	if err != nil {
		fmt.Println("Error searching for : " + searchTerm)
		return result, err
	}
	return result, err
}

func EscapeSpecialChars(input string) (string) {
	// Takes the output directory provided from the environment variable and adds escape characters
	// For Ex: F:\mypath\ becomes F:\\mypath\\
	var js json.RawMessage
	// Replace newlines with space rather than escaping them
	input = strings.ReplaceAll(input, "\n", " ")
	// Done to take the help of the json.Unmarshal function
	jsonString := createJsonString(input)
	byteValue := []byte(jsonString)
	err := json.Unmarshal(byteValue, &js)

	// Escape special characters only if JSON unmarshal results in an error
	if err != nil {
		out, err := json.Marshal(input)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			LogTxtHandler().Error("JSON marshalling failed with an error. " + strErr)
			return input
		} else {
			// JSON marshal quotes the entire string which results in double quotes at beginning/end of string
			return string(out[1 : len(out)-1])
		}
	}
	return input
}

func createJsonString(input string) string {
	// Used with EscapeSpecialChars function to properly format output directories that may include "\" in path
	jsonString := "{\"key\":\""
	endJson := "\"}"
	jsonString = jsonString + input + endJson
	return jsonString
}

func CheckPathType(path string) bool {
	// Checks path to see if path is Unix-based (has '/') or Windows-based (has '\')
	isWinPath := strings.Contains(path, "\\")
	return isWinPath
}

func StringCompare(inputStr, actualStr string) bool {
	// Performs case INSENSITIVE comparision of strings (like file names)
	// Does NOT do partial string comparisons
	if strings.EqualFold(inputStr, actualStr) {
		return true
	} else {   // Different strings
		return false
	}
}

func CheckAddSlashToPath(path string) string {
	lastChar := path[len(path)-1:]
	winPath := CheckPathType(path)

	if winPath == true {
		if lastChar == "\\" {
			LogTxtHandler().Debug("Path: '" + path + "' is formatted properly")
			return path
		} else {
			// Add backslash to path
			path = path + "\\"
			return path
		}
	} else {  // Unix Path
		if lastChar == "/" {
			LogTxtHandler().Debug("Path: '" + path + "' is formatted properly")
			return path
		} else {
			// Add forwardslash to path
			path = path + "/"
			return path
		}
	}
}

func ContainsSpecialChars(strings []string) bool {
    // Checks for the special characters disallowed by Artifactory in Properties
	// Returns true if ANY of the chars are found; false if not
	pattern := regexp.MustCompile("[(){}\\[\\]*+^$\\/~`!@#%&<>;, ]")  // add '=' back later
	for idx := 0; idx< len(strings); idx++ {
		if pattern.MatchString(strings[idx]) {
			return true
		}
	}
	return false
}

func SetLoggingLevel() slog.Level {
	level := util.Logging

	switch level {
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	case "DEBUG":
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}
	return logLevel
}

func LogTxtHandler() *slog.Logger {
	loggingLevel := SetLoggingLevel()
	opts := &slog.HandlerOptions{
		Level: slog.Level(loggingLevel),
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	txtLogger := slog.New(handler)
	return txtLogger
}

func LogJsonHandler() *slog.Logger {
	loggingLevel := SetLoggingLevel()
	opts := &slog.HandlerOptions{
		Level: slog.Level(loggingLevel),
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	jsonLogger := slog.New(handler)
	return jsonLogger
}