/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"unique"

	"github.com/spf13/cobra"
)

func hashFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't remove source file: %v", err)
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dedup",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		uniqueFiles := []unique.Handle[string]{}
		counter := 0
		duplicates := []string{}
		fmt.Printf("%d files processed (%d duplicates found)", counter, len(duplicates))
		filepath.WalkDir(sourceFolder, func(path string, d fs.DirEntry, err error) error {
			counter++
			if err != nil {
				return err
			}
			if !d.IsDir() {
				//println(path)
				h, _ := hashFileMD5(path)
				file := unique.Make(h)
				if slices.Contains(uniqueFiles, file) {
					//fmt.Printf("\n%s", path)
					duplicates = append(duplicates, path)
				} else {
					//println("Unique")
					uniqueFiles = append(uniqueFiles, file)
				}
			}
			fmt.Printf("\r")
			fmt.Printf("%d files processed (%d duplicates found)", counter, len(duplicates))
			return nil
		})
		if list {
			fmt.Println()
			for _, dup := range duplicates {
				fmt.Println(dup)
			}
		}
		if move {
			fmt.Println()
			fmt.Printf("Moving duplicate files to %s\n", destinationFolder)
			for _, dup := range duplicates {
				err := moveFile(dup, filepath.Join(destinationFolder, filepath.Base(dup)))
				if err != nil {
					fmt.Printf("Error: %s", err)
				}
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var sourceFolder string
var destinationFolder string
var list bool
var move bool

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&sourceFolder, "source", "s", "", "Source folder (to be dedup'ed)")
	rootCmd.PersistentFlags().StringVarP(&destinationFolder, "destination", "d", "", "Destination folder (for duplicated files)")
	rootCmd.MarkPersistentFlagRequired("source")

	rootCmd.PersistentFlags().BoolVar(&list, "list", false, "List duplicates")
	rootCmd.PersistentFlags().BoolVar(&move, "move", false, "Move duplicates")
	rootCmd.MarkFlagsMutuallyExclusive("list", "move")
	rootCmd.MarkFlagsRequiredTogether("move", "destination")
}
