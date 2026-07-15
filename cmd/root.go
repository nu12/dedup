/*
Copyright © 2026 nu12
*/
package cmd

import (
	"io"
	"os"

	"github.com/nu12/dedup/pkg/dedup"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dedup",
	Short: "Deduplicate files in a given directory",
	Long: `Deduplicate files in a given directory.
Examples: 

List duplicated files in the current directory
dedup -s . --list

Move duplicated files from the current directory to another one
dedup -s . --move -d ../destination-folder
`,
	Run: func(cmd *cobra.Command, args []string) {

		App := &dedup.Application{
			SourceFolder:      sourceFolder,
			DestinationFolder: destinationFolder,
			ListFlag:          list,
			MoveFlag:          move,

			OpenFunc:   func(path string) (io.ReadCloser, error) { return os.Open(path) },
			CopyFunc:   func(dst io.Writer, src io.Reader) (written int64, err error) { return io.Copy(dst, src) },
			CreateFunc: func(name string) (io.WriteCloser, error) { return os.Create(name) },
			RemoveFunc: func(name string) error { return os.Remove(name) },
		}
		App.Init().Run()
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
	rootCmd.PersistentFlags().StringVarP(&sourceFolder, "source", "s", "", "Source folder (to be dedup'ed)")
	rootCmd.PersistentFlags().StringVarP(&destinationFolder, "destination", "d", "", "Destination folder (for duplicated files)")
	rootCmd.MarkPersistentFlagRequired("source")

	rootCmd.PersistentFlags().BoolVar(&list, "list", false, "List duplicates")
	rootCmd.PersistentFlags().BoolVar(&move, "move", false, "Move duplicates")
	rootCmd.MarkFlagsMutuallyExclusive("list", "move")
	rootCmd.MarkFlagsRequiredTogether("move", "destination")

	rootCmd.AddCommand(versionCmd)
}
