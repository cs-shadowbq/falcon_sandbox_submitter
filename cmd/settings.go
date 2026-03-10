/*
Copyright © 2024-2026 CrowdStrike - Scott MacGregor scott.macgregor@crowdstrike.com
*/
package cmd

import (
	"fmt"
	"sort"

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
		var keysSorted []string
		for key := range keys {
			keysSorted = append(keysSorted, key)
		}
		sort.Strings(keysSorted)

		// get the keys and print them in sorted order
		for _, key := range keysSorted {
			if key == "clientsecret" {
				// Mask all but the first 4 characters of the secret.
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

		if verbose {
			fmt.Println("----------------------------")
		}

	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)

}
