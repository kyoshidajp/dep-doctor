package python

import (
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/python/pipenv"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type PipenvDoctor struct {
	HTTPClient http.Client
}

func NewPipenvDoctor() *PipenvDoctor {
	client := &http.Client{}
	return &PipenvDoctor{HTTPClient: *client}
}

func (d *PipenvDoctor) Libraries(r parser_io.ReadSeekerAt) []types.Library {
	p := pipenv.NewParser()
	libs, _, _ := p.Parse(r)
	return libs
}

func (d *PipenvDoctor) SourceCodeURL(lib types.Library) (string, error) {
	pypi := Pypi{name: lib.Name}
	url, err := pypi.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
