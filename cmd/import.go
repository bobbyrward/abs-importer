/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/bobbyrward/abs-importer/pkg/api/audible"
	"github.com/bobbyrward/abs-importer/pkg/config"
	"github.com/bobbyrward/abs-importer/pkg/metadata"
	"github.com/bobbyrward/abs-importer/pkg/source"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import source title",
	Short: "Imports an audiobook while attaching metadata",
	RunE:  printError,
	Args:  cobra.ExactArgs(2),
}

func printError(cmd *cobra.Command, args []string) error {
	err := runImport(cmd, args)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	return err
}

func runImport(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	libraryConfig, err := getLibraryConfig(cmd)
	if err != nil {
		return err
	}

	sourceInfo, err := source.AnalyzeSource(getSourcePathFromArgs(args))
	if err != nil {
		return err
	}

	selectedMetadata, err := lookupMetadataByTitle(getTitleFromArgs(args))
	if err != nil {
		return err
	}

	directoryName, err := selectedMetadata.GenerateDirectoryName()
	if err != nil {
		return err
	}

	directoryName = sanitizeName(directoryName)
	fullDirName := path.Join(libraryConfig.Path, directoryName)

	if !sourceInfo.IsDir {
		err = moveFileToLibrary(fullDirName, sourceInfo)
		if err != nil {
			return err
		}
	} else {
		err = moveFolderToLibrary(fullDirName, sourceInfo)
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

func moveFileToLibrary(fullDirName string, sourceInfo source.SourceInfo) error {
	err := os.MkdirAll(fullDirName, 0777)
	if err != nil {
		return err
	}

	fmt.Printf("Linking %s to %s\n", sourceInfo.Filename, path.Join(fullDirName, path.Base(sourceInfo.Filename)))

	err = os.Link(sourceInfo.Filename, path.Join(fullDirName, path.Base(sourceInfo.Filename)))
	if err != nil {
		return err
	}

	return nil
}

func moveFolderToLibrary(fullDirName string, sourceInfo source.SourceInfo) error {
	fmt.Printf("Copying %s to %s\n", sourceInfo.Filename, fullDirName)

	cmd := exec.Command("cp", "-al", sourceInfo.Filename, fullDirName)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func sanitizeName(name string) string {
	replacer := strings.NewReplacer("/", "-")

	return replacer.Replace(name)
}

func getSourcePathFromArgs(args []string) string {
	return args[0]
}

func getTitleFromArgs(args []string) string {
	return args[1]
}

func getLibraryConfig(cmd *cobra.Command) (*config.LibraryConfig, error) {
	libraryName, err := cmd.Flags().GetString("library")
	if err != nil {
		return nil, err
	}

	for _, library := range Config.Libraries {
		if library.Name == libraryName {
			return &library, nil
		}
	}

	return nil, fmt.Errorf("library `%s` not found", libraryName)
}

func getBookMetadataFromASINs(aac *audible.AudibleApiClient, asins []string) ([]metadata.BookMetadata, error) {
	metadatas := make([]metadata.BookMetadata, 0, len(asins))

	for _, asin := range asins {
		md, err := aac.GetMetadataFromAsin(asin)
		if err != nil {
			fmt.Printf("WARNING: Skipping %s: %v\n", asin, err)
			metadatas = append(metadatas, metadata.BookMetadata{
				Title: "Skipped",
			})
			continue
		}

		metadatas = append(metadatas, md)
	}

	return metadatas, nil
}

func summarizeBookMetadatas(metadatas []metadata.BookMetadata) ([]string, error) {
	bookChoices := make([]string, 0, len(metadatas))

	for _, md := range metadatas {
		bookChoices = append(bookChoices, summarizeBookMetadata(&md))
	}

	return bookChoices, nil
}

func lookupMetadataByTitle(title string) (metadata.BookMetadata, error) {
	aac := audible.NewAudibleApiClient()
	md := metadata.BookMetadata{}

	asins, err := aac.SearchByTitle(title)
	if err != nil {
		return md, err
	}

	metadatas, err := getBookMetadataFromASINs(aac, asins)
	if err != nil {
		return md, err
	}

	bookChoices, err := summarizeBookMetadatas(metadatas)
	if err != nil {
		return md, err
	}

	index, err := selectBook(bookChoices)
	if err != nil {
		return md, err
	}

	return metadatas[index], nil
}

func selectBook(books []string) (int, error) {
	prompt := promptui.Select{
		Label: "Select a book",
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

	if md.Authors == nil {
		return truncatedTitle
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
	importCmd.Flags().String("library", "", "The library to import into")
}
