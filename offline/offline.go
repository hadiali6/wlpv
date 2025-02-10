package offline

import (
	"wlpv/util"
	"wlpv/xmlparser"
)

const usrshare = "/usr/share/"
const wl = "wayland"
const stable = "wayland-protocols/stable"
const staging = "wayland-protocols/staging"
const unstable = "wayland-protocols/unstable"
const wlr = "wlr-protocols"
const kde = "plasma-wayland-protocols"
const weston = "libweston*"

func getProtocolsInDir(path string) ([]xmlparser.Protocol, error) {
	files, err := util.AllFilesInDir(path, ".xml")
	if err != nil {
		return nil, err
	}

	contents, err := util.ReadAllFiles(files)
	if err != nil {
		return nil, err
	}

	var protocols []xmlparser.Protocol
	for _, content := range contents {
		protocols = append(protocols, xmlparser.ParseProtocol(content))
	}

	return protocols, nil
}

func GetProtocolContents() (map[string][]xmlparser.Protocol, error) {
	waylandProtocols, err := getProtocolsInDir(usrshare + wl) // length should only be 1 (only wayland.xml)
	if err != nil {
		return nil, err
	}

	stableProtocols, err := getProtocolsInDir(usrshare + stable)
	if err != nil {
		return nil, err
	}

	stagingProtocols, err := getProtocolsInDir(usrshare + staging)
	if err != nil {
		return nil, err
	}

	unstableProtocols, err := getProtocolsInDir(usrshare + unstable)
	if err != nil {
		return nil, err
	}

	wlrProtocols, err := getProtocolsInDir(usrshare + wlr)
	if err != nil {
		return nil, err
	}

	kdeProtocols, err := getProtocolsInDir(usrshare + kde)
	if err != nil {
		return nil, err
	}

	westonPathMatches, err := util.FindMatchingDirs(usrshare + weston)
	if err != nil {
		return nil, err
	}

	westonProtocols, err := getProtocolsInDir(westonPathMatches[0])
	if err != nil {
		return nil, err
	}

	var result = map[string][]xmlparser.Protocol{
		"core":     waylandProtocols,
		"stable":   stableProtocols,
		"staging":  stagingProtocols,
		"unstable": unstableProtocols,
		"wlroots":  wlrProtocols,
		"weston":   westonProtocols,
		"kde":      kdeProtocols,
	}

	return result, nil
}
