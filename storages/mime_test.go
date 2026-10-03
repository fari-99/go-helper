package storages

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func headerFor(t *testing.T, path string) *multipart.FileHeader {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, _ := w.CreateFormFile("f", filepath.Base(path))
	part.Write(data)
	w.Close()
	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if err := req.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	return req.MultipartForm.File["f"][0]
}

func TestGetFileDataMedia(t *testing.T) {
	dir := os.Getenv("MIME_SAMPLES")
	if dir == "" {
		t.Skip("MIME_SAMPLES not set")
	}
	cases := map[string]string{"a.mp3": "audio/mpeg", "real.mp3": "audio/mpeg", "a.m4a": "audio/", "a.flac": "audio/", "a.opus": "audio/", "a.wav": "audio/", "a.aac": "audio/"}
	base := &StorageBase{}
	for name, want := range cases {
		data, err := base.getFileData(headerFor(t, filepath.Join(dir, name)))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		t.Logf("%s -> %s", name, data.ContentType)
		if len(data.ContentType) < len(want) || data.ContentType[:len(want)] != want {
			t.Errorf("%s: got %s want %s*", name, data.ContentType, want)
		}
	}
}
