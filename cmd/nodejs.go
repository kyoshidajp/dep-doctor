package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://docs.npmjs.com/cli/v8/using-npm/registry
const NODEJS_REGISTRY_API = "https://registry.npmjs.org/%s"

type NodejsRegistryResponse struct {
	Repository struct {
		Url string `json:"url"`
	}
}

type Nodejs struct {
	name string
}

func (n *Nodejs) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(NODEJS_REGISTRY_API, n.name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var NodejsRegistryResponse NodejsRegistryResponse
	err := json.Unmarshal(body, &NodejsRegistryResponse)
	if err != nil {
		return "", nil
	}

	return NodejsRegistryResponse.Repository.Url, nil
}
