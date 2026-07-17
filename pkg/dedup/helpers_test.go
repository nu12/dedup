package dedup

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"
	"unique"
)

type FakeRead struct {
	Text string
}

func (fake *FakeRead) Read(p []byte) (n int, err error) {
	return len(fake.Text), nil
}

func (fake *FakeRead) Close() error {
	return nil
}

func TestHash(t *testing.T) {
	expected := "d41d8cd98f00b204e9800998ecf8427e"
	App := &Application{
		OpenFunc: func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc: func(w io.Writer, r io.Reader) (written int64, err error) { return 0, nil },
	}

	actual, err := App.Hash("")
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("Hash doesn't match. Actual: %s, expected: %s", actual, expected)
	}
}

func TestHashErrorOpen(t *testing.T) {

	App := &Application{
		OpenFunc: func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, errors.New("Test error") },
	}

	_, err := App.Hash("")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestHashErrorCopy(t *testing.T) {
	App := &Application{
		OpenFunc: func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc: func(w io.Writer, r io.Reader) (written int64, err error) { return 0, errors.New("Test error") },
	}

	_, err := App.Hash("")
	if err == nil {
		t.Error(err)
	}
}

func TestLogProgress(t *testing.T) {
	oldStdout := os.Stdout
	nullFile, _ := os.Open(os.DevNull)
	os.Stdout = nullFile

	App := &Application{}
	App.LogProgress(false)
	App.LogProgress(true)

	os.Stdout = oldStdout
	nullFile.Close()
}

type FakeDirEntry struct {
	FileInfo *FakeFileInfo
	Error    error
}

func (fake *FakeDirEntry) Name() string               { return "" }
func (fake *FakeDirEntry) IsDir() bool                { return fake.FileInfo.IsDir() }
func (fake *FakeDirEntry) Type() fs.FileMode          { return fake.FileInfo.Mode() }
func (fake *FakeDirEntry) Info() (fs.FileInfo, error) { return fake.FileInfo, fake.Error }

type FakeFileInfo struct {
	FileName    string
	FileSize    int64
	FileMode    fs.FileMode
	FileModTime time.Time
	FileDir     bool
	FileSys     any
}

func (fake *FakeFileInfo) Name() string       { return fake.FileName }
func (fake *FakeFileInfo) Size() int64        { return fake.FileSize }
func (fake *FakeFileInfo) Mode() fs.FileMode  { return fake.FileMode }
func (fake *FakeFileInfo) ModTime() time.Time { return fake.FileModTime }
func (fake *FakeFileInfo) IsDir() bool        { return fake.FileDir }
func (fake *FakeFileInfo) Sys() any           { return fake.FileSys }

func TestWalkDir(t *testing.T) {
	App := &Application{
		counter:     0,
		uniqueFiles: []unique.Handle[string]{},
		duplicates:  []string{},

		OpenFunc: func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc: func(w io.Writer, r io.Reader) (written int64, err error) { return 0, nil },
	}
	entry := &FakeDirEntry{
		FileInfo: &FakeFileInfo{},
		Error:    nil,
	}

	if err := App.WalkDir("", entry, nil); err != nil {
		t.Error(err)
	}
	if App.counter != 1 {
		t.Errorf("Expected counter to be %d, got %d", 1, App.counter)
	}
	if len(App.uniqueFiles) != 1 {
		t.Errorf("Expected uniqueFiles length to be %d, got %d", 1, len(App.uniqueFiles))
	}
	if len(App.duplicates) != 0 {
		t.Errorf("Expected duplicates length to be %d, got %d", 0, len(App.duplicates))
	}
}

func TestWalkDirDuplicate(t *testing.T) {
	App := &Application{
		counter:     1,
		uniqueFiles: []unique.Handle[string]{unique.Make("d41d8cd98f00b204e9800998ecf8427e")},
		duplicates:  []string{},

		OpenFunc: func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc: func(w io.Writer, r io.Reader) (written int64, err error) { return 0, nil },
	}
	entry := &FakeDirEntry{
		FileInfo: &FakeFileInfo{},
		Error:    nil,
	}

	if err := App.WalkDir("", entry, nil); err != nil {
		t.Error(err)
	}
	if App.counter != 2 {
		t.Errorf("Expected counter to be %d, got %d", 2, App.counter)
	}
	if len(App.uniqueFiles) != 1 {
		t.Errorf("Expected uniqueFiles length to be %d, got %d", 1, len(App.uniqueFiles))
	}
	if len(App.duplicates) != 1 {
		t.Errorf("Expected duplicates length to be %d, got %d", 1, len(App.duplicates))
	}
}

func TestWalkDirError(t *testing.T) {
	App := &Application{
		counter:     0,
		uniqueFiles: []unique.Handle[string]{},
		duplicates:  []string{},

		OpenFunc: func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc: func(w io.Writer, r io.Reader) (written int64, err error) { return 0, nil },
	}
	entry := &FakeDirEntry{
		FileInfo: &FakeFileInfo{},
		Error:    nil,
	}

	if err := App.WalkDir("", entry, errors.New("Test error")); err == nil {
		t.Error("Expected error")
	}
}

type FakeWriter struct{}

func (FakeWriter) Close() error                      { return nil }
func (FakeWriter) Write(p []byte) (n int, err error) { return len(p), nil }

func TestMoveFile(t *testing.T) {
	App := &Application{
		OpenFunc:   func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc:   func(w io.Writer, r io.Reader) (written int64, err error) { return 0, nil },
		CreateFunc: func(s string) (io.WriteCloser, error) { return &FakeWriter{}, nil },
		RemoveFunc: func(s string) error { return nil },
	}
	err := App.MoveFile("", "")
	if err != nil {
		t.Error(err)
	}

	App.RemoveFunc = func(s string) error { return errors.New("Test error") }
	err = App.MoveFile("", "")
	if err.Error() != "couldn't remove source file: Test error" {
		t.Errorf("Expected %v, got %v", "couldn't remove source file: Test error", err.Error())
	}

	App.CopyFunc = func(w io.Writer, r io.Reader) (written int64, err error) { return 0, errors.New("Test error") }
	err = App.MoveFile("", "")
	if err.Error() != "couldn't copy to dest from source: Test error" {
		t.Errorf("Expected %v, got %v", "couldn't copy to dest from source: Test error", err.Error())
	}

	App.CreateFunc = func(s string) (io.WriteCloser, error) { return &FakeWriter{}, errors.New("Test error") }
	err = App.MoveFile("", "")
	if err.Error() != "couldn't open dest file: Test error" {
		t.Errorf("Expected %v, got %v", "couldn't open dest file: Test error", err.Error())
	}

	App.OpenFunc = func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, errors.New("Test error") }
	err = App.MoveFile("", "")
	if err.Error() != "couldn't open source file: Test error" {
		t.Errorf("Expected %v, got %v", "couldn't open source file: Test error", err.Error())
	}
}
