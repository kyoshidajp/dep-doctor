package golang

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

// https://proxy.golang.org/
const PROXY_GOLANG_REGISTRY_API = "https://proxy.golang.org/%s/@latest"

type ProxyGolangRegistryResponse struct {
	Origin struct {
		Vcs string `json:"VCS"`
		URL string `json:"URL"`
	} `json:"Origin"`
}

type ProxyGolang struct {
	lib types.Library
}

func (g *ProxyGolang) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(PROXY_GOLANG_REGISTRY_API, g.lib.Name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		m := fmt.Sprintf("Got status code: %d from %s", resp.StatusCode, PROXY_GOLANG_REGISTRY_API)
		return "", errors.New(m)
	}

	body, _ := io.ReadAll(resp.Body)

	var ProxyGolangRegistryResponse ProxyGolangRegistryResponse
	err = json.Unmarshal(body, &ProxyGolangRegistryResponse)
	if err != nil {
		return "", err
	}

	return ProxyGolangRegistryResponse.Origin.URL, nil
}
