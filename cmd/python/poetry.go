package python

import (
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/python/poetry"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type PoetryDoctor struct {
	HTTPClient http.Client
}

func NewPoetryDoctor() *PoetryDoctor {
	client := &http.Client{}
	return &PoetryDoctor{HTTPClient: *client}
}

func (d *PoetryDoctor) Libraries(r parser_io.ReadSeekerAt) []types.Library {
	p := poetry.NewParser()
	libs, _, _ := p.Parse(r)
	return libs
}

func (d *PoetryDoctor) SourceCodeURL(lib types.Library) (string, error) {
	pypi := Pypi{name: lib.Name}
	url, err := pypi.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
