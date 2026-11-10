package operations

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/raynaluzier/artifactory-go-sdk/util"
)

func TestUploadFileStreamsSourceAndUsesActualFilename(t *testing.T) {
	tempDir := t.TempDir()
	actualName := "Image.OVA"
	contents := "artifact bytes"
	if err := os.WriteFile(filepath.Join(tempDir, actualName), []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/artifactory/repo/images/"+actualName {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("authorization = %q", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
		}
		if string(body) != contents {
			t.Errorf("body = %q, want %q", body, contents)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{"downloadUri":"https://artifactory.example/repo/images/Image.OVA"}`)
	}))
	defer server.Close()

	previousServer, previousToken := util.ServerApi, util.Token
	util.ServerApi = server.URL + "/artifactory/api"
	util.Token = "test-token"
	t.Cleanup(func() {
		util.ServerApi, util.Token = previousServer, previousToken
	})

	got, err := UploadFile(filepath.Join(tempDir, "image.ova"), "/repo/images/")
	if err != nil {
		t.Fatal(err)
	}
	if want := "https://artifactory.example/repo/images/Image.OVA"; got != want {
		t.Fatalf("download URI = %q, want %q", got, want)
	}
}

func TestUploadFileReturnsErrorForRejectedUpload(t *testing.T) {
	tempDir := t.TempDir()
	source := filepath.Join(tempDir, "image.ova")
	if err := os.WriteFile(source, []byte("artifact bytes"), 0600); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "denied", http.StatusForbidden)
	}))
	defer server.Close()

	previousServer, previousToken := util.ServerApi, util.Token
	util.ServerApi = server.URL + "/artifactory/api"
	util.Token = "test-token"
	t.Cleanup(func() {
		util.ServerApi, util.Token = previousServer, previousToken
	})

	if _, err := UploadFile(source, "/repo/images/"); err == nil {
		t.Fatal("expected rejected upload to return an error")
	}
}

func TestRetrieveArtifactPreservesCaseAndRejectsHTTPFailure(t *testing.T) {
	outputDir := t.TempDir()
	requestedPath := "/Repo/MixedCase.OVA"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != requestedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, requestedPath)
		}
		if r.URL.Path == requestedPath {
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, "denied")
		}
	}))
	defer server.Close()

	previousOutputDir, previousToken := util.OutputDir, util.Token
	util.OutputDir = outputDir
	util.Token = "test-token"
	t.Cleanup(func() {
		util.OutputDir, util.Token = previousOutputDir, previousToken
	})

	if _, err := RetrieveArtifact(server.URL + requestedPath); err == nil {
		t.Fatal("expected non-success HTTP status to return an error")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "MixedCase.OVA")); !os.IsNotExist(err) {
		t.Fatalf("failed response should not create an artifact file, stat error = %v", err)
	}
}
