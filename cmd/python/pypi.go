package python

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// https://warehouse.pypa.io/api-reference/json.html
const PYPI_REGISTRY_API = "https://pypi.org/pypi/%s/json"

type PypiRegistryResponse struct {
	Info struct {
		ProjectUrls struct {
			SourceCode    string `json:"Source Code"`
			Source        string `json:"Source"`
			Code          string `json:"Code"`
			GitHubProject string `json:"GitHub Project"`
		} `json:"project_urls"`
	} `json:"info"`
}

type Pypi struct {
	name string
}

func (p *Pypi) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(PYPI_REGISTRY_API, p.name)
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
		m := fmt.Sprintf("Got status code: %d from %s", resp.StatusCode, PYPI_REGISTRY_API)
		return "", errors.New(m)
	}

	body, _ := io.ReadAll(resp.Body)

	var PypiRegistryResponse PypiRegistryResponse
	err = json.Unmarshal(body, &PypiRegistryResponse)
	if err != nil {
		return "", nil
	}

	if PypiRegistryResponse.Info.ProjectUrls.SourceCode != "" {
		return PypiRegistryResponse.Info.ProjectUrls.SourceCode, nil
	}

	if PypiRegistryResponse.Info.ProjectUrls.Code != "" {
		return PypiRegistryResponse.Info.ProjectUrls.Code, nil
	}

	if PypiRegistryResponse.Info.ProjectUrls.GitHubProject != "" {
		return PypiRegistryResponse.Info.ProjectUrls.GitHubProject, nil
	}

	return PypiRegistryResponse.Info.ProjectUrls.Source, nil
}
