package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// https://packagist.org/apidoc#get-package-data
const PACKAGIST_REGISTRY_API = "https://repo.packagist.org/p2/%s.json"

/*
	type PackagistRegistryResponse struct {
		Packages []struct {
			Source struct {
				URL string `json:"url"`
			} `json:"source"`
		} `json:"packages"`
	}
*/
type PackagistRegistryResponse struct {
	Packages map[string][]struct {
		Source struct {
			URL string `json:"url"`
		} `json:"source"`
	} `json:"packages"`
}

type Packagist struct {
	name string
}

func (p *Packagist) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(PACKAGIST_REGISTRY_API, p.name)
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
		m := fmt.Sprintf("Got status code: %d from %s", resp.StatusCode, RUBY_GEMS_REGISTRY_API)
		return "", errors.New(m)
	}

	body, _ := io.ReadAll(resp.Body)

	var registryResponse PackagistRegistryResponse
	err = json.Unmarshal(body, &registryResponse)
	if err != nil {
		return "", nil
	}

	return registryResponse.Packages[p.name][0].Source.URL, nil
}
