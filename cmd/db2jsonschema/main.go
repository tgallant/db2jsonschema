/*
Copyright © 2021 Tim Gallant <me@timgallant.us>

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
package main

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tgallant/db2jsonschema"
)

var (
	cfgFile    string
	driver     string
	dburl      string
	format     string
	outdir     string
	schematype string
	idtemplate string
	includes   []string
	excludes   []string
)

func HandleGenerate(cmd *cobra.Command, args []string) {
	if len(driver) == 0 || len(dburl) == 0 {
		cmd.Help()
		return
	}
	req := &db2jsonschema.Request{
		Driver:     driver,
		DataSource: dburl,
		Format:     format,
		Outdir:     outdir,
		SchemaType: schematype,
		IdTemplate: idtemplate,
		Includes:   includes,
		Excludes:   excludes,
	}
	err := req.Perform()
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "db2jsonschema",
	Short: "Generate JSON Schema definitions from database tables",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: HandleGenerate,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.db2jsonschema.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVar(&driver, "driver", "", "The DB Driver")
	rootCmd.Flags().StringVar(&dburl, "dburl", "", "The DB URL")
	rootCmd.Flags().StringVar(&format, "format", "", "The output format (json,yaml)")
	rootCmd.Flags().StringVar(&outdir, "outdir", "", "The output directory")
	rootCmd.Flags().StringVar(&schematype, "schematype", "", "The $schema value for the generated schemas")
	rootCmd.Flags().StringVar(&idtemplate, "idtemplate", "", "A template string for the $id value for the generated schemas")
	rootCmd.Flags().StringSliceVarP(&includes, "include", "", []string{}, "The tables to include")
	rootCmd.Flags().StringSliceVarP(&excludes, "exclude", "", []string{}, "The tables to exclude")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".db2jsonschema" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".db2jsonschema")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	Execute()
}
