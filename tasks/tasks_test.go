package tasks

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/raynaluzier/artifactory-go-sdk/util"
)

func TestDownloadArtifactsPreservesURLCase(t *testing.T) {
	outputDir := t.TempDir()
	requestedPath := "/Repo/MixedCase.OVA"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != requestedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, requestedPath)
		}
		_, _ = io.WriteString(w, "artifact bytes")
	}))
	defer server.Close()

	previousServer, previousToken, previousOutputDir := util.ServerApi, util.Token, util.OutputDir
	util.ServerApi = server.URL + "/artifactory/api"
	util.Token = "test-token"
	util.OutputDir = outputDir
	t.Cleanup(func() {
		util.ServerApi, util.Token, util.OutputDir = previousServer, previousToken, previousOutputDir
	})

	if got := DownloadArtifacts(util.ServerApi, util.Token, server.URL+requestedPath, outputDir); got != "End of download process" {
		t.Fatalf("result = %q, want successful completion", got)
	}
	contents, err := os.ReadFile(filepath.Join(outputDir, "MixedCase", "MixedCase.OVA"))
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "artifact bytes" {
		t.Fatalf("downloaded contents = %q", contents)
	}
}

func TestDownloadArtifactsReportsPrimaryDownloadFailure(t *testing.T) {
	outputDir := t.TempDir()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "denied", http.StatusForbidden)
	}))
	defer server.Close()

	previousServer, previousToken, previousOutputDir := util.ServerApi, util.Token, util.OutputDir
	util.ServerApi = server.URL + "/artifactory/api"
	util.Token = "test-token"
	util.OutputDir = outputDir
	t.Cleanup(func() {
		util.ServerApi, util.Token, util.OutputDir = previousServer, previousToken, previousOutputDir
	})

	if got := DownloadArtifacts(util.ServerApi, util.Token, server.URL+"/Repo/MixedCase.OVA", outputDir); got != "File download failed" {
		t.Fatalf("result = %q, want download failure", got)
	}
}
