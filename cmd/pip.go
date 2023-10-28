package cmd

import (
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/python/pip"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type PipDoctor struct {
}

func NewPipDoctor() *PipDoctor {
	return &PipDoctor{}
}

func (d *PipDoctor) Deps(r parser_io.ReadSeekerAt) []types.Library {
	p := pip.NewParser()
	deps, _, _ := p.Parse(r)
	return deps
}

func (d *PipDoctor) SourceCodeURL(name string) (string, error) {
	pypi := Pypi{name: name}
	url, err := pypi.fetchURLFromRepository()
	return url, err
}
