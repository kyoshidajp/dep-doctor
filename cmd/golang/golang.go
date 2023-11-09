package golang

import (
	"net/http"

	"github.com/aquasecurity/go-dep-parser/pkg/golang/mod"
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type GolangDoctor struct {
	HTTPClient http.Client
}

func NewGolangDoctor() *GolangDoctor {
	client := &http.Client{}
	return &GolangDoctor{HTTPClient: *client}
}

func (d *GolangDoctor) Parse(r parser_io.ReadSeekerAt) (types.Libraries, error) {
	p := &mod.Parser{}
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *GolangDoctor) SourceCodeURL(lib types.Library) (string, error) {
	proxyGolang := ProxyGolang{lib: lib}
	if len(lib.ExternalReferences) > 0 {
		return lib.ExternalReferences[0].URL, nil
	}

	url, err := proxyGolang.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
