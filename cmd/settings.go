/*
Copyright © 2024 CrowdStrike - Scott MacGregor scott.macgregor@crowdstrike.com
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version string // set by the build process
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Print the current configuration settings",
	Long: `Print the current configuration settings. This command is useful to see the final configuration once all the settings have been applied. 
	It also shows how to access the global flags and command flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("Settings called")

		if verbose {
			fmt.Println("--- Final configuration  ---")
		}
		//fmt.Printf("\tVersion: %s\n", Version)
		//for s, i := range viper.AllSettings() {
		//	fmt.Printf("\t%s: %v\n", s, i)
		//}

		keys := viper.AllSettings()
		//sort keys
		var keysSorted []string
		for key := range keys {
			keysSorted = append(keysSorted, key)
		}

		// get the keys and print them in sorted order
		for _, key := range keysSorted {
			if key == "clientsecret" {
				// Print the clientSecret first 4 characters then the rest as *
				fmt.Printf("\t%s: %v\n", key, viper.Get(key).(string)[:4]+"********")
			} else {
				fmt.Printf("\t%s: %v\n", key, viper.Get(key))
			}
		}

		if verbose {
			fmt.Println("----------------------------")
		}

	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)

}
