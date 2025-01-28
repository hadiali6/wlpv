package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hadiali6/wl-protocol-viewer/gitlab"
	"github.com/hadiali6/wl-protocol-viewer/offline"
	"github.com/hadiali6/wl-protocol-viewer/util"
)

const help = `usage: wl-protocol-viewer [options] [additional xml protocol file(s)]

    -h -help    Print this help message and exit.
    -v -version Print the version number and exit.
    -offline    Search for protocols found in /usr/share/* instead of fetching from git
`

const version = "0.0.1"

func main() {
	flag.Usage = func() {
		fmt.Print(help)
	}

	shortHelpFlag := flag.Bool("h", false, "")
	longHelpFlag := flag.Bool("help", false, "")

	shortVersionFlag := flag.Bool("v", false, "")
	longVersionFlag := flag.Bool("version", false, "")

	offlineFlag := flag.Bool("offline", false, "")

	flag.Parse()

	if *shortHelpFlag || *longHelpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *shortVersionFlag || *longVersionFlag {
		fmt.Printf("wl-protocol-viewer version %s\n", version)
		os.Exit(0)
	}

	allProtocolFiles := GetProtocolFilesToRead(flag.Args(), *offlineFlag)
	xmlFileContents, err := util.ReadAllFiles(allProtocolFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	xmlNetContents, err := GetProtocolsFromInternet()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting protocol from the internet: %s\n", err.Error())
		os.Exit(1)
	}

	xmlFileContents = append(xmlFileContents, xmlNetContents...)

	for _, file := range xmlFileContents {
		println(file)
	}
}

func GetProtocolsFromInternet() ([]string, error) {
	var contents []string

	{
		var wayland = gitlab.UrlConfig{
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wayland",
			Repository: "wayland",
			Branch:     "main",
			UrlType:    gitlab.UrlTypeFiles,
			Path:       "protocol/wayland.xml",
		}

		waylandResp, err := wayland.Get()
		if err != nil {
			return nil, err
		}
		defer waylandResp.Body.Close()

		content, err := wayland.FetchFiles(waylandResp)
		if err != nil {
			return nil, err
		}

		contents = append(contents, content)
	}

	{
		var extStable = gitlab.UrlConfig{
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wayland",
			Repository: "wayland-protocols",
			Branch:     "main",
			UrlType:    gitlab.UrlTypeTree,
			Path:       "stable",
		}

		extStableResp, err := extStable.Get()
		if err != nil {
			return nil, err
		}
		defer extStableResp.Body.Close()

		content, err := extStable.FetchTree(extStableResp)
		if err != nil {
			return nil, err
		}

		contents = append(contents, content...)
	}

	{
		var extStaging = gitlab.UrlConfig{
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wayland",
			Repository: "wayland-protocols",
			Branch:     "main",
			UrlType:    gitlab.UrlTypeTree,
			Path:       "staging",
		}

		extStagingResp, err := extStaging.Get()
		if err != nil {
			return nil, err
		}
		defer extStagingResp.Body.Close()

		content, err := extStaging.FetchTree(extStagingResp)
		if err != nil {
			return nil, err
		}

		contents = append(contents, content...)
	}

	{
		var extUnstable = gitlab.UrlConfig{
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wayland",
			Repository: "wayland-protocols",
			Branch:     "main",
			UrlType:    gitlab.UrlTypeTree,
			Path:       "unstable",
		}

		extUnstableResp, err := extUnstable.Get()
		if err != nil {
			return nil, err
		}
		defer extUnstableResp.Body.Close()

		content, err := extUnstable.FetchTree(extUnstableResp)
		if err != nil {
			return nil, err
		}

		contents = append(contents, content...)
	}

	return contents, nil
}

func GetProtocolFilesToRead(filePathsFromArgs []string, getProtocolsFromSystem bool) []string {
	var allFiles []string

	for _, path := range filePathsFromArgs {
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error matching path '%s'. error: '%s'\n", path, err.Error())
			os.Exit(1)
		}

		for _, match := range matches {
			fileInfo, err := os.Stat(match)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "Path '%s' doesn't exist\n", path)
					os.Exit(1)
				} else {
					fmt.Fprintf(
						os.Stderr,
						"Error checking path '%s'. error: '%s'\n",
						path,
						err.Error(),
					)
					os.Exit(1)
				}
			}

			if fileInfo.IsDir() {
				files, err := util.AllFilesInDir(match, ".xml")
				if err != nil {
					fmt.Fprintf(
						os.Stderr,
						"Error finding all '.xml' files in directory: '%s'. error: '%s'\n",
						match,
						err.Error(),
					)
					os.Exit(1)
				}

				allFiles = append(allFiles, files...)
			} else {
				if filepath.Ext(match) == ".xml" {
					allFiles = append(allFiles, match)
				}
			}
		}
	}

	if getProtocolsFromSystem {
		var protocolFiles = offline.FindAllProtocols()

		for _, filePathFromArgs := range allFiles {
			for index, filePathFromSystem := range protocolFiles {
				fileFromArgs := filepath.Base(filePathFromArgs)
				fileFromSystem := filepath.Base(filePathFromSystem)

				if fileFromArgs == fileFromSystem {
					fmt.Fprintf(
						os.Stderr,
						"Duplicate filenames! '%s' '%s'\n",
						filePathFromArgs,
						filePathFromSystem,
					)
					protocolFiles = util.RemoveAtIndex(protocolFiles, index)
				}
			}
		}

		allFiles = append(allFiles, protocolFiles...)
	}

	return allFiles
}
