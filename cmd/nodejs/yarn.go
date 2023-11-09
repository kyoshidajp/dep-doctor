package nodejs

import (
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/nodejs/yarn"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type YarnDoctor struct {
	HTTPClient http.Client
}

func NewYarnDoctor() *YarnDoctor {
	client := &http.Client{}
	return &YarnDoctor{HTTPClient: *client}
}

func (d *YarnDoctor) Parse(r parser_io.ReadSeekerAt) (types.Libraries, error) {
	p := yarn.NewParser()
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *YarnDoctor) SourceCodeURL(lib types.Library) (string, error) {
	nodejs := Nodejs{lib: lib}
	url, err := nodejs.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
