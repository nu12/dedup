package dedup

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"slices"
	"unique"
)

func (this *Application) Hash(path string) (string, error) {
	file, err := this.OpenFunc(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := this.CopyFunc(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (this *Application) LogProgress(clear bool) {
	if clear {
		fmt.Printf("\r")
	}
	fmt.Printf("%d files processed (%d duplicates found)", this.counter, len(this.duplicates))
}

func (this *Application) WalkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	this.counter++

	if !d.IsDir() {
		h, _ := this.Hash(path)
		file := unique.Make(h)
		if slices.Contains(this.uniqueFiles, file) {
			this.duplicates = append(this.duplicates, path)
		} else {
			this.uniqueFiles = append(this.uniqueFiles, file)
		}
	}
	this.LogProgress(true)
	return nil
}

func (this *Application) MoveFile(sourcePath, destPath string) error {
	inputFile, err := this.OpenFunc(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := this.CreateFunc(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = this.CopyFunc(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = this.RemoveFunc(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't remove source file: %v", err)
	}
	return nil
}
