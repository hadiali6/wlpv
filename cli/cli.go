package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"wlpv/util"
	"wlpv/xmlparser"
)

const help = `usage: wlpv [options] [protocol name]

    -h -help       Print this help message and exit.
    -v -version    Print the version number and exit.
    -a -add <path> Additional xml protocol file.
    -offline       Search for protocols found in /usr/share/* instead of fetching from git.
`

type Options struct {
	Help      bool                 // print help and exit(0)
	Version   bool                 // print version and exit(0)
	Offline   bool                 // offline mode
	Additions []xmlparser.Protocol // additional protocols
	Protocol  string               // name of protocol to directly open
}

type paths []string

func (p *paths) String() string {
	return strings.Join(*p, ", ")
}

func (p *paths) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func ParseArguments() (Options, error) {
	flag.Usage = func() {
		fmt.Print(help)
	}

	shortHelpFlag := flag.Bool("h", false, "")
	longHelpFlag := flag.Bool("help", false, "")

	shortVersionFlag := flag.Bool("v", false, "")
	longVersionFlag := flag.Bool("version", false, "")

	var paths paths
	flag.Var(&paths, "a", "")
	flag.Var(&paths, "add", "")

	offlineFlag := flag.Bool("offline", false, "")

	flag.Parse()

	var opts Options

	opts.Help = *shortHelpFlag || *longHelpFlag
	opts.Version = *shortVersionFlag || *longVersionFlag
	opts.Offline = *offlineFlag

	filePaths, err := getFilePaths(paths)
	if err != nil {
		return opts, err
	}

	contents, err := util.ReadAllFiles(filePaths)
	if err != nil {
		return opts, err
	}

	for _, content := range contents {
		opts.Additions = append(opts.Additions, xmlparser.ParseProtocol(content))
	}

	opts.Protocol = flag.Arg(0)

	return opts, nil
}

func getFilePaths(paths []string) ([]string, error) {
	var filePaths []string

	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				return nil, err
			}

			if info.IsDir() {
				files, err := util.AllFilesInDir(match, ".xml")
				if err != nil {
					return nil, err
				}

				filePaths = append(filePaths, files...)
			} else {
				if filepath.Ext(match) == ".xml" {
					filePaths = append(filePaths, match)
				}
			}
		}
	}

	return filePaths, nil
}
