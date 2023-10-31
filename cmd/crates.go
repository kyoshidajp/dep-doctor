package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const CRATES_REGISTRY_API = "https://crates.io/api/v1/crates/%s"

type CratesRegistryResponse struct {
	Crate struct {
		Repository string `json:"repository"`
	} `json:"crate"`
}

type Crate struct {
	name string
}

func (c *Crate) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(CRATES_REGISTRY_API, c.name)
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
		m := fmt.Sprintf("Got status code: %d from %s", resp.StatusCode, url)
		return "", errors.New(m)
	}

	body, _ := io.ReadAll(resp.Body)

	var CratesRegistryResponse CratesRegistryResponse
	err = json.Unmarshal(body, &CratesRegistryResponse)
	if err != nil {
		return "", err
	}

	return CratesRegistryResponse.Crate.Repository, nil
}
