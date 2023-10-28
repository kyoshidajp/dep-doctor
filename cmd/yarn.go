package cmd

import (
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/nodejs/yarn"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type YarnDoctor struct {
}

func NewYarnDoctor() *YarnDoctor {
	return &YarnDoctor{}
}

func (d *YarnDoctor) Deps(r parser_io.ReadSeekerAt) []types.Library {
	p := yarn.NewParser()
	deps, _, _ := p.Parse(r)
	return deps
}

func (d *YarnDoctor) SourceCodeURL(name string) (string, error) {
	nodejs := Nodejs{name: name}
	url, err := nodejs.fetchURLFromRegistry()
	return url, err
}
