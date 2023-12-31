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

	var PubRegistryResponse PubRegistryResponse
	err = json.Unmarshal(body, &PubRegistryResponse)
	if err != nil {
		return "", err
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

func (d *PubDoctor) Parse(r parser_io.ReadSeekerAt) (types.Libraries, error) {
	p := pub.NewParser()
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *PubDoctor) SourceCodeURL(lib types.Library) (string, error) {
	pub := Pub{lib: lib}
	url, err := pub.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
