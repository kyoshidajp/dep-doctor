package erlang_elixir

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const HEX_REGISTRY_API = "https://hex.pm/api/packages/%s"

type HexRegistryResponse struct {
	Meta struct {
		Links struct {
			GitHub string `json:"GitHub"`
		} `json:"links"`
	} `json:"meta"`
}

type Hex struct {
	name string
}

func (g *Hex) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(HEX_REGISTRY_API, g.name)
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

	var hex HexRegistryResponse
	err = json.Unmarshal(body, &hex)
	if err != nil {
		return "", err
	}

	return hex.Meta.Links.GitHub, nil
}
