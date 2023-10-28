package cmd

import (
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/nodejs/npm"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type NPMDoctor struct {
}

func NewNPMDoctor() *NPMDoctor {
	return &NPMDoctor{}
}

func (d *NPMDoctor) Deps(r parser_io.ReadSeekerAt) []types.Library {
	p := npm.NewParser()
	deps, _, _ := p.Parse(r)
	return deps
}

func (d *NPMDoctor) SourceCodeURL(name string) (string, error) {
	nodejs := Nodejs{name: name}
	url, err := nodejs.fetchURLFromRegistry()
	return url, err
}
