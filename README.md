# CrowdStrike Falcon Sandbox Submitter

[![Go Report Card](https://goreportcard.com/badge/github.com/cs-shadowbq/falcon_sandbox_submitter)](https://goreportcard.com/report/github.com/cs-shadowbq/falcon_sandbox_submitter)
[![GoDoc](https://godoc.org/github.com/cs-shadowbq/falcon_sandbox_submitter?status.svg)](https://godoc.org/github.com/cs-shadowbq/falcon_sandbox_submitter)

This tool is a golang implementation for submitting files to the CrowdStrike Falcon Sandbox. The CrowdStrike Falcon Sandbox is part of the CrowdStrike Counter Adversary Operations™ harnesses the power of the CrowdStrike Falcon® platform to provide a comprehensive, automated, and effective solution for identifying and stopping adversary operations. The CrowdStrike Falcon Sandbox provides deep analysis of evasive and unknown threats, enriches the results with intelligence, and delivers actionable indicators of compromise (IOCs) to the Falcon platform.

## Prerequisites

To use this tool, you need to have a CrowdStrike Falcon Platform account license with CrowdStrike Counter Adversary Operations Sandbox. With the account, you will have access to the CrowdStrike Falcon Sandbox API Write permissions. You will need to have the following information:

### Go Environment

Go must be installed and `GOPATH` must be set in your environment. If it is not set, `make` will fail during the dependency check. To configure it:

```bash
mkdir ~/gocode
export GOPATH=~/gocode
export PATH="$GOPATH/bin:$PATH"
```

Add those `export` lines to your shell profile (`~/.zshrc` or `~/.bash_profile`) to make them permanent. You can verify your current `GOPATH` with:

```bash
go env GOPATH
```

Then run `make` with it set:

```bash
GOPATH=$(go env GOPATH) make all
```

- API Client ID
- API Client Secret
- API Base URL
- API SCOPE:  
    "Sample uploads" Write *Yes*  
    "Sandbox (Falcon Intelligence)" Write *Yes*  

## Usage

```
Submit files to the CrowdStrike Falcon Sandbox for malware analysis. This command line tool allows you to submit files to the Falcon Sandbox for analysis against a variety of environments, and network settings.

Usage:
  falcon_sandbox [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  settings    Print the current configuration settings
  submit      SubCommand to submit a file to the CrowdStrike Falcon Sandbox for analysis.

Flags:
      --clientCloud string    Falcon CLIENT CLOUD API
      --clientId string       Falcon CLIENT API ID
      --clientSecret string   Falcon CLIENT SECRET API
      --config string         config file (default is $HOME/.falcon_sandbox.yaml)
      --debug                 debug output
  -h, --help                  help for falcon_sandbox
      --verbose               verbose output
```

## Installation

```bash
go get -u github.com/CrowdStrike/falcon-sandbox
```

## Usage Examples

Submit a file to the CrowdStrike Falcon Sandbox for analysis:

```bash
falcon_sandbox submit -f ~/CrowdStrikeTools/strings64.exe -s 110 --verbose
falcon_sandbox submit -f ~/CrowdStrikeTools/strings64.exe --environment 100 --verbose
falcon_sandbox submit -f ~/CrowdStrikeTools/strings64.exe --environment 110 --verbose
falcon_sandbox submit -a default_openie -e 160 -f ~/CrowdStrikeTools/strings64.exe  -n tor
```

To get help on the command line options, run the following command:

```bash
falcon_sandbox help
```

To view the current settings:

```bash
falcon_sandbox settings 
```

## Runtime Optional: Using the `.falcon_sandbox.yaml`

The runtime (not compile time) binary created, can leverage switches, ENV, as well as the `.falcon_sandbox.yaml`.

```yaml
verbose: false
clientSecret: SECRET-FROM-YAML-FILE
clientID: ID-FROM-YAML
```

## Submission Options

The `submit` subcommand has the following options:

```bash
Submit files to the CrowdStrike Falcon Sandbox for malware analysis. This command line tool allows you to submit files to the Falcon Sandbox for analysis against a variety of environments, and network settings.

Usage:
  falcon_sandbox submit [flags]

Flags:
  -a, --action_script string      Runtime script for sandbox analysis: ([default default_maxantievasion default_openie default_randomfiles default_randomtheme]) (default "default")
  -e, --environment int32         Specify the Environmental ID: ([100 110 140 160 200 300 400]) 
                                    400: MacOS Catalina 10.15, 64-bit
                                    300: Linux Ubuntu 16.04, 64-bit
                                    200: Android (static analysis)
                                    160: Windows 10, 64-bit
                                    140: Windows 11, 64-bit
                                    110: Windows 7, 64-bit
                                    100: Windows 7, 32-bit
                                    NOTE: GOV1 clouds support only: 140, 160 (default 160)
  -f, --filename string           The file to submit to the sandbox[ie. sample.exe]
  -h, --help                      help for submit
  -n, --network_settings string   Specifies the sandbox network_settings used for analysis: ([default offline simulated tor]) (default "default")
                                    NOTE: GOV1 clouds support only: default

Global Flags:
      --clientCloud string    Falcon CLIENT CLOUD API (us-1, us-2, eu-1, us-gov-1, gov1, *us-gov-2, *gov2)
      --clientId string       Falcon CLIENT API ID
      --clientSecret string   Falcon CLIENT SECRET API
      --config string         config file (default is $HOME/.falcon_sandbox.yaml)
      --debug                 debug output
      --verbose               verbose output
```

## Compile Optional: Compiling with API Keys for write access upload permissions via `.env`

The `.env` can be used to define the FALCON API Client, Secret, and Cloud settings for the `makefile` to compile into the binary. This functionality may be optimal for the environment where you are deploying the `falcon_sandbox` submission tool. 

The following commands can be used to compile the falcon-sandbox tool with the CLIENT_ID and CLIENT_SECRET:

```shell
$> git clone github.com/cs-shadowbq/falcon-sandbox
$> cd falcon-sandbox
```

Edit the `.env` file and add your `FALCON_CLIENT_ID` and `FALCON_CLIENT_SECRET` and `FALCON_API_BASE_URL`

```ini
FALCON_CLIENT_ID=aaaaaaa
FALCON_CLIENT_SECRET=bbbbbbb
FALCON_API_BASE_URL=us-1
```

Compile the binaries 

```shell
$> make all
```

Raw Compilation example without using the Makefile:

```shell
$> go build -ldflags "-X github.com/cs-shadowbq/falcon_sandbox_submitter/cmd.buildClientId=YOUR_CLIENT_ID -X github.com/cs-shadowbq/falcon_sandbox_submitter/cmd.buildClientSecret=YOUR_CLIENT_SECRET" -o falcon-sandbox main.go
```

## Cross-Compilation

The following commands can be used to cross-compile the falcon-sandbox tool for different operating systems and architectures.

> **Note:** `GOPATH` must be set in your environment before running `make`. If it is not already exported in your shell profile, prefix the command: `GOPATH=$(go env GOPATH) make all`. See [Prerequisites](#prerequisites) for setup instructions.

```shell
$> GOPATH=$(go env GOPATH) make all
Product Version 1.0.0

Checking Build Dependencies ---->

Cleaning Build ---->
rm -f -rf pkg/*
rm -f -rf build/*
rm -f -rf tmp/*

Building ---->
env GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/cs-shadowbq/falcon_sandbox_submitter/cmd.Version=x.y.z" -o build/falcon_sandbox_linux_amd64 main.go
env GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/cs-shadowbq/falcon_sandbox_submitter/cmd.Version=x.y.z" -o build/falcon_sandbox.exe main.go
env GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/cs-shadowbq/falcon_sandbox_submitter/cmd.Version=x.y.z" -o build/falcon_sandbox_darwin_amd64 main.go
```

You can get a list of supported cross-compilation targets by running the following command:

```bash
go tool dist list
```

## Code Signing

To sign the binary, you will need to have a valid code signing certificate. Edit the makefile section `codesign` and list your certificate name. Then run the following command:

```shell
make codesign
```
