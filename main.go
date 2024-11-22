package main

import (
	_ "artifactory/common"
	"artifactory/operations"
	"artifactory/search"
	"fmt"
)



func main(){

	//prop := "channel"
	//val  := "windows-prod"
	//prop := "release"
	//val := "latest-stable"
	//val := "stable"
	//val := ""
	//artifName := "vmxt"
	//artifName := "W22_X64_STD_24_09_02"
	//artifName := "w22_X64_STD"
	//artifName := "win-22-882tv73c2001482aasxvn908"
	//artifName := "banana"
	//artifName := "test-artifact"
	//artifName := "WIN-22-882tv73c2001482aasxvn908.vmxt"

	//listProps := []string{}
	//listProps = append(listProps, "channel")

	//artifPath := "api-testing/test-artifact"
	//downloadUri := "http://server.com:8082/artifactory/repo-key/folder/win-22-882tv73c2001482aasxvn908.ova"

	//sourcePath := "H:\\repos\\artifactory-go\\output-testing\\another3.txt"
	//sourcePath := ""
	//targetPath := "/repo/api-testing/another"
	//targetPath := ""

	//artifUri := "http://server.com:8082/artifactory/repo-key/folder/win-22-882tv73c2001482aasxvn908.vmxt"

	
	kvProps := []string{}
	kvProps = append(kvProps, "release=stable")
	
	//ext := "vmxt"
	ext := "txt"
	//artifName := "win-22"
	artifName := "W22"
	listArtifacts, err := search.GetArtifactsByName(artifName)
	if err != nil {
		fmt.Println("some error")
	}

	newArtifacts, err := (search.FilterListByFileType(ext, listArtifacts))
	if err != nil {
		fmt.Println("some other error")
	}

	fmt.Println(operations.FilterListByProps(newArtifacts, kvProps))


}
