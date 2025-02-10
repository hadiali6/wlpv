package gitlab

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"wlpv/xmlparser"
)

type FetchResult struct {
	Namespace string
	Protocols []xmlparser.Protocol
}

type fileResponse struct {
	Content string `json:"content"`
}

type treeResponse struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type UrlType uint8

const (
	UrlTypeTree UrlType = iota
	UrlTypeFiles
)

type UrlConfig struct {
	Origin     string
	Namespace  string
	Repository string
	Branch     string
	Path       string
	UrlType
}

func (u UrlConfig) url() string {
	base := fmt.Sprintf("%s/api/v4/projects/%s%%2F%s/repository",
		u.Origin,
		u.Namespace,
		u.Repository,
	)

	formattedPath := url.PathEscape(u.Path)

	switch u.UrlType {
	case UrlTypeFiles:
		return fmt.Sprintf("%s/files/%s?ref=%s", base, formattedPath, u.Branch)
	case UrlTypeTree:
		return fmt.Sprintf(
			"%s/tree?path=%s&ref=%s&per_page=100&recursive=true",
			base,
			formattedPath,
			u.Branch,
		)
	}

	return ""
}

func (u UrlConfig) Fetch(wg *sync.WaitGroup, ch chan<- FetchResult, namespace string) {
	defer wg.Done()

	resp, err := http.Get(u.url())
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	defer resp.Body.Close()

	result := FetchResult{Namespace: namespace}

	switch u.UrlType {
	case UrlTypeFiles:
		protocol, err := u.handleFile(resp)
		if err != nil {
			return
		}

		result.Protocols = append(result.Protocols, protocol)

	case UrlTypeTree:
		filePaths, err := u.handleTree(resp)
		if err != nil {
			return
		}

		var fileWg sync.WaitGroup
		fileCh := make(chan xmlparser.Protocol)

		for _, path := range filePaths {
			fileWg.Add(1)

			fileUrlConfig := u
			fileUrlConfig.UrlType = UrlTypeFiles
			fileUrlConfig.Path = path

			go func(cu UrlConfig, cwg *sync.WaitGroup, cch chan<- xmlparser.Protocol) {
				defer cwg.Done()

				resp, err := http.Get(cu.url())
				if err != nil || resp.StatusCode != http.StatusOK {
					return
				}
				defer resp.Body.Close()

				protocol, err := cu.handleFile(resp)
				if err != nil {
					return
				}

				cch <- protocol
			}(fileUrlConfig, &fileWg, fileCh)
		}

		go func() {
			fileWg.Wait()
			close(fileCh)
		}()

		for protocol := range fileCh {
			result.Protocols = append(result.Protocols, protocol)
		}
	}

	ch <- result
}

func (u UrlConfig) handleFile(resp *http.Response) (xmlparser.Protocol, error) {
	defer resp.Body.Close()

	var fileResp fileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return xmlparser.Protocol{}, err
	}

	decodedContent, err := base64.StdEncoding.DecodeString(fileResp.Content)
	if err != nil {
		return xmlparser.Protocol{}, err
	}

	return xmlparser.ParseProtocol(decodedContent), nil
}

func (u UrlConfig) handleTree(resp *http.Response) ([]string, error) {
	defer resp.Body.Close()

	var nodes []treeResponse
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return nil, err
	}

	var filePaths []string
	for _, node := range nodes {
		if node.Type == "blob" && len(node.Path) > 4 && node.Path[len(node.Path)-4:] == ".xml" {
			filePaths = append(filePaths, node.Path)
		}
	}

	return filePaths, nil
}
