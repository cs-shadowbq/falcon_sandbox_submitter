/*
Copyright © 2024-2026 CrowdStrike - Scott MacGregor scott.macgregor@crowdstrike.com
*/

// Package cmd provides the CLI commands for the falcon_sandbox submitter tool.
package cmd

import (
	"fmt"
	"sort"

	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/cs-shadowbq/falcon_sandbox_submitter/sandbox"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	filename        string
	sandboxEnvID    int32
	actionScript    string
	networkSettings string
)

// ValidEnvValues holds the set of valid sandbox environment IDs.
type ValidEnvValues struct {
	Values map[int32]bool
}

// NewValidEnvValues returns a ValidEnvValues populated with all supported environment IDs.
func NewValidEnvValues() *ValidEnvValues {
	return &ValidEnvValues{
		Values: map[int32]bool{
			100: true,
			110: true,
			140: true,
			160: true,
			200: true,
			300: true,
			400: true,
		},
	}
}

// NewGov1ValidEnvValues returns the restricted set of environment IDs supported
// by the GOV1 cloud (us-gov-1 / gov1). Only Windows 10 and Windows 11 64-bit
// detonation environments are available in the government cloud.
func NewGov1ValidEnvValues() *ValidEnvValues {
	return &ValidEnvValues{
		Values: map[int32]bool{
			140: true, // Windows 11, 64-bit
			160: true, // Windows 10, 64-bit
		},
	}
}

// NewGov1ValidNetworkSettings returns the restricted set of network_settings
// supported by the GOV1 cloud. Only "default" is available.
func NewGov1ValidNetworkSettings() *ValidNetworkSettings {
	return &ValidNetworkSettings{
		Values: map[string]bool{
			"default": true,
		},
	}
}

// ValidActionValues holds the set of valid action script names.
type ValidActionValues struct {
	Values map[string]bool
}

// NewValidActionSettings returns a ValidActionValues populated with all supported action script names.
func NewValidActionSettings() *ValidActionValues {
	return &ValidActionValues{
		Values: map[string]bool{
			"default":                true,
			"default_maxantievasion": true,
			"default_randomfiles":    true,
			"default_randomtheme":    true,
			"default_openie":         true,
		},
	}
}

// ValidNetworkSettings holds the set of valid network settings values.
type ValidNetworkSettings struct {
	Values map[string]bool
}

// NewValidNetworkSettings returns a ValidNetworkSettings populated with all supported network settings.
func NewValidNetworkSettings() *ValidNetworkSettings {
	return &ValidNetworkSettings{
		Values: map[string]bool{
			"default":   true,
			"tor":       true,
			"simulated": true,
			"offline":   true,
		},
	}
}

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "SubCommand to submit a file to the CrowdStrike Falcon Sandbox for analysis.",
	Long:  `Submit files to the CrowdStrike Falcon Sandbox for malware analysis. This command line tool allows you to submit files to the Falcon Sandbox for analysis against a variety of environments, and network settings.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return fmt.Errorf("failed to get verbose flag: %w", err)
		}
		sub := sandbox.CmdSubmission{
			FalconClientID:     viper.GetString("clientId"),
			FalconClientSecret: viper.GetString("clientSecret"),
			ClientCloud:        viper.GetString("clientCloud"),
			Filename:           filename,
			SandboxEnvID:       sandboxEnvID,
			NetworkSettings:    networkSettings,
			ActionScript:       actionScript,
		}
		return sub.SubmitFile(verbose)
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
	submitCmd.Flags().Int32VarP(&sandboxEnvID, "environment", "e", 160, fmt.Sprintf("Specify the Environmental ID: (%v) \n  400: MacOS Catalina 10.15, 64-bit\n  300: Linux Ubuntu 16.04, 64-bit\n  200: Android (static analysis)\n  160: Windows 10, 64-bit\n  140: Windows 11, 64-bit\n  110: Windows 7, 64-bit\n  100: Windows 7, 32-bit\n  NOTE: GOV1 clouds support only: 140, 160", getValidEnvValues(validEnvValues)))
	if err := submitCmd.MarkFlagRequired("filename"); err != nil {
		panic(err)
	}

	submitCmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return fmt.Errorf("invalid value for a flag: %v", err)
	})

	submitCmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		// Validate the cloud string before doing anything else.
		cloud, err := falcon.CloudValidate(viper.GetString("clientCloud"))
		if err != nil {
			return fmt.Errorf("invalid clientCloud value: %w", err)
		}

		// GOV2 has no CAO Sandbox capability at all.
		if cloud == falcon.CloudUsGov2 || cloud == falcon.CloudGov2 {
			return fmt.Errorf("cloud %q does not support CrowdStrike CAO Falcon Sandbox submissions", cloud.String())
		}

		// Choose the allowed environment and network_settings lists based on cloud region.
		allowedEnvs := validEnvValues
		allowedNetworkSettings := validNetworkSettings
		if cloud == falcon.CloudUsGov1 || cloud == falcon.CloudGov1 {
			allowedEnvs = NewGov1ValidEnvValues()
			allowedNetworkSettings = NewGov1ValidNetworkSettings()
		}

		environment, err := cmd.Flags().GetInt32("environment")
		if err != nil {
			return fmt.Errorf("failed to get environment flag: %w", err)
		}
		if !allowedEnvs.Values[environment] {
			return fmt.Errorf("invalid value for environment. Allowed values for cloud %q: %v", cloud.String(), getValidEnvValues(allowedEnvs))
		}

		actionScript, err := cmd.Flags().GetString("action_script")
		if err != nil {
			return fmt.Errorf("failed to get action_script flag: %w", err)
		}
		if !validActionValues.Values[actionScript] {
			return fmt.Errorf("invalid value for action_script. It must be %v", getValidActionValues(validActionValues))
		}

		networkSettings, err := cmd.Flags().GetString("network_settings")
		if err != nil {
			return fmt.Errorf("failed to get network_settings flag: %w", err)
		}
		if !allowedNetworkSettings.Values[networkSettings] {
			return fmt.Errorf("invalid value for network_settings. Allowed values for cloud %q: %v", cloud.String(), getValidNetworkSettings(allowedNetworkSettings))
		}

		return nil
	}
}
