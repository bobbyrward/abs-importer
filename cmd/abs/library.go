package abs

// import (
// 	"fmt"
//
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
//
// 	"github.com/bobbyrward/abs-importer/cmd/output"
// 	"github.com/bobbyrward/abs-importer/pkg/api/audiobookshelf"
// )
//
// // importCmd represents the import command
// var libraryCmd = &cobra.Command{
// 	Use:   "library",
// 	Short: "Commands for interacting with Audiobookshelf Libraries",
// }
//
// var libraryLsCmd = &cobra.Command{
// 	Use:   "ls",
// 	Short: "List libraries in Audiobookshelf",
// 	RunE:  listLibraries,
// }
//
// var libraryScanCmd = &cobra.Command{
// 	Use:   "scan",
// 	Short: "Scan a library's folders",
// 	RunE:  scanLibrary,
// 	Args:  cobra.RangeArgs(0, 1),
// }
//
// var libraryGetCmd = &cobra.Command{
// 	Use:   "get",
// 	Short: "Get a library",
// 	RunE:  getLibrary,
// 	Args:  cobra.RangeArgs(0, 1),
// }
//
// func getLibraryId(args []string) (string, error) {
// 	var libraryId string
//
// 	if len(args) == 1 {
// 		libraryId = args[0]
// 	} else {
// 		libraryId = viper.GetString("libraryId")
// 		if libraryId == "" {
// 			return "", fmt.Errorf("Library Id not configured")
// 		}
// 	}
//
// 	return libraryId, nil
// }
//
// func scanLibrary(cmd *cobra.Command, args []string) error {
// 	aac := audiobookshelf.NewClientFromViper()
//
// 	force, err := cmd.Flags().GetBool("force")
// 	if err != nil {
// 		return err
// 	}
//
// 	libraryId, err := getLibraryId(args)
// 	if err != nil {
// 		return err
// 	}
//
// 	return aac.ScanLibrary(libraryId, force)
// }
//
// func getLibrary(cmd *cobra.Command, args []string) error {
// 	aac := audiobookshelf.NewClientFromViper()
//
// 	libraryId, err := getLibraryId(args)
// 	if err != nil {
// 		return err
// 	}
//
// 	library, err := aac.GetLibrary(libraryId, false)
// 	if err != nil {
// 		return err
// 	}
//
// 	gout, err := output.NewGout(cmd)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = gout.Print(library)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func listLibraries(cmd *cobra.Command, args []string) error {
// 	aac := audiobookshelf.NewClientFromViper()
//
// 	libraries, err := aac.ListLibraries()
// 	if err != nil {
// 		return err
// 	}
//
// 	gout, err := output.NewGout(cmd)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = gout.Print(libraries)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func init() {
// 	libraryCmd.AddCommand(libraryLsCmd)
// 	output.AddOutputFields(libraryLsCmd)
//
// 	libraryCmd.AddCommand(libraryGetCmd)
// 	output.AddOutputFields(libraryGetCmd)
//
// 	libraryScanCmd.Flags().Bool("force", false, "Whether to force a rescan for all of a library's items")
//
// 	libraryCmd.AddCommand(libraryScanCmd)
// }
