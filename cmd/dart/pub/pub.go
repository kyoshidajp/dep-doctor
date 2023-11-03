package dart

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aquasecurity/go-dep-parser/pkg/dart/pub"
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

const PUB_REGISTRY_API = "https://pub.dev/api/packages/%s"

type PubRegistryResponse struct {
	Latest struct {
		Pubspec struct {
			Repository   string `json:"repository"`
			Homepage     string `json:"homepage"`
			IssueTracker string `json:"issue_tracker"`
		} `json:"pubspec"`
	} `json:"latest"`
}

type Pub struct {
	lib types.Library
}

func (n *Pub) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(PUB_REGISTRY_API, n.lib.Name)
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

	var PubRegistryResponse PubRegistryResponse
	err = json.Unmarshal(body, &PubRegistryResponse)
	if err != nil {
		return "", nil
	}

	if PubRegistryResponse.Latest.Pubspec.Repository != "" {
		return PubRegistryResponse.Latest.Pubspec.Repository, nil
	}

	if PubRegistryResponse.Latest.Pubspec.Homepage != "" {
		return PubRegistryResponse.Latest.Pubspec.Homepage, nil
	}

	return PubRegistryResponse.Latest.Pubspec.IssueTracker, nil
}

type PubDoctor struct {
	HTTPClient http.Client
}

func NewPubDoctor() *PubDoctor {
	return &PubDoctor{}
}

func (d *PubDoctor) Libraries(r parser_io.ReadSeekerAt) []types.Library {
	p := pub.NewParser()
	libs, _, _ := p.Parse(r)
	return libs
}

func (d *PubDoctor) SourceCodeURL(lib types.Library) (string, error) {
	pub := Pub{lib: lib}
	url, err := pub.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
