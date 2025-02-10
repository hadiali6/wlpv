package main

import (
	"flag"
	"fmt"
	"os"
	"wlpv/cli"
	"wlpv/inet"
	"wlpv/offline"
	"wlpv/tui"
	"wlpv/xmlparser"
)

func main() {
	opts, err := cli.ParseArguments()
	if err != nil {
		os.Exit(1)
	}

	if opts.Help {
		flag.Usage()
		os.Exit(0)
	}

	if opts.Version {
		fmt.Println("0.0.1")
		os.Exit(0)
	}

	protocols := make(map[string][]xmlparser.Protocol)
	protocols["User"] = opts.Additions

	if opts.Offline {
		protocolsFromSystem, err := offline.GetProtocolContents()
		if err != nil {
			os.Exit(1)
		}

		for namespace, protocolGroup := range protocolsFromSystem {
			protocols[namespace] = protocolGroup
		}
	} else {
		protocolsFromNet, err := inet.GetProtocolContents()
		if err != nil {
			os.Exit(1)
		}

		for namespace, protocolGroup := range protocolsFromNet {
			protocols[namespace] = protocolGroup
		}
	}

	if err := tui.Run(opts.Protocol, protocols); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
