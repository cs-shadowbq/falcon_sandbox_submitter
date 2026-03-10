/*
Copyright © 2024-2026 CrowdStrike - Scott MacGregor scott.macgregor@crowdstrike.com

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

// Package cmd provides the CLI commands for the falcon_sandbox submitter tool.
package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const licenseText = `Copyright © 2024-2026 CrowdStrike - Scott MacGregor

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
THE SOFTWARE.`

var (
	cfgFile           string
	verbose           bool
	debug             bool
	showVersion       bool
	buildClientID     string // Set by the build process
	clientID          string
	buildClientSecret string // Set by the build process
	clientSecret      string
	buildAPIBaseURL   string // Set by the build process
	clientCloud       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "falcon_sandbox",
	Short: "Submit files to the CrowdStrike Falcon Sandbox for analysis.",
	Long:  `Submit files to the CrowdStrike Falcon Sandbox for malware analysis. This command line tool allows you to submit files to the Falcon Sandbox for analysis against a variety of environments, and network settings.`,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		if showVersion {
			fmt.Printf("falcon_sandbox version %s\n\n%s\n", Version, licenseText)
			os.Exit(0)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, _ []string) {
		// reached only when no subcommand and --version not set
		cmd.Help() //nolint:errcheck
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

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.falcon_sandbox.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug output")
	rootCmd.PersistentFlags().BoolVar(&showVersion, "version", false, "print version and license then exit")

	rootCmd.PersistentFlags().StringVar(&clientID, "clientId", "", "Falcon CLIENT API ID")
	rootCmd.PersistentFlags().StringVar(&clientSecret, "clientSecret", "", "Falcon CLIENT SECRET API")
	rootCmd.PersistentFlags().StringVar(&clientCloud, "clientCloud", "", "Falcon CLIENT CLOUD API (us-1, us-2, eu-1, us-gov-1, gov1, *us-gov-2, *gov2)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// bind the configuration to file/environment variables
	cobra.CheckErr(viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")))
	viper.SetDefault("verbose", false)
	cobra.CheckErr(viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")))
	viper.SetDefault("debug", false)

	// if buildClientID is nil set it to ""
	/*
		if buildClientID == "" {
			buildClientID = "exClientId"
		}

		if buildClientSecret == "" {
			buildClientSecret = "exClientSecret"
		}

		if buildClientCloud == "" {
			buildClientCloud = "exClientCloud"
		}
	*/
	cobra.CheckErr(viper.BindPFlag("clientId", rootCmd.PersistentFlags().Lookup("clientId")))
	viper.SetDefault("clientId", buildClientID)
	cobra.CheckErr(viper.BindPFlag("clientSecret", rootCmd.PersistentFlags().Lookup("clientSecret")))
	viper.SetDefault("clientSecret", buildClientSecret)
	cobra.CheckErr(viper.BindPFlag("clientCloud", rootCmd.PersistentFlags().Lookup("clientCloud")))
	viper.SetDefault("clientCloud", buildAPIBaseURL)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".falcon_sandbox" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".falcon_sandbox")
	}

	// if debug is enabled, print the configuration
	// get the keys and print them in sorted order
	// Print the clientSecret first 4 characters then the rest as *
	debugOut("--- Compile and Switch Configuration ---")

	viper.AutomaticEnv() // read in environment variables that match

	debugOut("--- Environment Configuration ---")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

		debugOut("--- Config File Configuration ---")
	}
}

func debugOut(phase string) {
	if viper.GetBool("debug") {
		fmt.Println(phase)

		keys := viper.AllKeys()
		sort.Strings(keys)

		for _, key := range keys {
			if key == "clientsecret" {
				secret, _ := viper.Get(key).(string)
				if len(secret) > 4 {
					fmt.Printf("\t%s: %v\n", key, secret[:4]+"********")
				} else {
					fmt.Printf("\t%s: %v\n", key, "********")
				}
			} else {
				fmt.Printf("\t%s: %v\n", key, viper.Get(key))
			}
		}

		fmt.Println("-----------------------------------")
	}
}
