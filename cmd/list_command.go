/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cmd

import (
	"ballerina-update-tool/constants"
	"ballerina-update-tool/exceptions"
	"ballerina-update-tool/utils"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// NewListCmd creates a new list command using Cobra and CommandBase
func NewListCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"list",
		"List Ballerina Distributions",
		`List locally and remotely available Ballerina distributions.

This command displays all Ballerina distributions that are installed locally
and those that are available for download from the remote server.`,
		printStream,
	)

	// Define flags
	var allFlag bool
	var preReleasesFlag bool

	// Add flags
	cmd.Flags().BoolVarP(&allFlag, "all", "a", false, "List all distributions (not just the most recent ones)")
	cmd.Flags().BoolVarP(&preReleasesFlag, "pre-releases", "p", false, "Include pre-release versions in the listings")

	// Add example
	cmd.Example = `  # List local and recent remote distributions
  bal dist list

  # List all available distributions
  bal dist list --all

  # List including pre-release versions
  bal dist list --pre-releases

  # List all including pre-release versions
  bal dist list --all --pre-releases`

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// Check for too many arguments
		if len(args) > 0 {
			panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", constants.BallerinaCliCommands.LIST))
		}

		// Execute list command
		listDistributions(cmdBase.GetPrintStream(), allFlag, preReleasesFlag)
	}

	return cmd
}

// For backward compatibility with the existing command system

// ListCommandStruct is a wrapper for backward compatibility
type ListCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	ListCommands    []string
	HelpFlag        bool
	AllFlag         bool
	PreReleasesFlag bool
	ParentCmdParser interface{}
}

// NewList creates a new ListCommand for backward compatibility
func NewList(printStream io.Writer) *ListCommandStruct {
	// Create the Cobra command
	cobraCmd := NewListCmd(printStream)

	// Create the wrapper
	cmd := &ListCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *ListCommandStruct) Execute() {
	if cmd.HelpFlag {
		cmd.cmdBase.PrintUsageInfo(constants.CommandToolConstants.CliHelpFilePrefix + cmd.GetName())
		return
	}

	// Apply flags to the cobra command
	cmd.cobraCmd.Flags().Set("all", fmt.Sprintf("%v", cmd.AllFlag))
	cmd.cobraCmd.Flags().Set("pre-releases", fmt.Sprintf("%v", cmd.PreReleasesFlag))

	if cmd.ListCommands == nil {
		listDistributions(cmd.cmdBase.GetPrintStream(), cmd.AllFlag, cmd.PreReleasesFlag)
		return
	}

	if len(cmd.ListCommands) > 0 {
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
	}
}

// GetName returns the name of the command
func (cmd *ListCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.LIST
}

// PrintLongDesc prints the long description of the command
func (cmd *ListCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *ListCommandStruct) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal dist list\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *ListCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *ListCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *ListCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}

// listDistributions lists distributions in the local and remote
func listDistributions(outStream io.Writer, allFlag bool, prFlag bool) {
	currentBallerinaVersion := utils.ToolUtil.GetCurrentBallerinaVersion()
	folder := filepath.FromSlash(utils.ToolUtil.GetDistributionsPath())
	listOfFiles, _ := ioutil.ReadDir(folder)
	maxListingDistributions := 10

	try := func(fn func()) (cmdErr *exceptions.CommandException) {
		defer func() {
			if r := recover(); r != nil {
				if ce, ok := r.(exceptions.CommandException); ok {
					// Create a new instance to get a pointer
					cmdErr = exceptions.New()
					// Copy messages
					for _, msg := range ce.GetMessages() {
						cmdErr.AddMessage(msg)
					}
				} else if cePtr, ok := r.(*exceptions.CommandException); ok {
					cmdErr = cePtr
				} else {
					panic(r)
				}
			}
		}()
		fn()
		return nil
	}

	cmdErr := try(func() {
		distList := make(map[string]interface{})
		channelsArr := make([]interface{}, 0)
		channels := utils.ToolUtil.GetDistributions(outStream)

		if len(listOfFiles) > 0 {
			fmt.Fprintln(outStream, "Distributions available locally: \n")
		}

		for _, channel := range channels {
			channelJson := make(map[string]interface{})
			releases := make([]interface{}, 0)
			channelJson["name"] = channel.Name
			channelDistListLocal := channel.Distributions

			if !strings.Contains(channel.Name, constants.CommandToolConstants.PreRelease) {
				channelDistListLocal = getSortedDistList(channelDistListLocal)
			}

			if len(listOfFiles) > 0 {
				for _, distribution := range channelDistListLocal {
					for _, file := range listOfFiles {
						if file.IsDir() {
							version := ""
							parts := strings.Split(file.Name(), "-")
							if len(parts) == 2 {
								version = parts[1]
							}
							if version == distribution.Version {
								versionName := distribution.Name
								versionId := distribution.Version
								versionInfo := make(map[string]interface{})
								versionInfo["name"] = versionName
								versionInfo["version"] = versionId
								fmt.Fprintln(outStream, markVersion(currentBallerinaVersion, versionId))
								releases = append(releases, versionInfo)
							}
						}
					}
				}
			}
			channelJson["releases"] = releases
			channelsArr = append(channelsArr, channelJson)
		}
		distList["channels"] = channelsArr
		writeLocalDistsIntoJson(distList)

		fmt.Fprintln(outStream, "\nDistributions available remotely:")
		for _, channel := range channels {
			if strings.Contains(channel.Name, constants.CommandToolConstants.PreRelease) && !prFlag {
				continue
			} else {
				fmt.Fprintf(outStream, "\n%s\n\n", channel.Name)
				channelDistList := channel.Distributions
				if !strings.Contains(channel.Name, constants.CommandToolConstants.PreRelease) {
					channelDistList = getSortedDistList(channelDistList)
				}
				if !allFlag {
					if len(channelDistList) > maxListingDistributions {
						recentDistributions := channelDistList[:maxListingDistributions]
						for _, distribution := range recentDistributions {
							fmt.Fprintln(outStream, markVersion(currentBallerinaVersion, distribution.Version, channelDistList[0].Version))
						}
					} else {
						for _, distribution := range channelDistList {
							fmt.Fprintln(outStream, markVersion(currentBallerinaVersion, distribution.Version, channelDistList[0].Version))
						}
					}
				} else {
					for _, distribution := range channelDistList {
						fmt.Fprintln(outStream, markVersion(currentBallerinaVersion, distribution.Version, channelDistList[0].Version))
					}
				}
			}
		}
	})

	if cmdErr != nil {
		fmt.Fprintln(outStream, "Distributions available locally: \n")
		if len(listOfFiles) > 0 && !isUpdated(listOfFiles) {
			listLocalDists(listOfFiles, outStream, currentBallerinaVersion)
		} else {
			readLocalDistsFromJson(outStream, currentBallerinaVersion)
		}
		fmt.Fprintln(outStream, "\nDistributions available remotely: \n")
		utils.ErrorUtil.PrintLauncherException(*cmdErr, outStream)
	}

	fmt.Fprintln(outStream)
	if !allFlag {
		fmt.Fprintln(outStream, "Use 'bal dist list -a' to list all the distributions under each channel. ")
	}
	fmt.Fprintln(outStream, "Use 'bal help dist' for more information on specific commands.")
}

// markVersion checks used Ballerina version and mark the output
func markVersion(used string, current string, latest ...string) string {
	if len(latest) == 1 {
		usedMarker := "  "
		latestMarker := ""
		if used == current {
			usedMarker = "* "
		}
		if current == latest[0] {
			latestMarker = " - latest"
		}
		return usedMarker + current + latestMarker
	} else {
		if used == current {
			return "* " + current
		} else {
			return "  " + current
		}
	}
}

// writeLocalDistsIntoJson writes the locally available distributions into a json file
func writeLocalDistsIntoJson(distList map[string]interface{}) {
	try := func() {
		distListFilePath := utils.OSUtils.GetBallerinaDistListFilePath()
		file, err := os.Stat(distListFilePath)

		if os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(distListFilePath), 0755)
			_, err = os.Create(distListFilePath)
			if err != nil {
				panic(utils.ErrorUtil.CreateCommandException("failed to create file: " + err.Error()))
			}
			utils.ToolUtil.AddWritePermissionToFile(utils.OSUtils.GetBallerinaHomePath())
			utils.ToolUtil.AddWritePermissionToFile(distListFilePath)
		}

		if file != nil && file.Mode().Perm()&(1<<1) == 0 {
			panic(utils.ErrorUtil.CreateCommandException("permission denied: you do not have write access to '" + distListFilePath + "'"))
		}

		jsonBytes, err := json.Marshal(distList)
		if err != nil {
			panic(utils.ErrorUtil.CreateCommandException("failed to write in the file: " + err.Error()))
		}

		err = ioutil.WriteFile(distListFilePath, jsonBytes, 0644)
		if err != nil {
			panic(utils.ErrorUtil.CreateCommandException("failed to write in the file: " + err.Error()))
		}
	}

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(exceptions.CommandException); ok {
				panic(r)
			}
		}
	}()

	try()
}

// readLocalDistsFromJson reads the locally available distributions from previously saved json file
func readLocalDistsFromJson(outStream io.Writer, currentBallerinaVersion string) {
	try := func() {
		jsonFile, err := os.Open(utils.OSUtils.GetBallerinaDistListFilePath())
		if err != nil {
			panic(utils.ErrorUtil.CreateCommandException("failed to read the file: " + err.Error()))
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var result map[string]interface{}
		err = json.Unmarshal(byteValue, &result)
		if err != nil {
			panic(utils.ErrorUtil.CreateCommandException("failed to parse the content of the file: " + err.Error()))
		}

		channels := result["channels"].([]interface{})
		for _, channel := range channels {
			channelObj := channel.(map[string]interface{})
			releases := channelObj["releases"].([]interface{})
			for _, release := range releases {
				versionInfo := release.(map[string]interface{})
				fmt.Fprintln(outStream, markVersion(currentBallerinaVersion, versionInfo["version"].(string)))
			}
		}
	}

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(exceptions.CommandException); ok {
				panic(r)
			}
		}
	}()

	try()
}

// listLocalDists lists the locally available distributions from the distributions directory
func listLocalDists(listOfFiles []os.FileInfo, outStream io.Writer, currentBallerinaVersion string) {
	for _, file := range listOfFiles {
		if file.IsDir() {
			version := ""
			parts := strings.Split(file.Name(), "-")
			if len(parts) == 2 {
				version = parts[1]
			}
			fmt.Fprintln(outStream, markVersion(currentBallerinaVersion, version))
		}
	}
}

// getSortedDistList sorts the distributions under a channel according to their semver version
func getSortedDistList(channelDistList []utils.Distribution) []utils.Distribution {
	sort.Slice(channelDistList, func(i, j int) bool {
		versionI := strings.Split(channelDistList[i].Version, ".")
		versionJ := strings.Split(channelDistList[j].Version, ".")

		majorI, _ := strconv.Atoi(versionI[0])
		minorI, _ := strconv.Atoi(versionI[1])
		patchI, _ := strconv.Atoi(versionI[2])

		majorJ, _ := strconv.Atoi(versionJ[0])
		minorJ, _ := strconv.Atoi(versionJ[1])
		patchJ, _ := strconv.Atoi(versionJ[2])

		if majorI != majorJ {
			return majorI > majorJ
		}
		if minorI != minorJ {
			return minorI > minorJ
		}
		return patchI > patchJ
	})

	return channelDistList
}

// isUpdated checks if the distributions list file is updated
func isUpdated(listOfFiles []os.FileInfo) bool {
	distListFilePath := utils.OSUtils.GetBallerinaDistListFilePath()
	if _, err := os.Stat(distListFilePath); os.IsNotExist(err) {
		return false
	}

	jsonFile, err := os.Open(distListFilePath)
	if err != nil {
		panic(utils.ErrorUtil.CreateCommandException("failed to read the file: " + err.Error()))
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		panic(utils.ErrorUtil.CreateCommandException("failed to parse the content of the file: " + err.Error()))
	}

	jsonString := string(byteValue)
	for _, file := range listOfFiles {
		if file.IsDir() {
			version := ""
			parts := strings.Split(file.Name(), "-")
			if len(parts) == 2 {
				version = parts[1]
			}
			if !strings.Contains(jsonString, version) {
				return false
			}
		}
	}
	return true
}
