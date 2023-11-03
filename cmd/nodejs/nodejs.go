package nodejs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

// https://docs.npmjs.com/cli/v8/using-npm/registry
const NODEJS_REGISTRY_API = "https://registry.npmjs.org/%s"

type NodejsRegistryResponse struct {
	Repository struct {
		Url string `json:"url"`
	}
}

type Nodejs struct {
	lib types.Library
}

func (n *Nodejs) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(NODEJS_REGISTRY_API, n.lib.Name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		m := fmt.Sprintf("Got status code: %d from %s", resp.StatusCode, url)
		return "", errors.New(m)
	}

	body, _ := io.ReadAll(resp.Body)

	var NodejsRegistryResponse NodejsRegistryResponse
	err = json.Unmarshal(body, &NodejsRegistryResponse)
	if err != nil {
		return "", nil
	}

	return NodejsRegistryResponse.Repository.Url, nil
}
