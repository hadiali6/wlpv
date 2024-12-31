package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const help = `usage: wl-protocol-viewer [options] [additional xml protocol file(s)]

    -h -help    Print this help message and exit.
    -v -version Print the version number and exit.
    -no-system  Disable usage of protocols found in /usr/share/*.
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

	noSystemFlag := flag.Bool("no-system", false, "")

	flag.Parse()

	if *shortHelpFlag || *longHelpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *shortVersionFlag || *longVersionFlag {
		fmt.Printf("wl-protocol-viewer version %s\n", version)
		os.Exit(0)
	}

	protocolFiles := flag.Args()

	for _, protocolFile := range protocolFiles {
		fmt.Printf("%s\n", protocolFile)
	}

	var xmlProtocolFiles []string

	const usrshare = "/usr/share/"
	const wl = "wayland"
	const ext = "wayland-protocols"
	const wlr = "wlr-protocols"
	const kde = "plasma-wayland-protocols"
	const weston = "libweston*"

	if *noSystemFlag {
		fmt.Println("no system")
	} else {
		wl_xml, err := findAllFiles(usrshare+wl, ".xml")
		if err == nil {
			xmlProtocolFiles = append(xmlProtocolFiles, wl_xml...)
		}

		ext_xml, err := findAllFiles(usrshare+ext, ".xml")
		if err == nil {
			xmlProtocolFiles = append(xmlProtocolFiles, ext_xml...)
		}

		wlr_xml, err := findAllFiles(usrshare+wlr, ".xml")
		if err == nil {
			xmlProtocolFiles = append(xmlProtocolFiles, wlr_xml...)
		}

		kde_xml, err := findAllFiles(usrshare+kde, ".xml")
		if err == nil {
			xmlProtocolFiles = append(xmlProtocolFiles, kde_xml...)
		}

		westonPathMatches, _ := findMatchingDirs(usrshare + weston)

		weston_xml, err := findAllFiles(westonPathMatches[0], ".xml")
		if err == nil {
			xmlProtocolFiles = append(xmlProtocolFiles, weston_xml...)
		}
	}

	fmt.Println(xmlProtocolFiles)

}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func findAllFiles(path string, extension string) ([]string, error) {
	var xmlFiles []string

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error: accessing path of %s %w\n", path, err)
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), extension) {
			xmlFiles = append(xmlFiles, path)
		}

		return nil
	})

	return xmlFiles, err
}

func findMatchingDirs(patternPath string) ([]string, error) {
	matches, err := filepath.Glob(patternPath)
	if err != nil {
		return nil, fmt.Errorf("error while matching pattern: %w", err)
	}

	var dirs []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return nil, fmt.Errorf("error while checking path %q: %w", match, err)
		}

		if info.IsDir() {
			dirs = append(dirs, match)
		}
	}

	return dirs, nil
}
