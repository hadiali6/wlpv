package gitlab

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type FileResponse struct {
	Content string `json:"content"`
}

type TreeResponse struct {
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

func (uc UrlConfig) ToUrl() (urlstr string, err error) {
	base := fmt.Sprintf("%s/api/v4/projects/%s%%2F%s/repository",
		uc.Origin,
		uc.Namespace,
		uc.Repository,
	)

	var result string

	formattedPath := url.PathEscape(uc.Path)
	switch uc.UrlType {
	case UrlTypeFiles:
		result = fmt.Sprintf("%s/files/%s?ref=%s", base, formattedPath, uc.Branch)
	case UrlTypeTree:
		result = fmt.Sprintf(
			"%s/tree?path=%s&ref=%s&per_page=100&recursive=true",
			base,
			formattedPath,
			uc.Branch,
		)
	default:
		return "", fmt.Errorf("invalid UrlType")
	}

	return result, nil
}

func (uc UrlConfig) Get() (*http.Response, error) {
	u, err := uc.ToUrl()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file: %s", resp.Status)
	}

	return resp, nil
}

func (uc UrlConfig) FetchFiles(resp *http.Response) (string, error) {
	content, err := handleFileResponse(resp)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (uc UrlConfig) FetchTree(resp *http.Response) ([]string, error) {
	files, err := handleTreeResponse(resp)
	if err != nil {
		return nil, err
	}

	var contents []string
	for _, file := range files {
		urlConfig := uc
		urlConfig.UrlType = UrlTypeFiles
		urlConfig.Path = file

		resp, err := urlConfig.Get()
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		content, err := urlConfig.FetchFiles(resp)
		if err != nil {
			return nil, err
		}

		contents = append(contents, content)
	}

	return contents, nil
}

func handleFileResponse(resp *http.Response) (string, error) {
	var fileResp FileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(fileResp.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %v", err)
	}

	return string(decodedContent), nil
}

func handleTreeResponse(resp *http.Response) ([]string, error) {
	var nodes []TreeResponse
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return nil, err
	}

	var files []string
	for _, node := range nodes {
		if node.Type == "blob" && len(node.Path) > 4 && node.Path[len(node.Path)-4:] == ".xml" {
			files = append(files, node.Path)
		}
	}

	return files, nil
}
