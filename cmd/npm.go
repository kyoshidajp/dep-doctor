package cmd

import (
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/nodejs/npm"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type NPMDoctor struct {
	HTTPClient http.Client
}

func NewNPMDoctor() *NPMDoctor {
	return &NPMDoctor{}
}

func (d *NPMDoctor) Deps(r parser_io.ReadSeekerAt) []types.Library {
	p := npm.NewParser()
	deps, _, _ := p.Parse(r)
	return deps
}

func (d *NPMDoctor) SourceCodeURL(lib types.Library) (string, error) {
	nodejs := Nodejs{lib: lib}
	url, err := nodejs.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
