/*
Copyright © 2024 CrowdStrike - Scott MacGregor scott.macgregor@crowdstrike.com
*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/cs-shadowbq/falcon_sandbox/sandbox"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	filename        string
	sandboxEnvId    int32
	actionScript    string
	networkSettings string
)

type ValidEnvValues struct {
	Values map[int32]bool
}

func NewValidEnvValues() *ValidEnvValues {
	return &ValidEnvValues{
		Values: map[int32]bool{
			100: true,
			110: true,
			140: true,
			160: true,
			200: true,
			300: true,
			400: true},
	}
}

type ValidActionValues struct {
	Values map[string]bool
}

func NewValidActionSettings() *ValidActionValues {
	return &ValidActionValues{
		Values: map[string]bool{
			"default":                true,
			"default_maxantievasion": true,
			"default_randomfiles":    true,
			"default_randomtheme":    true,
			"default_openie":         true},
	}
}

type ValidNetworkSettings struct {
	Values map[string]bool
}

func NewValidNetworkSettings() *ValidNetworkSettings {
	return &ValidNetworkSettings{
		Values: map[string]bool{
			"default":   true,
			"tor":       true,
			"simulated": true,
			"offline":   true},
	}
}

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "SubCommand to submit a file to the CrowdStrike Falcon Sandbox for analysis.",
	Long:  `Submit files to the CrowdStrike Falcon Sandbox for malware analysis. This command line tool allows you to submit files to the Falcon Sandbox for analysis against a variety of environments, and network settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if global flag was set
		verbose, _ := cmd.Flags().GetBool("verbose")
		sub := sandbox.CmdSubmission{
			FalconClientId:     viper.GetString("clientId"),
			FalconClientSecret: viper.GetString("clientSecret"),
			ClientCloud:        viper.GetString("clientCloud"),
			Filename:           filename,
			SandboxEnvId:       sandboxEnvId,
			NetworkSettings:    networkSettings,
			ActionScript:       actionScript,
		}
		sub.SubmitFile(verbose)
	},
}

type sortedInt32Array []int32

func (f sortedInt32Array) Len() int {
	return len(f)
}

func (f sortedInt32Array) Less(i, j int) bool {
	return f[i] < f[j]
}

func (f sortedInt32Array) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func getValidEnvValues(validValues *ValidEnvValues) []int32 {
	var values sortedInt32Array
	for value := range validValues.Values {
		values = append(values, value)
	}
	// sort the values
	sort.Sort(values)
	return values
}

func getValidActionValues(validValues *ValidActionValues) []string {
	var values []string
	for value := range validValues.Values {
		values = append(values, value)
	}
	// sort the values
	sort.Strings(values)
	return values
}

func getValidNetworkSettings(validValues *ValidNetworkSettings) []string {
	var values []string
	for value := range validValues.Values {
		values = append(values, value)
	}
	// sort the values
	sort.Strings(values)
	return values
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	/*
		Environment ID for the sandbox. Values:

			- `400`: MacOS Catalina 10.15, 64-bit
			- `300`: Linux Ubuntu 16.04, 64-bit
			- `200`: Android (static analysis)
			- `160`: Windows 10, 64-bit
			- `140`: Windows 11, 64-bit
			- `110`: Windows 7, 64-bit
			- `100`: Windows 7, 32-bit
	*/

	/*
		**`action_script`** (optional): Runtime script for sandbox analysis. Values:

			- `default`
			- `default_maxantievasion`
			- `default_randomfiles`
			- `default_randomtheme`
			- `default_openie`
	*/

	/*
		**`network_settings`** (optional): Specifies the sandbox network_settings used for analysis. Values:

			- `default`: Fully operating network
			- `tor`: Route network traffic via TOR
			- `simulated`: Simulate network traffic
			- `offline`: No network traffic

	*/

	validEnvValues := NewValidEnvValues()
	validActionValues := NewValidActionSettings()
	validNetworkSettings := NewValidNetworkSettings()

	submitCmd.Flags().StringVarP(&filename, "filename", "f", "", "The file to submit to the sandbox[ie. sample.exe]")
	submitCmd.Flags().StringVarP(&actionScript, "action_script", "a", "default", fmt.Sprintf("Runtime script for sandbox analysis: (%v)", getValidActionValues(validActionValues)))
	submitCmd.Flags().StringVarP(&networkSettings, "network_settings", "n", "default", fmt.Sprintf("Specifies the sandbox network_settings used for analysis: (%v)", getValidNetworkSettings(validNetworkSettings)))
	submitCmd.Flags().Int32VarP(&sandboxEnvId, "environment", "e", 160, fmt.Sprintf("Specify the Environmental ID: (%v) \n  400: MacOS Catalina 10.15, 64-bit\n  300: Linux Ubuntu 16.04, 64-bit\n  200: Android (static analysis)\n  160: Windows 10, 64-bit\n  140: Windows 11, 64-bit\n  110: Windows 7, 64-bit\n  100: Windows 7, 32-bit", getValidEnvValues(validEnvValues)))
	submitCmd.MarkFlagRequired("filename")
	//submitCmd.MarkFlagRequired("environment")

	submitCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		//return fmt.Errorf("invalid value for environment. It must be %v", getValidEnvValues(validEnvValues))
		return fmt.Errorf("invalid value for a flag: %v", err)
	})

	submitCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		environment, _ := cmd.Flags().GetInt32("environment")
		if !validEnvValues.Values[environment] {
			return fmt.Errorf("invalid value for environment. It must be %v", getValidEnvValues(validEnvValues))
		}

		actionScript, _ := cmd.Flags().GetString("action_script")
		if !validActionValues.Values[actionScript] {
			return fmt.Errorf("invalid value for action_script. It must be %v", getValidActionValues(validActionValues))
		}

		networkSettings, _ := cmd.Flags().GetString("network_settings")
		if !validNetworkSettings.Values[networkSettings] {
			return fmt.Errorf("invalid value for network_settings. It must be %v", getValidNetworkSettings(validNetworkSettings))
		}

		return nil
	}

	//validate that sandboxEnvId is contains one of the following values 100, 110, 160, 200, 300

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// submitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
