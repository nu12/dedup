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

func (this *Application) Init() *Application {
	this.counter = 0
	this.duplicates = []string{}
	this.uniqueFiles = []unique.Handle[string]{}
	return this
}

func (this *Application) Run() error {
	if err := filepath.WalkDir(this.SourceFolder, this.WalkDir); err != nil {
		return err
	}
	this.List()
	if err := this.Move(); err != nil {
		return err
	}
	return nil
}

func (this *Application) List() {
	if this.ListFlag {
		fmt.Println()
		for _, dup := range this.duplicates {
			fmt.Println(dup)
		}
	}
}

func (this *Application) Move() error {
	if this.MoveFlag {
		fmt.Println()
		fmt.Printf("Moving duplicate files to %s\n", this.DestinationFolder)
		for _, dup := range this.duplicates {
			err := this.MoveFile(dup, filepath.Join(this.DestinationFolder, filepath.Base(dup)))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
