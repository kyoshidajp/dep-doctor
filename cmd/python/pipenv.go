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

func (d *PipenvDoctor) Parse(r parser_io.ReadSeekerAt) (types.Libraries, error) {
	p := pipenv.NewParser()
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *PipenvDoctor) SourceCodeURL(lib types.Library) (string, error) {
	pypi := Pypi{name: lib.Name}
	url, err := pypi.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
