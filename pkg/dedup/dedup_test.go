package dedup

import (
	"errors"
	"io"
	"testing"
	"unique"
)

func TestInit(t *testing.T) {
	App := &Application{}

	App.Init()
	if App.counter != 0 {
		t.Errorf("Expected counter to be %d, got %d", 0, App.counter)
	}

	if len(App.uniqueFiles) != 0 {
		t.Errorf("Expected uniqueFiles length to be %d, got %d", 0, len(App.uniqueFiles))
	}

	if len(App.duplicates) != 0 {
		t.Errorf("Expected duplicates length to be %d, got %d", 0, len(App.duplicates))
	}
}

func TestRun(t *testing.T) {
	App := &Application{
		counter:           0,
		uniqueFiles:       []unique.Handle[string]{},
		duplicates:        []string{""},
		ListFlag:          false,
		MoveFlag:          false,
		OpenFunc:          func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc:          func(w io.Writer, r io.Reader) (written int64, err error) { return int64(len([]byte("a"))), nil },
		CreateFunc:        func(s string) (io.WriteCloser, error) { return &FakeWriter{}, nil },
		RemoveFunc:        func(s string) error { return nil },
		SourceFolder:      ".",
		DestinationFolder: ".",
	}
	if err := App.Run(); err != nil {
		t.Error(err)
	}
}

func TestRunErrorOnMove(t *testing.T) {
	App := &Application{
		MoveFlag:     true,
		OpenFunc:     func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc:     func(w io.Writer, r io.Reader) (written int64, err error) { return int64(len([]byte("a"))), nil },
		CreateFunc:   func(s string) (io.WriteCloser, error) { return &FakeWriter{}, nil },
		RemoveFunc:   func(s string) error { return errors.New("Test error") },
		SourceFolder: ".",
	}

	if err := App.Run(); err == nil {
		t.Error("Expected error")
	}
}

func TestRunErrorSourceFolder(t *testing.T) {
	App := &Application{
		counter:           0,
		uniqueFiles:       []unique.Handle[string]{},
		duplicates:        []string{""},
		ListFlag:          false,
		MoveFlag:          false,
		OpenFunc:          func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc:          func(w io.Writer, r io.Reader) (written int64, err error) { return int64(len([]byte("a"))), nil },
		CreateFunc:        func(s string) (io.WriteCloser, error) { return &FakeWriter{}, nil },
		RemoveFunc:        func(s string) error { return nil },
		SourceFolder:      "",
		DestinationFolder: ".",
	}
	App.SourceFolder = ""
	if err := App.Run(); err == nil {
		t.Errorf("Expected error when SourceFolder is not a valid path")
	}

}

func TestList(t *testing.T) {
	App := &Application{
		ListFlag:   true,
		duplicates: []string{""},
	}
	App.List()
}

func TestMove(t *testing.T) {
	App := &Application{
		duplicates: []string{""},
		MoveFlag:   true,
		OpenFunc:   func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc:   func(w io.Writer, r io.Reader) (written int64, err error) { return 0, nil },
		CreateFunc: func(s string) (io.WriteCloser, error) { return &FakeWriter{}, nil },
		RemoveFunc: func(s string) error { return nil },
	}

	if err := App.Move(); err != nil {
		t.Error(err)
	}
}

func TestMoveError(t *testing.T) {
	App := &Application{
		duplicates: []string{""},
		MoveFlag:   true,
		OpenFunc:   func(path string) (io.ReadCloser, error) { return &FakeRead{Text: "a"}, nil },
		CopyFunc:   func(w io.Writer, r io.Reader) (written int64, err error) { return 0, nil },
		CreateFunc: func(s string) (io.WriteCloser, error) { return &FakeWriter{}, nil },
		RemoveFunc: func(s string) error { return errors.New("Test error") },
	}

	if err := App.Move(); err == nil {
		t.Errorf("Expected error")
	}
}
