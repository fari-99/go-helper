package wsstorage

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func localConfig(t *testing.T) Config {
	t.Helper()
	return Config{Driver: DriverLocal, LocalPath: t.TempDir()}
}

func TestUploadLocalStoresContentAndSniffsMime(t *testing.T) {
	cfg := localConfig(t)
	pngHeader := []byte("\x89PNG\r\n\x1a\n" + strings.Repeat("x", 100))

	result, err := Upload(context.Background(), cfg, "avatars", "me.PNG", bytes.NewReader(pngHeader), 1024)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	if result.Mime != "image/png" {
		t.Errorf("mime = %q, want image/png (sniffed, not client supplied)", result.Mime)
	}
	if !strings.HasSuffix(result.Filename, ".png") {
		t.Errorf("filename = %q, want .png extension", result.Filename)
	}
	if result.Size != int64(len(pngHeader)) {
		t.Errorf("size = %d, want %d", result.Size, len(pngHeader))
	}

	stored, err := os.ReadFile(filepath.Join(cfg.LocalPath, "avatars", result.Path, result.Filename))
	if err != nil {
		t.Fatalf("stored file missing: %v", err)
	}
	if !bytes.Equal(stored, pngHeader) {
		t.Error("stored content differs from input (images must be stored as-is)")
	}
}

func TestUploadRejectsOverLimitAndLeavesNothing(t *testing.T) {
	cfg := localConfig(t)

	_, err := Upload(context.Background(), cfg, "docs", "big.txt", strings.NewReader(strings.Repeat("a", 2000)), 1000)
	if !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("err = %v, want ErrFileTooLarge", err)
	}

	assertNoFiles(t, cfg.LocalPath)
}

func TestUploadExactlyAtLimitIsAccepted(t *testing.T) {
	cfg := localConfig(t)

	if _, err := Upload(context.Background(), cfg, "docs", "ok.txt", strings.NewReader(strings.Repeat("a", 1000)), 1000); err != nil {
		t.Fatalf("upload at exact limit failed: %v", err)
	}
}

func TestUploadRejectsEmptyAndBadType(t *testing.T) {
	cfg := localConfig(t)

	if _, err := Upload(context.Background(), cfg, "docs", "e.txt", strings.NewReader(""), 10); !errors.Is(err, ErrEmptyFile) {
		t.Errorf("empty: err = %v, want ErrEmptyFile", err)
	}

	for _, bad := range []string{"", "../etc", "a/b", "a b", strings.Repeat("a", 65)} {
		if _, err := Upload(context.Background(), cfg, bad, "f.txt", strings.NewReader("x"), 10); !errors.Is(err, ErrInvalidType) {
			t.Errorf("file type %q: err = %v, want ErrInvalidType", bad, err)
		}
	}
}

func TestUploadCancelledContextRemovesPartialFile(t *testing.T) {
	cfg := localConfig(t)
	ctx, cancel := context.WithCancel(context.Background())

	reader := &cancelAfterFirstRead{cancel: cancel}
	if _, err := Upload(ctx, cfg, "docs", "f.txt", reader, 0); err == nil {
		t.Fatal("expected an error from the cancelled upload")
	}

	assertNoFiles(t, cfg.LocalPath)
}

type cancelAfterFirstRead struct {
	cancel func()
	n      int
}

func (c *cancelAfterFirstRead) Read(p []byte) (int, error) {
	c.n++
	if c.n == 1 {
		return copy(p, "hello"), nil
	}
	c.cancel()
	return 0, context.Canceled
}

func TestSafeExtension(t *testing.T) {
	cases := map[string]string{
		"a.JPG":                 ".jpg",
		"a":                     "",
		"a.tar.gz":              ".gz",
		"../../a.txt":           ".txt",
		`c:\x\evil.exe`:         ".exe",
		"a.":                    "",
		"a.bad ext":             "",
		"a.waytoolongextension": "",
	}

	for input, want := range cases {
		if got := safeExtension(input); got != want {
			t.Errorf("safeExtension(%q) = %q, want %q", input, got, want)
		}
	}
}

func assertNoFiles(t *testing.T, root string) {
	t.Helper()
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			t.Errorf("unexpected file left behind: %s", path)
		}
		return nil
	})
}
