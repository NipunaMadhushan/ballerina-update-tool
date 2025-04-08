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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// OSUtils provides operating system-related utilities, corresponding to OSUtils.java
var OSUtils = osUtilsStruct{}

type osUtilsStruct struct{}

const (
	BallerinaHomeDir  = ".ballerina"
	BallerinaConfig   = "ballerina-version"
	InstallerVersion  = "installer-version"
	BallerinaListJSON = "local-dists.json"
	UpdateNoticeFile  = "command-notice"
	BirCache          = "bir_cache"
	JarCache          = "jar_cache"
	Repositories      = "repositories"
	CentralCacheRoot  = "central.ballerina.io"
	LocalCacheRoot    = "local"
	Cache             = "cache"
)

// GetInstallationPath returns the installation path
func (o osUtilsStruct) GetInstallationPath() (string, error) {
	// In Go, we'd typically use a known directory like /usr/local/ballerina on Unix systems
	// or C:\Program Files\Ballerina on Windows, or possibly a path based on the executable
	// For this conversion, we'll use a placeholder approach
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Go two directories up from the executable to get the installation root
	installPath := filepath.Dir(filepath.Dir(execPath))
	return installPath, nil
}

// GetExecutableFileName returns the name of the Ballerina executable for the current OS
func (o osUtilsStruct) GetExecutableFileName(distribution string) string {
	fileName := "bal"
	if o.IsWindows() {
		fileName = "bal.bat"
	}

	// First check for the new-style executable
	distPath := ToolUtil.GetDistributionsPath()
	filePath := filepath.Join(distPath, ToolUtil.GetType(distribution)+"-"+distribution, "bin", fileName)

	if fileExists(filePath) {
		return fileName
	}

	// Fall back to old-style executable
	if o.IsWindows() {
		return "ballerina.bat"
	}
	return "ballerina"
}

// GetInstallScriptFileName returns the name of the install script for the current OS
func (o osUtilsStruct) GetInstallScriptFileName() string {
	if o.IsWindows() {
		return "install.bat"
	}
	return "install"
}

// GetDebugAdapterName returns the name of the debug adapter script for the current OS
func (o osUtilsStruct) GetDebugAdapterName() string {
	if o.IsWindows() {
		return "debug-adapter-launcher.bat"
	}
	return "debug-adapter-launcher.sh"
}

// GetLangServerLauncherName returns the name of the language server launcher for the current OS
func (o osUtilsStruct) GetLangServerLauncherName() string {
	if o.IsWindows() {
		return "language-server-launcher.bat"
	}
	return "language-server-launcher.sh"
}

// GetBallerinaHomePath returns the Ballerina home directory path
func (o osUtilsStruct) GetBallerinaHomePath() string {
	return filepath.Join(o.GetUserHome(), BallerinaHomeDir)
}

// GetBallerinaVersionFilePath returns the path to the Ballerina version file
func (o osUtilsStruct) GetBallerinaVersionFilePath() string {
	ballerinaHome := filepath.Join(o.GetUserHome(), BallerinaHomeDir)
	ballerinaVersionFile := filepath.Join(ballerinaHome, BallerinaConfig)

	if !fileExists(ballerinaVersionFile) {
		// Create the file and parent directories if they don't exist
		os.MkdirAll(ballerinaHome, 0755)
		emptyFile, err := os.Create(ballerinaVersionFile)
		if err == nil {
			emptyFile.Close()
			// Set file permissions
			os.Chmod(ballerinaHome, 0755)
			os.Chmod(ballerinaVersionFile, 0644)

			// Set the default version
			ToolUtil.setVersion(ballerinaVersionFile, ToolUtil.GetCurrentInstalledBallerinaVersion())
		}
	}

	return ballerinaVersionFile
}

// GetInstallerVersionFilePath returns the path to the installer version file
func (o osUtilsStruct) GetInstallerVersionFilePath() string {
	ballerinaHome := filepath.Join(o.GetUserHome(), BallerinaHomeDir)
	installerVersionFile := filepath.Join(ballerinaHome, InstallerVersion)

	installedVersionPath := o.GetInstalledInstallerVersionPath()

	// If the installed version doesn't exist but the user's version does, delete the user's version
	if !fileExists(installedVersionPath) && fileExists(installerVersionFile) {
		os.Remove(installerVersionFile)
	}

	// If the installed version exists but the user's version doesn't, create the user's version
	if fileExists(installedVersionPath) && !fileExists(installerVersionFile) {
		os.MkdirAll(ballerinaHome, 0755)
		emptyFile, err := os.Create(installerVersionFile)
		if err == nil {
			emptyFile.Close()
			// Set file permissions
			os.Chmod(ballerinaHome, 0755)
			os.Chmod(installerVersionFile, 0644)

			// Set the installer version
			ToolUtil.setInstallerVersion(installerVersionFile)
			ToolUtil.setVersion(o.GetBallerinaVersionFilePath(), ToolUtil.GetCurrentInstalledBallerinaVersion())
		}
	}

	return installerVersionFile
}

// GetBallerinaDistListFilePath returns the path to the Ballerina distribution list file
func (o osUtilsStruct) GetBallerinaDistListFilePath() string {
	return filepath.Join(o.GetUserHome(), BallerinaHomeDir, BallerinaListJSON)
}

// GetInstalledConfigPath returns the path to the installed configuration file
func (o osUtilsStruct) GetInstalledConfigPath() string {
	return filepath.Join(ToolUtil.GetDistributionsPath(), BallerinaConfig)
}

// GetInstalledInstallerVersionPath returns the path to the installed installer version file
func (o osUtilsStruct) GetInstalledInstallerVersionPath() string {
	return filepath.Join(ToolUtil.GetDistributionsPath(), InstallerVersion)
}

// GetUpdateNoticePath returns the path to the update notice file
func (o osUtilsStruct) GetUpdateNoticePath() string {
	return filepath.Join(o.GetUserHome(), BallerinaHomeDir, UpdateNoticeFile)
}

// UpdateNotice checks if the update notice should be shown
func (o osUtilsStruct) UpdateNotice() bool {
	today := time.Now().Format("2006-01-02") // ISO format date
	noticePath := o.GetUpdateNoticePath()

	// If the file doesn't exist, create it and return true
	if !fileExists(noticePath) {
		// Create parent dirs if needed
		os.MkdirAll(filepath.Dir(noticePath), 0755)

		// Write today's date to the file
		err := ioutil.WriteFile(noticePath, []byte(today+"\n"), 0644)
		if err != nil {
			return false
		}
		return true
	}

	// Read the last updated date from the file
	data, err := ioutil.ReadFile(noticePath)
	if err != nil {
		return false
	}

	lastUpdatedDate := strings.TrimSpace(string(data))
	lastDate, err := time.Parse("2006-01-02", lastUpdatedDate)
	if err != nil {
		return false
	}

	// Check if more than 1 day has passed
	daysSinceLastUpdate := int(time.Since(lastDate).Hours() / 24)
	showNotice := daysSinceLastUpdate > 1

	// If we're showing the notice, update the file
	if showNotice {
		err := ioutil.WriteFile(noticePath, []byte(today+"\n"), 0644)
		if err != nil {
			return false
		}
	}

	return showNotice
}

// ClearBirCacheLocation clears the BIR cache
func (o osUtilsStruct) ClearBirCacheLocation(outStream io.Writer) error {
	birCachePath := filepath.Join(o.GetUserHome(), BallerinaHomeDir, BirCache)
	return o.DeleteDirectory(birCachePath, outStream)
}

// ClearJarCacheLocation clears the JAR cache
func (o osUtilsStruct) ClearJarCacheLocation(outStream io.Writer) error {
	jarCachePath := filepath.Join(o.GetUserHome(), BallerinaHomeDir, JarCache)
	return o.DeleteDirectory(jarCachePath, outStream)
}

// DeleteDirectory removes a directory and all its contents
func (o osUtilsStruct) DeleteDirectory(path string, outStream io.Writer) error {
	if !dirExists(path) {
		return nil
	}

	// First check for write permissions
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if we can write to this file/directory
		if !isWritable(filePath) {
			return fmt.Errorf("permission denied: you do not have write access to '%s'", filePath)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Now actually delete the files
	return os.RemoveAll(path)
}

// GetUserAgent returns the user agent string for HTTP requests
func (o osUtilsStruct) GetUserAgent(ballerinaVersion, toolVersion, distributionType string) string {
	os := "none"
	if o.IsWindows() {
		os = "win-64"
	} else if o.IsUnix() || o.IsSolaris() {
		os = "linux-64"
	} else if o.IsMac() {
		if o.IsArmArchitecture() {
			os = "macos-arm-64"
		} else {
			os = "macos-64"
		}
	}

	return fmt.Sprintf("%s/%s (%s) Updater/%s", distributionType, ballerinaVersion, os, toolVersion)
}

// IsWindows checks if the OS is Windows
func (o osUtilsStruct) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsMac checks if the OS is macOS
func (o osUtilsStruct) IsMac() bool {
	return runtime.GOOS == "darwin"
}

// IsUnix checks if the OS is a Unix-like system (excluding macOS)
func (o osUtilsStruct) IsUnix() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" || runtime.GOOS == "netbsd"
}

// IsSolaris checks if the OS is Solaris
func (o osUtilsStruct) IsSolaris() bool {
	return runtime.GOOS == "solaris"
}

// IsArmArchitecture checks if the architecture is ARM (for macOS)
func (o osUtilsStruct) IsArmArchitecture() bool {
	macArchitecture := os.Getenv("BALLERINA_MAC_ARCHITECTURE")
	return macArchitecture == "arm64"
}

// GetUserHome returns the user's home directory
func (o osUtilsStruct) GetUserHome() string {
	homeDir := os.Getenv("HOME")

	// Handle special case for running as root on Unix
	if o.IsUnix() && strings.Contains(homeDir, "root") {
		sudoUser := os.Getenv("SUDO_USER")
		if sudoUser != "" {
			return "/home/" + sudoUser
		}
	}

	// If HOME not set, use user.Current()
	if homeDir == "" {
		usr, err := user.Current()
		if err == nil {
			homeDir = usr.HomeDir
		} else {
			// Fallback to system property
			homeDir = os.Getenv("USERPROFILE") // Windows
			if homeDir == "" {
				homeDir = os.Getenv("HOME") // Unix
			}
		}
	}

	return homeDir
}

// DeleteFiles removes a directory and all its contents
func (o osUtilsStruct) DeleteFiles(dirPath string) error {
	if dirPath == "" || !dirExists(dirPath) {
		return nil
	}

	// Check permissions first
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !isWritable(path) {
			return fmt.Errorf("permission denied: you do not have write access to '%s'", dirPath)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Now remove everything
	return os.RemoveAll(dirPath)
}

// DeleteCaches deletes the repository caches in the user's home directory
func (o osUtilsStruct) DeleteCaches(version string, outStream io.Writer) error {
	centralCache := filepath.Join(o.GetUserHome(), BallerinaHomeDir, Repositories, CentralCacheRoot, Cache+"-"+version)
	localCache := filepath.Join(o.GetUserHome(), BallerinaHomeDir, Repositories, LocalCacheRoot, Cache+"-"+version)

	err := o.DeleteDirectory(centralCache, outStream)
	if err != nil {
		return err
	}

	return o.DeleteDirectory(localCache, outStream)
}
