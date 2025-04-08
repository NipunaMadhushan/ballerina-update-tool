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

package utils

import (
	"archive/zip"
	"ballerina-update-tool/constants"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Constants for server URLs and version info
const (
	ConnectionErrorMessage  = "connection to the remote server failed"
	ProxyErrorMessage       = "connection to the remote server through proxy server failed"
	DefaultBallerinaVersion = "0.0.0"
)

// Environment variables
var (
	BallerinaStagingUpdate = os.Getenv("BALLERINA_STAGING_UPDATE") == "true"
	BallerinaDevUpdate     = os.Getenv("BALLERINA_DEV_UPDATE") == "true"
	TestMode               = os.Getenv("TEST_MODE_ACTIVE") == "true"
)

// ToolUtil provides utility functions for Ballerina tools, corresponding to ToolUtil.java
var ToolUtil = toolUtilStruct{}

type toolUtilStruct struct{}

// Helper function to check if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !info.IsDir()
}

// Helper function to check if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return info.IsDir()
}

// Helper function to check if a file is writable
func isWritable(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			dir := filepath.Dir(path)
			return isWritable(dir)
		}
		return false
	}
	if fileInfo.IsDir() {
		tempFile := filepath.Join(path, fmt.Sprintf(".write_test_%d", time.Now().UnixNano()))
		file, err := os.OpenFile(tempFile, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
		if err != nil {
			return false
		}
		file.Close()
		os.Remove(tempFile)
		return true
	}
	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// GetCurrentBallerinaVersion returns the currently used Ballerina version
func (t toolUtilStruct) GetCurrentBallerinaVersion() string {
	installerVersionFilePath := OSUtils.GetInstallerVersionFilePath()

	if fileExists(installerVersionFilePath) {
		installedInstallerVersion, err := t.getInstallerVersion(OSUtils.GetInstalledInstallerVersionPath())
		if err != nil {
			panic(ErrorUtil.CreateCommandException("current Ballerina version not found: " + err.Error()))
		}

		currentInstallerVersion, err := t.getInstallerVersion(installerVersionFilePath)
		if err != nil {
			panic(ErrorUtil.CreateCommandException("current Ballerina version not found: " + err.Error()))
		}

		if installedInstallerVersion != currentInstallerVersion {
			t.SetCurrentBallerinaVersion(t.GetCurrentInstalledBallerinaVersion())
			t.setInstallerVersion(installerVersionFilePath)
		}
	}

	ballerinaVersionFilePath := OSUtils.GetBallerinaVersionFilePath()
	if !fileExists(ballerinaVersionFilePath) {
		defaultBallerinaVersion := DefaultBallerinaVersion
		t.SetCurrentBallerinaVersion(defaultBallerinaVersion)
		t.setInstallerVersion(installerVersionFilePath)
	}

	userVersion, err := t.getVersion(ballerinaVersionFilePath)
	if err != nil {
		panic(ErrorUtil.CreateCommandException("current Ballerina version not found: " + err.Error()))
	}

	if t.CheckDistributionAvailable(userVersion) {
		return userVersion
	}

	return t.GetCurrentInstalledBallerinaVersion()
}

// GetCurrentInstalledBallerinaVersion returns the installed Ballerina version
func (t toolUtilStruct) GetCurrentInstalledBallerinaVersion() string {
	version, err := t.getVersion(OSUtils.GetInstalledConfigPath())
	if err != nil {
		// If files do not exist, return empty version
		return DefaultBallerinaVersion
	}
	return version
}

// SetCurrentBallerinaVersion sets the current Ballerina version
func (t toolUtilStruct) SetCurrentBallerinaVersion(version string) {
	ballerinaVersionFilePath := OSUtils.GetBallerinaVersionFilePath()
	if !isWritable(ballerinaVersionFilePath) {
		panic(ErrorUtil.CreateCommandException(fmt.Sprintf("permission denied: you do not have write access to '%s'",
			ballerinaVersionFilePath)))
	}

	err := t.setVersion(ballerinaVersionFilePath, version)
	if err != nil {
		panic(ErrorUtil.CreateCommandException("failed to set the Ballerina version: " + err.Error()))
	}
}

// ClearCache clears the BIR and JAR caches
func (t toolUtilStruct) ClearCache(outStream io.Writer) {
	err := OSUtils.ClearBirCacheLocation(outStream)
	if err != nil {
		panic(ErrorUtil.CreateCommandException("failed to clear the caches."))
		return
	}

	err = OSUtils.ClearJarCacheLocation(outStream)
	if err != nil {
		panic(ErrorUtil.CreateCommandException("failed to clear the caches."))
	}
}

// GetCurrentToolsVersion reads the version from properties file
func (t *toolUtilStruct) GetCurrentToolsVersion() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	// Find resources directory next to the executable
	propsPath := filepath.Join(filepath.Dir(filepath.Dir(execPath)), "resources", "tool.properties")

	data, err := os.ReadFile(propsPath)
	if err != nil {
		panic(ErrorUtil.CreateCommandException("version info not available"))
	}

	// Parse properties file content
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "command.version=") {
			return strings.TrimPrefix(line, "command.version=")
		}
	}

	panic(ErrorUtil.CreateCommandException("version info not available"))
}

// Helper functions for file operations
func (t toolUtilStruct) getVersion(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("empty version file")
	}

	parts := strings.Split(lines[0], "-")
	if len(parts) == 2 {
		return parts[1], nil
	}

	return "", nil
}

func (t toolUtilStruct) getInstallerVersion(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("empty installer version file")
	}

	return lines[0], nil
}

func (t toolUtilStruct) setVersion(path string, version string) error {
	versionType := t.GetType(version)
	content := fmt.Sprintf("%s-%s\n", versionType, version)
	return ioutil.WriteFile(path, []byte(content), 0644)
}

func (t toolUtilStruct) setInstallerVersion(path string) error {
	installerVersion, err := t.getInstallerVersion(OSUtils.GetInstalledInstallerVersionPath())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, []byte(installerVersion+"\n"), 0644)
}

// UseBallerinaVersion sets the current Ballerina version and clears cache
func (t toolUtilStruct) UseBallerinaVersion(printStream io.Writer, distribution string) {
	t.SetCurrentBallerinaVersion(distribution)
	t.ClearCache(printStream)
}

// CheckDistributionAvailable checks if a distribution is available locally
func (t toolUtilStruct) CheckDistributionAvailable(distribution string) bool {
	installFile := filepath.Join(t.GetDistributionsPath(), t.GetType(distribution)+"-"+distribution)
	return dirExists(installFile)
}

// CheckDependencyAvailable checks if a dependency is available locally
func (t toolUtilStruct) CheckDependencyAvailable(dependency string) bool {
	dependencyLocation := filepath.Join(t.GetDependencyPath(), dependency)
	return dirExists(dependencyLocation)
}

// GetServerURLWithProxy creates an HTTP client with proxy configuration
func (t toolUtilStruct) GetServerURLWithProxy() (*http.Client, error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	if t.CheckProxyConfigsDefinition() {
		proxyConfigs, err := t.GetProxyConfigs()
		if err != nil {
			return client, err
		}

		proxyHost, _ := proxyConfigs["host"].(string)
		proxyPortStr, _ := proxyConfigs["port"].(string)
		proxyUser, _ := proxyConfigs["user"].(string)
		proxyPassword, _ := proxyConfigs["password"].(string)

		if proxyHost != "" && proxyPortStr != "" {
			proxyPort, err := strconv.Atoi(proxyPortStr)
			if err == nil && proxyPort > 0 && proxyPort < 65536 {
				proxyURL := &url.URL{
					Scheme: "http",
					Host:   fmt.Sprintf("%s:%d", proxyHost, proxyPort),
				}

				if proxyUser != "" && proxyPassword != "" {
					proxyURL.User = url.UserPassword(proxyUser, proxyPassword)
				}

				transport := &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				}

				client.Transport = transport
			}
		}
	}

	return client, nil
}

// GetProxyConfigs reads and returns proxy configurations
func (t toolUtilStruct) GetProxyConfigs() (map[string]interface{}, error) {
	settingsFile := filepath.Join(OSUtils.GetBallerinaHomePath(), constants.CommandToolConstants.BallerinaSettingsFile)

	// Simple TOML parser for proxy settings - minimal implementation
	content, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return nil, err
	}

	proxyConfig := make(map[string]interface{})
	inProxySection := false
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "[proxy]" {
			inProxySection = true
			continue
		} else if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inProxySection = false
			continue
		}

		if inProxySection && line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Remove quotes if present
				value = strings.Trim(value, "\"")
				proxyConfig[key] = value
			}
		}
	}

	return proxyConfig, nil
}

// CheckProxyConfigsDefinition checks if proxy configs are defined
func (t toolUtilStruct) CheckProxyConfigsDefinition() bool {
	settingsFile := filepath.Join(OSUtils.GetBallerinaHomePath(), constants.CommandToolConstants.BallerinaSettingsFile)
	if !fileExists(settingsFile) {
		return false
	}

	content, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return false
	}

	// Simple check for [proxy] section
	return strings.Contains(string(content), "[proxy]")
}

// GetDistributions gets available distributions from the server
func (t toolUtilStruct) GetDistributions(printStream io.Writer) []Channel {
	client, err := t.GetServerURLWithProxy()
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}

	url := t.GetServerURL() + "/distributions"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}

	req.Header.Set("user-agent", OSUtils.GetUserAgent(t.GetCurrentBallerinaVersion(),
		t.GetCurrentToolsVersion(), "jballerina"))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(ErrorUtil.CreateCommandException(t.getServerRequestFailedErrorMessage(resp)))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}

	// In Go, we'll parse the JSON directly rather than using regex
	var jsonResp struct {
		Distributions []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Type    string `json:"type"`
			Channel string `json:"channel"`
		} `json:"list"`
	}

	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}

	channels := make(map[string]*Channel)
	distributions := []Distribution{}

	// Process distributions
	for _, d := range jsonResp.Distributions {
		dist := Distribution{
			Name:    d.Name,
			Version: d.Version,
			Type:    d.Type,
			Channel: d.Channel,
		}
		distributions = append(distributions, dist)
	}

	// Group distributions by channel
	for _, distribution := range distributions {
		channelName := distribution.Channel
		var channel *Channel

		if existingChannel, ok := channels[channelName]; ok {
			channel = existingChannel
		} else {
			newChannel := NewChannelWithName(channelName)
			channel = newChannel
			channels[channelName] = channel
		}

		// Add to beginning of distributions list
		channel.Distributions = append([]Distribution{distribution}, channel.Distributions...)
	}

	// Convert map to slice
	result := []Channel{}
	for _, ch := range channels {
		result = append([]Channel{*ch}, result...) // Add to beginning
	}

	return result
}

// GetLatest gets the latest version for a given distribution type
func (t toolUtilStruct) GetLatest(currentVersion string, distType string) string {
	client, err := t.GetServerURLWithProxy()
	if err != nil {
		return ""
	}

	url := fmt.Sprintf("%s/distributions/latest?version=%s&type=%s",
		t.GetServerURL(), currentVersion, distType)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("user-agent", OSUtils.GetUserAgent(t.GetCurrentBallerinaVersion(),
		t.GetCurrentToolsVersion(), "jballerina"))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return t.getValue(distType, t.convertStreamToString(resp.Body))
	}

	if resp.StatusCode == 404 {
		return ""
	}

	panic(ErrorUtil.CreateCommandException(t.getServerRequestFailedErrorMessage(resp)))
}

// Helper to extract a value from JSON
func (t toolUtilStruct) getValue(key, json string) string {
	pattern := fmt.Sprintf(`"%s":"(.*?)"`, key)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(json)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// GetLatestToolVersion gets the latest tool version
func (t toolUtilStruct) GetLatestToolVersion() *Tool {
	client, err := t.GetServerURLWithProxy()
	if err != nil {
		panic(ErrorUtil.CreateCommandException(err.Error()))
	}

	url := t.GetServerURL() + "/versions/latest"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(ErrorUtil.CreateCommandException(err.Error()))
	}

	req.Header.Set("user-agent", OSUtils.GetUserAgent(t.GetCurrentBallerinaVersion(),
		t.GetCurrentToolsVersion(), "jballerina"))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		tool := NewTool()
		json := t.convertStreamToString(resp.Body)

		// Extract version
		versionPattern := regexp.MustCompile(`"version":"(.*?)"`)
		versionMatches := versionPattern.FindStringSubmatch(json)
		if len(versionMatches) > 1 {
			tool.SetVersion(versionMatches[1])
		}

		// Extract compatibility
		compatPattern := regexp.MustCompile(`"compatibility":(true|false)`)
		compatMatches := compatPattern.FindStringSubmatch(json)
		if len(compatMatches) > 1 {
			tool.SetCompatibility(compatMatches[1])
		}

		return tool
	}

	if resp.StatusCode == 404 {
		return nil
	}

	panic(ErrorUtil.CreateCommandException(t.getServerRequestFailedErrorMessage(resp)))
}

// Helper to convert a reader to a string
func (t toolUtilStruct) convertStreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}

// GetDistributionsPath returns the path to distributions
func (t toolUtilStruct) GetDistributionsPath() string {
	installPath, err := OSUtils.GetInstallationPath()
	if err != nil {
		panic(ErrorUtil.CreateCommandException("failed to get the path of the distributions"))
	}

	distDirectory := filepath.Join(installPath, "distributions")
	if !dirExists(distDirectory) {
		err := os.MkdirAll(distDirectory, 0755)
		if err != nil {
			panic(ErrorUtil.CreateCommandException("failed to create distributions directory"))
		}
	}

	return distDirectory
}

// GetDependencyPath returns the path to dependencies
func (t toolUtilStruct) GetDependencyPath() string {
	installPath, err := OSUtils.GetInstallationPath()
	if err != nil {
		panic(ErrorUtil.CreateCommandException("failed to get the path of the distributions"))
	}

	depDirectory := filepath.Join(installPath, "dependencies")
	if !dirExists(depDirectory) {
		err := os.MkdirAll(depDirectory, 0755)
		if err != nil {
			panic(ErrorUtil.CreateCommandException("failed to create dependencies directory"))
		}
	}

	return depDirectory
}

// GetToolUnzipLocation returns the path where the tool is unzipped
func (t toolUtilStruct) GetToolUnzipLocation() string {
	installPath, err := OSUtils.GetInstallationPath()
	if err != nil {
		panic(ErrorUtil.CreateCommandException(
			"failed to get a temporary directory to unzip the update tool zip to"))
	}

	return filepath.Join(installPath, "ballerina-command-tmp")
}

// CheckForUpdate checks if an update is available
func (t toolUtilStruct) CheckForUpdate(printStream io.Writer) {
	defer func() {
		// Recover from any panics to make this function optional
		if r := recover(); r != nil {
			// Just ignore any errors
		}
	}()

	version := t.GetCurrentBallerinaVersion()
	if OSUtils.UpdateNotice() {
		latestVersion := t.GetLatest(version, "patch")
		// For 1.0.x releases we support through jballerina distribution
		if latestVersion == "" || strings.HasPrefix(latestVersion, constants.CommandToolConstants.Ballerina1XVersions) {
			return
		}
		if latestVersion != version {
			fmt.Fprintf(printStream, "A new version of Ballerina is available: %s\n", latestVersion)
			fmt.Fprintf(printStream, "Use 'bal dist pull %s' to download and use the distribution\n\n", latestVersion)
		}
	}
}

// DownloadDistribution downloads a distribution
func (t toolUtilStruct) DownloadDistribution(printStream io.Writer, distribution, distributionType, distributionVersion string, testMode bool) bool {
	if !t.CheckDistributionAvailable(distribution) {
		client, err := t.GetServerURLWithProxy()
		if err != nil {
			panic(ErrorUtil.CreateCommandException(err.Error()))
		}
		// Set a longer timeout for downloading distributions
		client.Timeout = time.Minute * 10

		url := t.GetServerURL() + "/distributions/" + distributionVersion
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(ErrorUtil.CreateCommandException(err.Error()))
		}

		req.Header.Set("user-agent", OSUtils.GetUserAgent(t.GetCurrentBallerinaVersion(),
			t.GetCurrentToolsVersion(), distributionType))
		req.Header.Set("Accept", "application/json")

		if testMode || TestMode {
			req.Header.Set("testMode", "true")
		}

		resp, err := client.Do(req)
		if err != nil {
			if t.CheckProxyConfigsDefinition() {
				panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
			} else {
				panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode == 302 {
			fmt.Fprintf(printStream, "Fetching the '%s' distribution from the remote server...\n", distribution)
			newUrl := resp.Header.Get("Location")
			redirectReq, err := http.NewRequest("GET", newUrl, nil)
			if err != nil {
				panic(ErrorUtil.CreateCommandException(err.Error()))
			}

			redirectReq.Header.Set("content-type", "binary/data")
			redirectResp, err := client.Do(redirectReq)
			if err != nil {
				panic(ErrorUtil.CreateCommandException(err.Error()))
			}
			defer redirectResp.Body.Close()

			t.DownloadDistributionZip(printStream, redirectResp, distribution)
			dependencyForDistribution := t.GetDependency(printStream, distribution, distributionType, distributionVersion)
			if !t.CheckDependencyAvailable(dependencyForDistribution) {
				t.DownloadDependency(printStream, dependencyForDistribution, distributionType, distributionVersion)
			}
			t.SetupDistribution(distribution, dependencyForDistribution)
			return false
		} else if resp.StatusCode == 200 {
			fmt.Fprintf(printStream, "Fetching the '%s' distribution from the remote server...\n", distribution)
			t.DownloadDistributionZip(printStream, resp, distribution)
			dependencyForDistribution := t.GetDependency(printStream, distribution, distributionType, distributionVersion)
			if !t.CheckDependencyAvailable(dependencyForDistribution) {
				t.DownloadDependency(printStream, dependencyForDistribution, distributionType, distributionVersion)
			}
			t.SetupDistribution(distribution, dependencyForDistribution)
			return false
		} else {
			panic(ErrorUtil.CreateDistributionNotFoundException(distribution))
		}
	} else {
		fmt.Fprintf(printStream, "'%s' is already available locally\n", distribution)
		return true
	}
}

// DownloadDistributionZip downloads the distribution zip file
func (t toolUtilStruct) DownloadDistributionZip(printStream io.Writer, resp *http.Response, distribution string) {
	zipFileLocation := filepath.Join(t.GetDistributionsPath(), distribution+".zip")
	err := t.downloadFile(resp, zipFileLocation, distribution, printStream)
	if err != nil {
		panic(ErrorUtil.CreateCommandException("failed to download distribution: " + err.Error()))
	}
}

// SetupDistribution sets up the downloaded distribution
func (t toolUtilStruct) SetupDistribution(distribution string, dependency string) {
	distPath := t.GetDistributionsPath()
	zipFileLocation := filepath.Join(t.GetDistributionsPath(), distribution+".zip")

	if !t.CheckDependencyAvailable(dependency) {
		os.Remove(zipFileLocation)
		panic(ErrorUtil.CreateCommandException(fmt.Sprintf(
			"The required dependency '%s' is not available locally. Please try reinstalling the distribution.",
			dependency)))
	}

	t.unzip(zipFileLocation, distPath)

	// Add executable permissions to distribution binary
	execFileName := OSUtils.GetExecutableFileName(distribution)
	execFilePath := filepath.Join(distPath, t.GetType(distribution)+"-"+distribution, "bin", execFileName)
	t.addExecutablePermissionToFile(execFilePath)

	// Add executable permissions to language server and debug adapter if they exist
	langServerPath := filepath.Join(distPath, distribution, "lib", "tools")
	launcherServer := filepath.Join(langServerPath, "lang-server", "launcher", OSUtils.GetLangServerLauncherName())
	debugAdapter := filepath.Join(langServerPath, "debug-adapter", "launcher", OSUtils.GetDebugAdapterName())

	if fileExists(debugAdapter) {
		t.addExecutablePermissionToFile(debugAdapter)
	}

	if fileExists(launcherServer) {
		t.addExecutablePermissionToFile(launcherServer)
	}

	// Remove the zip file
	os.Remove(zipFileLocation)
}

// GetDependency gets the dependency for a distribution
func (t toolUtilStruct) GetDependency(printStream io.Writer, distribution, distributionType, distributionVersion string) string {
	client, err := t.GetServerURLWithProxy()
	if err != nil {
		fmt.Fprintln(printStream, "Error setting up HTTP client:", err)
		return ""
	}

	url := t.GetServerURL() + "/distributions"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintln(printStream, "Error creating request:", err)
		return ""
	}

	req.Header.Set("user-agent", OSUtils.GetUserAgent(distributionVersion,
		t.GetCurrentToolsVersion(), distributionType))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		responseStr := t.convertStreamToString(resp.Body)

		// Find the distribution info containing our version
		infoRegex := regexp.MustCompile(`\{.*?\]\}`)
		infoMatches := infoRegex.FindAllString(responseStr, -1)

		for _, distInfo := range infoMatches {
			if strings.Contains(distInfo, distributionVersion) {
				// Extract the dependency from this distribution
				dependencyRegex := regexp.MustCompile(`"(jdk.*?)"`)
				dependencyMatches := dependencyRegex.FindStringSubmatch(distInfo)
				if len(dependencyMatches) > 1 {
					return dependencyMatches[1]
				}
			}
		}
	}

	return ""
}

// DownloadDependency downloads a dependency
func (t toolUtilStruct) DownloadDependency(printStream io.Writer, dependency, distributionType, distributionVersion string) {
	fmt.Fprintf(printStream, "\nFetching the dependencies for '%s' from the remote server...\n", distributionVersion)

	encodedDependencyName := t.encodePlusCharacters(dependency)
	url := t.GetServerURL() + "/dependencies/" + encodedDependencyName

	client, err := t.GetServerURLWithProxy()
	if err != nil {
		panic(ErrorUtil.CreateCommandException(err.Error()))
	}
	// Set a longer timeout for downloading dependencies
	client.Timeout = time.Minute * 5

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(ErrorUtil.CreateCommandException(err.Error()))
	}

	req.Header.Set("user-agent", OSUtils.GetUserAgent(distributionVersion,
		t.GetCurrentToolsVersion(), distributionType))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		newUrl := resp.Header.Get("Location")
		redirectReq, err := http.NewRequest("GET", newUrl, nil)
		if err != nil {
			panic(ErrorUtil.CreateCommandException(err.Error()))
		}

		redirectReq.Header.Set("content-type", "binary/data")
		redirectResp, err := client.Do(redirectReq)
		if err != nil {
			panic(ErrorUtil.CreateCommandException(err.Error()))
		}
		defer redirectResp.Body.Close()

		t.downloadAndSetupDependency(redirectResp, printStream, dependency)
	} else if resp.StatusCode == 200 {
		t.downloadAndSetupDependency(resp, printStream, dependency)
	} else {
		panic(ErrorUtil.CreateDependencyNotFoundException(dependency))
	}
}

// downloadAndSetupDependency downloads and sets up a dependency
func (t toolUtilStruct) downloadAndSetupDependency(resp *http.Response, printStream io.Writer, dependency string) {
	dependencyLocation := t.GetDependencyPath()
	zipFileLocation := filepath.Join(dependencyLocation, dependency+".zip")

	// Download the dependency
	err := t.downloadFile(resp, zipFileLocation, dependency, printStream)
	if err != nil {
		panic(ErrorUtil.CreateCommandException("Failed to download dependency: " + err.Error()))
	}

	// Unzip the dependency
	t.unzip(zipFileLocation, dependencyLocation)

	// Add executable permissions to the dependency directory
	t.AddExecutablePermissionToDirectory(filepath.Join(dependencyLocation, dependency))

	// Delete the zip file
	if fileExists(zipFileLocation) {
		os.Remove(zipFileLocation)
	}
}

// RemoveUnusedDependencies removes dependencies that are no longer used
func (t toolUtilStruct) RemoveUnusedDependencies(distributionVersion string, printStream io.Writer) {
	dependencyForDistribution := ""
	channelForDistribution := ""
	distributionsWithDependency := []string{}

	channels := t.GetDistributions(printStream)

	// Find the dependency for the specified distribution
	for _, channel := range channels {
		for _, distribution := range channel.Distributions {
			if distribution.Version == distributionVersion {
				dependencyForDistribution = distribution.Dependency
				channelForDistribution = channel.Name
				break
			}
		}
		if dependencyForDistribution != "" {
			break
		}
	}

	if dependencyForDistribution == "" {
		fmt.Fprintln(printStream, "No dependency found for the given distribution version")
		return
	}

	// Find all distributions that use this dependency
	for _, channel := range channels {
		if channel.Name == channelForDistribution {
			for _, distribution := range channel.Distributions {
				if distribution.Dependency == dependencyForDistribution {
					distributionsWithDependency = append(distributionsWithDependency, distribution.Version)
				}
			}
		}
	}

	// Get local distributions
	localDistributions := t.getLocalDistributions()

	// Find the intersection of local distributions and distributions with dependency
	inUse := false
	for _, local := range localDistributions {
		for _, withDep := range distributionsWithDependency {
			if local == withDep {
				inUse = true
				break
			}
		}
		if inUse {
			break
		}
	}

	// If no distributions with this dependency are installed locally, delete the dependency
	if !inUse {
		fmt.Fprintf(printStream, "No local distributions found for the dependency '%s'\n"+
			"Deleting the dependency '%s'\n", dependencyForDistribution, dependencyForDistribution)
		dependencyToDelete := filepath.Join(t.GetDependencyPath(), dependencyForDistribution)
		if dirExists(dependencyToDelete) {
			err := OSUtils.DeleteFiles(dependencyToDelete)
			if err != nil {
				fmt.Fprintf(printStream, "Error occurred while deleting the dependency '%s': %v\n",
					dependencyForDistribution, err)
			}
		}
	}
}

// getLocalDistributions gets the list of locally installed distributions
func (t toolUtilStruct) getLocalDistributions() []string {
	localDistributions := []string{}
	folder := t.GetDistributionsPath()

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return localDistributions
	}

	// Sort files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() {
			parts := strings.Split(file.Name(), "-")
			if len(parts) == 2 {
				localDistributions = append(localDistributions, parts[1])
			}
		}
	}

	return localDistributions
}

// DownloadTool downloads a tool
func (t toolUtilStruct) DownloadTool(printStream io.Writer, toolVersion string) {
	url := t.GetServerURL() + "/versions/" + toolVersion

	client, err := t.GetServerURLWithProxy()
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}
	// Set a longer timeout for downloading tools
	client.Timeout = time.Minute * 2

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}

	req.Header.Set("user-agent", OSUtils.GetUserAgent(t.GetCurrentBallerinaVersion(),
		t.GetCurrentToolsVersion(), "jballerina"))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if t.CheckProxyConfigsDefinition() {
			panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
		} else {
			panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		newUrl := resp.Header.Get("Location")
		redirectReq, err := http.NewRequest("GET", newUrl, nil)
		if err != nil {
			if t.CheckProxyConfigsDefinition() {
				panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
			} else {
				panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
			}
		}

		redirectReq.Header.Set("content-type", "binary/data")
		redirectResp, err := client.Do(redirectReq)
		if err != nil {
			if t.CheckProxyConfigsDefinition() {
				panic(ErrorUtil.CreateCommandException(ProxyErrorMessage))
			} else {
				panic(ErrorUtil.CreateCommandException(ConnectionErrorMessage))
			}
		}
		defer redirectResp.Body.Close()

		t.downloadAndSetupTool(printStream, redirectResp, "ballerina-command-"+toolVersion)
	} else if resp.StatusCode == 200 {
		t.downloadAndSetupTool(printStream, resp, "ballerina-command-"+toolVersion)
	} else {
		panic(ErrorUtil.CreateCommandException(fmt.Sprintf("tool version '%s' not found", toolVersion)))
	}
}

// downloadAndSetupTool downloads and sets up a tool
func (t toolUtilStruct) downloadAndSetupTool(printStream io.Writer, resp *http.Response, toolFileName string) {
	toolUnzipLocation := t.GetToolUnzipLocation()

	// Create/clear the unzip directory
	if dirExists(toolUnzipLocation) {
		os.RemoveAll(toolUnzipLocation)
	}
	os.MkdirAll(toolUnzipLocation, 0755)

	zipFileLocation := filepath.Join(toolUnzipLocation, toolFileName+".zip")

	// Download the tool
	err := t.downloadFile(resp, zipFileLocation, toolFileName, printStream)
	if err != nil {
		os.RemoveAll(toolUnzipLocation)
		if fileExists(zipFileLocation) {
			os.Remove(zipFileLocation)
		}
		panic(ErrorUtil.CreateCommandException("Failed to download tool: " + err.Error()))
	}

	// Unzip the tool
	t.unzip(zipFileLocation, toolUnzipLocation)

	// Copy and setup scripts
	t.copyScripts(toolUnzipLocation, toolFileName)

	// Clean up
	if fileExists(zipFileLocation) {
		os.Remove(zipFileLocation)
	}
}

// downloadFile downloads a file from an HTTP response to a local path with in-place progress updates
func (t toolUtilStruct) downloadFile(resp *http.Response, zipFileLocation, fileName string, printStream io.Writer) error {
	// Create the file
	out, err := os.Create(zipFileLocation)
	if err != nil {
		return fmt.Errorf("failed to download file %s to %s: %v", fileName, zipFileLocation, err)
	}
	defer out.Close()

	// Get file size for progress reporting
	totalSize := resp.ContentLength
	downloaded := int64(0)
	totalSizeInMB := float64(totalSize) / (1024 * 1024)

	fmt.Fprintf(printStream, "Downloading %s (%.2f MB)...\n", fileName, totalSizeInMB)

	// Create a buffer to read into
	buf := make([]byte, 32*1024) // 32KB buffer
	lastPercentReported := -1

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			// Write to file
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}

			// Update progress
			downloaded += int64(n)
			if totalSize > 0 {
				percentComplete := int((float64(downloaded) / float64(totalSize)) * 100)

				// Update progress in-place if it has changed
				if percentComplete != lastPercentReported {
					fmt.Fprintf(printStream, "\rProgress: %d%%", percentComplete)
					lastPercentReported = percentComplete
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				// Print a newline before the final message since we've been using carriage returns
				fmt.Fprintln(printStream, "\nDownload complete")
				break
			}
			return err
		}
	}

	return nil
}

// unzip extracts a zip file to a destination directory
func (t toolUtilStruct) unzip(zipFilePath string, destDirectory string) {
	// Ensure destination directory exists
	if !dirExists(destDirectory) {
		os.MkdirAll(destDirectory, 0755)
	}

	// Open the zip file
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		panic(ErrorUtil.CreateCommandException(
			"failed to unzip the zip file in '" + zipFilePath + "' to '" + destDirectory + "'"))
	}
	defer zipReader.Close()

	// Extract each file
	for _, file := range zipReader.File {
		filePath := filepath.Join(destDirectory, file.Name)

		// Create directory for file if needed
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}

		// Ensure parent directory exists
		os.MkdirAll(filepath.Dir(filePath), 0755)

		// Create the file
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			panic(ErrorUtil.CreateCommandException("failed to create output file: " + err.Error()))
		}

		// Open the file in the archive
		inFile, err := file.Open()
		if err != nil {
			outFile.Close()
			panic(ErrorUtil.CreateCommandException("failed to open file in zip: " + err.Error()))
		}

		// Copy the file contents
		_, err = io.Copy(outFile, inFile)
		outFile.Close()
		inFile.Close()
		if err != nil {
			panic(ErrorUtil.CreateCommandException("failed to extract file: " + err.Error()))
		}
	}
}

// copyScripts copies install scripts from the unzipped tool package
func (t toolUtilStruct) copyScripts(unzippedUpdateToolPath string, ballerinaCommandDir string) {
	installScriptFileName := OSUtils.GetInstallScriptFileName()
	installScriptDest := filepath.Join(unzippedUpdateToolPath, installScriptFileName)
	installScriptSrc := filepath.Join(unzippedUpdateToolPath, ballerinaCommandDir, "scripts", installScriptFileName)

	// Copy the install script
	input, err := ioutil.ReadFile(installScriptSrc)
	if err != nil {
		panic(ErrorUtil.CreateCommandException(fmt.Sprintf(
			"failed to copy the update scripts to temporary directory '%s': %v",
			unzippedUpdateToolPath, err)))
	}

	err = ioutil.WriteFile(installScriptDest, input, 0755)
	if err != nil {
		panic(ErrorUtil.CreateCommandException(fmt.Sprintf(
			"failed to copy the update scripts to temporary directory '%s': %v",
			unzippedUpdateToolPath, err)))
	}

	// Add executable permissions
	t.addExecutablePermissionToFile(installScriptDest)
}

// addExecutablePermissionToFile adds executable permissions to a file
func (t toolUtilStruct) addExecutablePermissionToFile(path string) {
	// Use chmod to set permissions
	if fileExists(path) {
		os.Chmod(path, 0755) // read+write+execute for owner, read+execute for group and others
	}
}

// AddWritePermissionToFile adds write permissions to a file
func (t toolUtilStruct) AddWritePermissionToFile(path string) {
	// Use chmod to set permissions
	if fileExists(path) {
		os.Chmod(path, 0644) // read+write for owner, read for group and others
	}
}

// AddExecutablePermissionToDirectory adds executable permissions to a directory recursively
func (t toolUtilStruct) AddExecutablePermissionToDirectory(path string) {
	if !dirExists(path) {
		return
	}

	// On Unix-like systems, use the chmod command for recursive permission changes
	if OSUtils.IsWindows() {
		// On Windows, walk the directory and set permissions manually
		filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				os.Chmod(filePath, 0755)
			} else {
				os.Chmod(filePath, 0755)
			}
			return nil
		})
	} else {
		// On Unix-like systems, use chmod -R
		cmd := exec.Command("chmod", "-R", "755", path)
		err := cmd.Run()
		if err != nil {
			panic(ErrorUtil.CreateCommandException("permission denied: you do not have write access to '" + path + "'"))
		}
	}
}

// ReadFileAsString reads a file and returns its contents as a string
func (t toolUtilStruct) ReadFileAsString(path string) (string, error) {
	// In Go, we'd typically use the embedded files or specific file paths
	// This implementation is different from Java's ClassLoader.getSystemResourceAsStream

	// Try to find the file in the current directory
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("path cannot be found in %s: %v", path, err)
	}

	return string(data), nil
}

// HandleInstallDirPermission checks if the installation directory is writable
func (t toolUtilStruct) HandleInstallDirPermission() {
	installationPath, err := OSUtils.GetInstallationPath()
	if err != nil {
		panic(ErrorUtil.CreateCommandException("failed to get the path to the Ballerina installation directory"))
	}

	if !isWritable(installationPath) {
		panic(ErrorUtil.CreateCommandException(fmt.Sprintf(
			"permission denied: you do not have write access to '%s'", installationPath)))
	}
}

// GetServerURL returns the server URL based on environment variables
func (t toolUtilStruct) GetServerURL() string {
	url := constants.CommandToolConstants.ProductionUrl
	if BallerinaStagingUpdate {
		url = constants.CommandToolConstants.StagingUrl
	}
	if BallerinaDevUpdate {
		url = constants.CommandToolConstants.DevUrl
	}
	return url
}

// getServerRequestFailedErrorMessage creates an error message for failed server requests
func (t toolUtilStruct) getServerRequestFailedErrorMessage(resp *http.Response) string {
	responseMessage := resp.Status
	if responseMessage == "" {
		responseMessage = fmt.Sprintf("%d", resp.StatusCode)
	}
	return "server request failed: " + responseMessage
}

// GetType returns the distribution type based on version
func (t toolUtilStruct) GetType(version string) string {
	parts := strings.Split(version, ".")
	if len(parts) > 0 && parts[0] == "1" {
		return "jballerina"
	}
	return "ballerina"
}

// GetTypeName returns a formatted type name based on version
func (t toolUtilStruct) GetTypeName(version string) string {
	if len(version) == 0 {
		return ""
	}

	lastChar := version[len(version)-1]

	if strings.Contains(version, "1.") {
		return "jballerina" + " version " + version
	} else if strings.Contains(version, "slp") {
		return fmt.Sprintf(" Preview %c", lastChar)
	} else {
		if len(version) <= 2 {
			return version
		}
		versionID := version[2:len(version)-1] + " " + string(lastChar)
		return strings.ToUpper(versionID[:1]) + versionID[1:]
	}
}

// EncodePlusCharacters encodes + characters in a dependency name
func (t toolUtilStruct) encodePlusCharacters(dependency string) string {
	encodedDependency := url.QueryEscape(dependency)
	encodedDependency = strings.ReplaceAll(encodedDependency, "+", "%2B")
	return encodedDependency
}

// UpdateTool updates the update tool if any latest version available
func (t *toolUtilStruct) UpdateTool(printStream io.Writer) Tool {
	version := t.GetCurrentToolsVersion()
	toolDetails := Tool{}

	fmt.Fprintln(printStream, "Checking for newer versions of the update tool...")

	// Use deferred function to handle panics
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Just recover from the panic, we'll return toolDetails with empty values
				fmt.Fprintln(printStream, "Error checking for tool updates")
			}
		}()

		latestVersionInfo := t.GetLatestToolVersion()

		latestVersion := ""
		backwardCompatibility := ""
		if latestVersionInfo != nil {
			latestVersion = latestVersionInfo.Version
			backwardCompatibility = latestVersionInfo.Compatibility
		}

		toolDetails.Version = latestVersion
		toolDetails.Compatibility = backwardCompatibility

		if latestVersion == "" {
			fmt.Fprintln(printStream, "Failed to find the latest update tool version")
		} else if latestVersion != version {
			if backwardCompatibility == "" {
				fmt.Fprintln(printStream, "Failed to find the compatibility of the latest update tool version")
			} else if backwardCompatibility != "true" {
				fmt.Fprintln(printStream)
				fmt.Fprintln(printStream, "ERROR: Outdated Ballerina update tool version found")
				fmt.Fprintln(printStream, "Use the following command to update the Ballerina update tool")
				fmt.Fprintln(printStream, "   bal update")
				fmt.Fprintln(printStream)
			} else {
				t.HandleInstallDirPermission()

				// Download tool and handle potential panic
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Fprintln(printStream, "Tool download failed")
						}
					}()
					t.DownloadTool(printStream, latestVersion)
				}()

				// Execute file
				err := executeFile(printStream)
				if err != nil {
					fmt.Fprintln(printStream, "Update failed due to errors")
				} else {
					fmt.Fprintln(printStream, "Update successfully completed")
				}

				// Clean up regardless of success or failure
				err = OSUtils.DeleteFiles(getToolUnzipLocation())
				if err != nil {
					fmt.Fprintln(printStream, "Error occurred while removing files")
				}

				fmt.Fprintln(printStream)
			}
		}
	}()

	return toolDetails
}

// executeFile executes the installation script
func executeFile(printStream io.Writer) error {
	filePath := filepath.Join(getToolUnzipLocation(), OSUtils.GetInstallScriptFileName())

	cmd := exec.Command(filePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	output := ""
	for scanner.Scan() {
		line := scanner.Text()
		output += line + "\n"
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	fmt.Fprint(printStream, output)
	return nil
}

// getToolUnzipLocation returns the location where the tool is unzipped
func getToolUnzipLocation() string {
	// Implementation needs to match the Java counterpart
	return filepath.Join(OSUtils.GetBallerinaHomePath(), "temp")
}
