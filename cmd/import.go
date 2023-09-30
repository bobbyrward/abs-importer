/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/bobbyrward/abs-importer/pkg"
	"github.com/bobbyrward/abs-importer/pkg/api/audible"
	"github.com/bobbyrward/abs-importer/pkg/metadata"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports an audiobook while attaching metadata",
	RunE:  runImport,
	Args:  cobra.ExactArgs(2),
}

func runImport(cmd *cobra.Command, args []string) error {
	sourceInfo, err := pkg.AnalyzeSource(args[0])
	if err != nil {
		return err
	}

	aac := audible.NewAudibleApiClient()
	asins, err := aac.SearchByTitle(args[1])
	if err != nil {
		return err
	}

	metadatas := make([]metadata.BookMetadata, 0, len(asins))

	for _, asin := range asins {
		md, err := aac.GetMetadataFromAsin(asin)
		if err != nil {
			return err
		}

		metadatas = append(metadatas, md)
	}

	bookChoices := make([]string, 0, len(asins))

	for _, md := range metadatas {
		summary := summarizeBookMetadata(&md)
		bookChoices = append(bookChoices, summary)
	}

	index, err := selectBook(bookChoices)
	if err != nil {
		return err
	}

	selectedMetadata := metadatas[index]

	directoryName, err := selectedMetadata.GenerateDirectoryName()
	if err != nil {
		return err
	}

	fullDirName := path.Join(viper.GetString("libraryRoot"), directoryName)

	if !sourceInfo.IsDir {
		err = os.MkdirAll(fullDirName, 0777)
		if err != nil {
			return err
		}

		err = os.Link(sourceInfo.Filename, path.Join(fullDirName, path.Base(sourceInfo.Filename)))
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command("cp", "-al", sourceInfo.Filename, fullDirName)
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	err = selectedMetadata.WriteOpf(path.Join(fullDirName, "metadata.opf"))
	if err != nil {
		return err
	}

	return nil
}

func selectBook(books []string) (int, error) {
	prompt := promptui.Select{
		Label: "Select a  book",
		Items: books,
	}

	index, _, err := prompt.Run()

	if err != nil {
		return 0, err
	}

	return index, nil
}

func summarizeBookMetadata(md *metadata.BookMetadata) string {
	var buffer strings.Builder

	truncatedTitle := md.Title

	if len(truncatedTitle) > 80 {
		truncatedTitle = truncatedTitle[:77] + "..."
	}

	buffer.WriteString(truncatedTitle)
	buffer.WriteString(" by ")
	buffer.WriteString(md.Authors[0].Name)

	for i := 1; i < len(md.Authors); i++ {
		buffer.WriteString(" & ")
		buffer.WriteString(md.Authors[i].Name)
	}

	if md.PrimarySeries != nil {
		buffer.WriteString(" - ")
		buffer.WriteString(md.PrimarySeries.Name)

		if md.PrimarySeries.Position != nil {
			buffer.WriteString(" ")
			buffer.WriteString(*md.PrimarySeries.Position)
		}
	}

	return buffer.String()
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
