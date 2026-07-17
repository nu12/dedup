package dedup

import (
	"fmt"
	"io"
	"path/filepath"
	"unique"
)

type Application struct {
	SourceFolder      string
	DestinationFolder string
	ListFlag          bool
	MoveFlag          bool

	OpenFunc   func(string) (io.ReadCloser, error)
	CopyFunc   func(io.Writer, io.Reader) (written int64, err error)
	CreateFunc func(string) (io.WriteCloser, error)
	RemoveFunc func(string) error

	counter     int
	uniqueFiles []unique.Handle[string]
	duplicates  []string
}

func (app *Application) Init() *Application {
	app.counter = 0
	app.duplicates = []string{}
	app.uniqueFiles = []unique.Handle[string]{}
	return app
}

func (app *Application) Run() error {
	if err := filepath.WalkDir(app.SourceFolder, app.WalkDir); err != nil {
		return err
	}
	app.List()
	if err := app.Move(); err != nil {
		return err
	}
	return nil
}

func (app *Application) List() {
	if app.ListFlag {
		fmt.Println()
		for _, dup := range app.duplicates {
			fmt.Println(dup)
		}
	}
}

func (app *Application) Move() error {
	if app.MoveFlag {
		fmt.Println()
		fmt.Printf("Moving duplicate files to %s\n", app.DestinationFolder)
		for _, dup := range app.duplicates {
			err := app.MoveFile(dup, filepath.Join(app.DestinationFolder, filepath.Base(dup)))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
