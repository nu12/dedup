package dedup

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"slices"
	"unique"
)

func (app *Application) Hash(path string) (string, error) {
	file, err := app.OpenFunc(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := app.CopyFunc(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (app *Application) LogProgress(clear bool) {
	if clear {
		fmt.Printf("\r")
	}
	fmt.Printf("%d files processed (%d duplicates found)", app.counter, len(app.duplicates))
}

func (app *Application) WalkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	app.counter++

	if !d.IsDir() {
		h, _ := app.Hash(path)
		file := unique.Make(h)
		if slices.Contains(app.uniqueFiles, file) {
			app.duplicates = append(app.duplicates, path)
		} else {
			app.uniqueFiles = append(app.uniqueFiles, file)
		}
	}
	app.LogProgress(true)
	return nil
}

func (app *Application) MoveFile(sourcePath, destPath string) error {
	inputFile, err := app.OpenFunc(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := app.CreateFunc(destPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = app.CopyFunc(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = app.RemoveFunc(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't remove source file: %v", err)
	}
	return nil
}
