package nodejs

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

func (d *NPMDoctor) Parse(r parser_io.ReadSeekerAt) (types.Libraries, error) {
	p := npm.NewParser()
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *NPMDoctor) SourceCodeURL(lib types.Library) (string, error) {
	nodejs := Nodejs{lib: lib}
	url, err := nodejs.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
