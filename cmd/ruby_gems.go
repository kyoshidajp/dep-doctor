package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// https://guides.rubygems.org/rubygems-org-api/
const RUBY_GEMS_REGISTRY_API = "https://rubygems.org/api/v1/gems/%s.json"

type RubyGemsRegistryResponse struct {
	Name          string `json:"name"`
	SourceCodeUri string `json:"source_code_uri"`
	HomepageUri   string `json:"homepage_uri"`
}

type RubyGems struct {
	name string
}

func (g *RubyGems) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(RUBY_GEMS_REGISTRY_API, g.name)
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

	var Gem RubyGemsRegistryResponse
	err = json.Unmarshal(body, &Gem)
	if err != nil {
		return "", err
	}

	if Gem.SourceCodeUri != "" {
		return Gem.SourceCodeUri, nil
	} else if Gem.HomepageUri != "" {
		return Gem.HomepageUri, nil
	}

	return "", errors.New("source code URL is not found")
}
