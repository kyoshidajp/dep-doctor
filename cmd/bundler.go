package cmd

import (
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/ruby/bundler"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type BundlerDoctor struct {
	HTTPClient http.Client
}

func NewBundlerDoctor() *BundlerDoctor {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConnsPerHost = -1
	client := &http.Client{Transport: t}
	return &BundlerDoctor{HTTPClient: *client}
}

func (d *BundlerDoctor) Deps(r parser_io.ReadSeekerAt) []types.Library {
	p := &bundler.Parser{}
	deps, _, _ := p.Parse(r)
	return deps
}

func (d *BundlerDoctor) SourceCodeURL(name string) (string, error) {
	rubyGems := RubyGems{name: name}
	url, err := rubyGems.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
