package main

import (
	"fmt"
	"os"

	_ "github.com/raynaluzier/artifactory-go-sdk/common"
	_ "github.com/raynaluzier/artifactory-go-sdk/operations"
	_ "github.com/raynaluzier/artifactory-go-sdk/search"
	"github.com/raynaluzier/artifactory-go-sdk/tasks"
	"github.com/raynaluzier/artifactory-go-sdk/util"
)

func main(){

	// -------------- TESTING --------------------------------------------------
	// mimicking passing in vars and then assigning them to the global vars
	serverApi 	:= os.Getenv("ARTIFACTORY_SERVER")
	token 		:= os.Getenv("ARTIFACTORY_TOKEN")
	logLevel 	:= os.Getenv("ARTIFACTORY_LOGGING")
	//outputDir 	:= os.Getenv("ARTIFACTORY_OUTPUTDIR")   //c:\lab\output-testing\ or /lab/output-testing

	util.ServerApi = serverApi
	util.Token     = token
	util.Logging   = logLevel
	//util.OutputDir = outputDir
	// -------------------------------------------------------------------------
	outputDir := "c:\\lab"
	//downloadUri := "https://riverpointtechnology.jfrog.io/artifactory/rpt-libs-local/image9012/image9012.ova"
	//downloadUri := "https://riverpointtechnology.jfrog.io/artifactory/rpt-libs-local/image5678/image5678.ovf"
	downloadUri := "https://riverpointtechnology.jfrog.io/artifactory/rpt-libs-local/image1234/image1234.vmtx"
	//fmt.Println(operations.GetArtifact(downloadUri))

	fmt.Println(tasks.DownloadArtifacts(serverApi, token, downloadUri, outputDir))

	//imageType := "vmtx"
	//imageName := "image1234"
	//sourceDir := "c:\\lab\\" + imageName
	//targetDir := "/rpt-libs-local"
	//fileSuffix := ""
	//fmt.Println(tasks.UploadArtifacts(serverApi, token, imageType, imageName, sourceDir, targetDir, fileSuffix))

}
