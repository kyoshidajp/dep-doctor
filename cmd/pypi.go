package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://warehouse.pypa.io/api-reference/json.html
const PYPI_REGISTRY_API = "https://pypi.org/pypi/%s/json"

type PypiRegistryResponse struct {
	Info struct {
		ProjectUrls struct {
			SourceCode string `json:"Source Code"`
			Source     string `json:"Source"`
		} `json:"project_urls"`
	} `json:"info"`
}

type Pypi struct {
}

func (p *Pypi) fetchURLFromRepository(name string) (string, error) {
	url := fmt.Sprintf(PYPI_REGISTRY_API, name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	client := new(http.Client)
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var PypiRegistryResponse PypiRegistryResponse
	err := json.Unmarshal(body, &PypiRegistryResponse)
	if err != nil {
		return "", nil
	}

	if PypiRegistryResponse.Info.ProjectUrls.SourceCode != "" {
		return PypiRegistryResponse.Info.ProjectUrls.SourceCode, nil
	}

	return PypiRegistryResponse.Info.ProjectUrls.Source, nil
}