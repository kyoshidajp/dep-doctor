package ruby

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
	client := &http.Client{}
	return &BundlerDoctor{HTTPClient: *client}
}

func (d *BundlerDoctor) Libraries(r parser_io.ReadSeekerAt) []types.Library {
	p := &bundler.Parser{}
	libs, _, _ := p.Parse(r)
	return libs
}

func (d *BundlerDoctor) SourceCodeURL(lib types.Library) (string, error) {
	rubyGems := RubyGems{name: lib.Name}
	url, err := rubyGems.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
