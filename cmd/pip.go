package cmd

import (
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/python/pip"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type PipDoctor struct {
	HTTPClient http.Client
}

func NewPipDoctor() *PipDoctor {
	client := &http.Client{}
	return &PipDoctor{HTTPClient: *client}
}

func (d *PipDoctor) Deps(r parser_io.ReadSeekerAt) []types.Library {
	p := pip.NewParser()
	deps, _, _ := p.Parse(r)
	return deps
}

func (d *PipDoctor) SourceCodeURL(name string) (string, error) {
	pypi := Pypi{name: name}
	url, err := pypi.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
