package cmd

import (
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/ruby/bundler"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type BundlerDoctor struct {
}

func NewBundlerDoctor() *BundlerDoctor {
	return &BundlerDoctor{}
}

func (d *BundlerDoctor) Deps(r parser_io.ReadSeekerAt) []types.Library {
	p := &bundler.Parser{}
	deps, _, _ := p.Parse(r)
	return deps
}

func (d *BundlerDoctor) SourceCodeURL(name string) (string, error) {
	rubyGems := RubyGems{name: name}
	url, err := rubyGems.fetchURLFromRegistry()
	return url, err
}
