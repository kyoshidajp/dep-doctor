package php

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

// https://packagist.org/apidoc#get-package-data
const PACKAGIST_REGISTRY_API = "https://repo.packagist.org/p2/%s.json"

type PackagistRegistryResponse struct {
	Packages map[string][]struct {
		Source struct {
			URL string `json:"url"`
		} `json:"source"`
	} `json:"packages"`
}

type Packagist struct {
	lib types.Library
}

func (p *Packagist) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(PACKAGIST_REGISTRY_API, p.lib.Name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		m := fmt.Sprintf("Got status code: %d from %s", resp.StatusCode, url)
		return "", errors.New(m)
	}

	body, _ := io.ReadAll(resp.Body)

	var registryResponse PackagistRegistryResponse
	err = json.Unmarshal(body, &registryResponse)
	if err != nil {
		return "", nil
	}

	return registryResponse.Packages[p.lib.Name][0].Source.URL, nil
}
