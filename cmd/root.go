/*
Copyright © 2024 CrowdStrike - Scott MacGregor scott.macgregor@crowdstrike.com

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
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile           string
	verbose           bool
	debug             bool
	buildClientId     string // Set by the build process
	clientId          string
	buildClientSecret string // Set by the build process
	clientSecret      string
	buildApiBaseUrl   string // Set by the build process
	clientCloud       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "falcon_sandbox",
	Short: "Submit files to the CrowdStrike Falcon Sandbox for analysis.",
	Long:  `Submit files to the CrowdStrike Falcon Sandbox for malware analysis. This command line tool allows you to submit files to the Falcon Sandbox for analysis against a variety of environments, and network settings.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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

	rootCmd.PersistentFlags().StringVar(&clientId, "clientId", "", "Falcon CLIENT API ID")
	rootCmd.PersistentFlags().StringVar(&clientSecret, "clientSecret", "", "Falcon CLIENT SECRET API")
	rootCmd.PersistentFlags().StringVar(&clientCloud, "clientCloud", "", "Falcon CLIENT CLOUD API")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// bind the configuration to file/environment variables
	cobra.CheckErr(viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")))
	viper.SetDefault("verbose", false)
	cobra.CheckErr(viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")))
	viper.SetDefault("debug", false)

	// if buildClientId is nil set it to ""
	/*
		if buildClientId == "" {
			buildClientId = "exClientId"
		}

		if buildClientSecret == "" {
			buildClientSecret = "exClientSecret"
		}

		if buildClientCloud == "" {
			buildClientCloud = "exClientCloud"
		}
	*/
	cobra.CheckErr(viper.BindPFlag("clientId", rootCmd.PersistentFlags().Lookup("clientId")))
	viper.SetDefault("clientId", buildClientId)
	cobra.CheckErr(viper.BindPFlag("clientSecret", rootCmd.PersistentFlags().Lookup("clientSecret")))
	viper.SetDefault("clientSecret", buildClientSecret)
	cobra.CheckErr(viper.BindPFlag("clientCloud", rootCmd.PersistentFlags().Lookup("clientCloud")))
	viper.SetDefault("clientCloud", buildApiBaseUrl)
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
	//var keys []string
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

		var keys []string = viper.AllKeys()
		sort.Strings(keys)

		for _, key := range keys {
			if key == "clientsecret" {

				fmt.Printf("\t%s: %v\n", key, viper.Get(key).(string)[:4]+"********")
			} else {
				fmt.Printf("\t%s: %v\n", key, viper.Get(key))
			}
		}

		fmt.Println("-----------------------------------")
	}
}
