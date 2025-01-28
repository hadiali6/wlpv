package offline

import "github.com/hadiali6/wl-protocol-viewer/util"

const usrshare = "/usr/share/"
const wl = "wayland"
const ext = "wayland-protocols"
const wlr = "wlr-protocols"
const kde = "plasma-wayland-protocols"
const weston = "libweston*"

func FindAllProtocols() []string {
	var xmlProtocolFiles []string

	wl_xml, err := util.AllFilesInDir(usrshare+wl, ".xml")
	if err == nil {
		xmlProtocolFiles = append(xmlProtocolFiles, wl_xml...)
	}

	ext_xml, err := util.AllFilesInDir(usrshare+ext, ".xml")
	if err == nil {
		xmlProtocolFiles = append(xmlProtocolFiles, ext_xml...)
	}

	wlr_xml, err := util.AllFilesInDir(usrshare+wlr, ".xml")
	if err == nil {
		xmlProtocolFiles = append(xmlProtocolFiles, wlr_xml...)
	}

	kde_xml, err := util.AllFilesInDir(usrshare+kde, ".xml")
	if err == nil {
		xmlProtocolFiles = append(xmlProtocolFiles, kde_xml...)
	}

	westonPathMatches, _ := util.FindMatchingDirs(usrshare + weston)

	weston_xml, err := util.AllFilesInDir(westonPathMatches[0], ".xml")
	if err == nil {
		xmlProtocolFiles = append(xmlProtocolFiles, weston_xml...)
	}

	return xmlProtocolFiles
}
