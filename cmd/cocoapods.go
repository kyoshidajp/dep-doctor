package cmd

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/swift/cocoapods"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

// It will redirected from https://cdn.cocoapods.org/
// Should be used origin URL?
const COCOA_PODS_REGISTRY_API = "https://cdn.jsdelivr.net/cocoa/Specs/%s"

type CocoaPodsDoctor struct {
	HTTPClient http.Client
}

func NewCococaPodsDoctor() *CocoaPodsDoctor {
	client := &http.Client{}
	return &CocoaPodsDoctor{HTTPClient: *client}
}

func (d *CocoaPodsDoctor) Libraries(r parser_io.ReadSeekerAt) []types.Library {
	p := &cocoapods.Parser{}
	libs, _, _ := p.Parse(r)
	return libs
}

func (d *CocoaPodsDoctor) SourceCodeURL(lib types.Library) (string, error) {
	pod := CocoaPod{
		name:    lib.Name,
		version: lib.Version,
	}
	url, err := pod.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}

type CocoaPodsRegistryResponse struct {
	Homepage string `json:"homepage"`
	Source   struct {
		Git string `json:"git"`
	} `json:"source"`
}

type CocoaPod struct {
	name    string
	version string
}

func (c *CocoaPod) BaseName() string {
	splittedName := strings.Split(c.name, "/")
	return splittedName[0]
}

func (c *CocoaPod) PodspecPath() string {
	md5 := md5.Sum([]byte(c.BaseName()))
	hashString := fmt.Sprintf("%x", md5)
	return fmt.Sprintf("%c/%c/%c/%s/%s/%s.podspec.json", hashString[0], hashString[1], hashString[2], c.BaseName(), c.version, c.BaseName())
}

func (c *CocoaPod) fetchURLFromRegistry(client http.Client) (string, error) {
	url := fmt.Sprintf(COCOA_PODS_REGISTRY_API, c.PodspecPath())
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

	var CocoaPodsRegistryResponse CocoaPodsRegistryResponse
	err = json.Unmarshal(body, &CocoaPodsRegistryResponse)
	if err != nil {
		return "", err
	}

	if CocoaPodsRegistryResponse.Source.Git != "" {
		return CocoaPodsRegistryResponse.Source.Git, nil
	}
	return CocoaPodsRegistryResponse.Homepage, nil
}
