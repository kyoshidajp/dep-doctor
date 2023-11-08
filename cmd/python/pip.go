package python

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

func (d *PipDoctor) Parse(r parser_io.ReadSeekerAt) (types.Libraries, error) {
	p := pip.NewParser()
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *PipDoctor) SourceCodeURL(lib types.Library) (string, error) {
	pypi := Pypi{name: lib.Name}
	url, err := pypi.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
