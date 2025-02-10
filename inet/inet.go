package inet

import (
	"fmt"
	"sync"
	"wlpv/gitlab"
	"wlpv/xmlparser"
)

func GetProtocolContents() (map[string][]xmlparser.Protocol, error) {
	var wg sync.WaitGroup

	ch := make(chan gitlab.FetchResult)

	urls := map[string]gitlab.UrlConfig{
		"core": {
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wayland",
			Repository: "wayland",
			Branch:     "main",
			UrlType:    gitlab.UrlTypeFiles,
			Path:       "protocol/wayland.xml",
		},
		"stable": {
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wayland",
			Repository: "wayland-protocols",
			Branch:     "main",
			UrlType:    gitlab.UrlTypeTree,
			Path:       "stable",
		},
		"staging": {
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wayland",
			Repository: "wayland-protocols",
			Branch:     "main",
			UrlType:    gitlab.UrlTypeTree,
			Path:       "staging",
		},
		"unstable": {
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wayland",
			Repository: "wayland-protocols",
			Branch:     "main",
			UrlType:    gitlab.UrlTypeTree,
			Path:       "unstable",
		},
		"wlroots": {
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wlroots",
			Repository: "wlr-protocols",
			Branch:     "master",
			UrlType:    gitlab.UrlTypeTree,
			Path:       "unstable",
		},
		"weston": {
			Origin:     "https://gitlab.freedesktop.org",
			Namespace:  "wayland",
			Repository: "weston",
			Branch:     "main",
			UrlType:    gitlab.UrlTypeTree,
			Path:       "protocol",
		},
		"kde": {
			Origin:     "https://invent.kde.org",
			Namespace:  "libraries",
			Repository: "plasma-wayland-protocols",
			Branch:     "master",
			UrlType:    gitlab.UrlTypeTree,
			Path:       "src/protocols",
		},
	}

	for ns, uc := range urls {
		wg.Add(1)
		go uc.Fetch(&wg, ch, ns)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	protocols := make(map[string][]xmlparser.Protocol)
	for fetchResult := range ch {
		protocols[fetchResult.Namespace] = fetchResult.Protocols
	}

	if len(protocols) == 0 {
		return nil, fmt.Errorf("fetch failed or returned no results")
	}

	return protocols, nil
}
