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
}

func (g *RubyGems) fetchURLFromRegistry(name string) (string, error) {
	url := fmt.Sprintf(RUBY_GEMS_REGISTRY_API, name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	client := new(http.Client)
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var Gem RubyGemsRegistryResponse
	err := json.Unmarshal(body, &Gem)
	if err != nil {
		return "", errors.New("error: Unknown response")
	}

	if Gem.SourceCodeUri != "" {
		return Gem.SourceCodeUri, nil
	} else if Gem.HomepageUri != "" {
		return Gem.HomepageUri, nil
	}

	return "", nil
}
