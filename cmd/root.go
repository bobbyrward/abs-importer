/*
Copyright © 2023 Bobby R. Ward

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/abs-importer/cmd/abs"
	"github.com/bobbyrward/abs-importer/pkg/config"
)

var (
	cfgFile string
	Config  *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "abs-importer",
	Short: "A tool to help with importing into ABS",
	Long:  "A tool to help with importing into ABS",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.abs-importer.yaml)")

	// rootCmd.PersistentFlags().String("api-token", "", "The ABS api token")
	// rootCmd.PersistentFlags().String("library", "", "The root of your library")
	// rootCmd.PersistentFlags().String("libraryId", "", "The id of your audiobook library")

	// viper.BindPFlag("apiToken", rootCmd.PersistentFlags().Lookup("api-token"))
	// viper.BindPFlag("libraryRoot", rootCmd.PersistentFlags().Lookup("libraryRoot"))
	// viper.BindPFlag("libraryId", rootCmd.PersistentFlags().Lookup("libraryId"))

	rootCmd.AddCommand(abs.AbsCmd)
}

func initConfig() {
	cfg, err := config.NewConfig(config.WithConfigFilename(cfgFile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: unable to load config: %v", err)
		os.Exit(1)
	}

	Config = cfg
}
