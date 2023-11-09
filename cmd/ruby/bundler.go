package ruby

import (
	"net/http"

	"github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/ruby/bundler"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type BundlerDoctor struct {
	HTTPClient http.Client
}

func NewBundlerDoctor() *BundlerDoctor {
	client := &http.Client{}
	return &BundlerDoctor{HTTPClient: *client}
}

func (d *BundlerDoctor) Parse(r io.ReadSeekerAt) (types.Libraries, error) {
	p := &bundler.Parser{}
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *BundlerDoctor) SourceCodeURL(lib types.Library) (string, error) {
	rubyGems := RubyGems{name: lib.Name}
	url, err := rubyGems.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
