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

func (d *YarnDoctor) Libraries(r parser_io.ReadSeekerAt) []types.Library {
	p := yarn.NewParser()
	libs, _, _ := p.Parse(r)
	return libs
}

func (d *YarnDoctor) SourceCodeURL(lib types.Library) (string, error) {
	nodejs := Nodejs{lib: lib}
	url, err := nodejs.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}