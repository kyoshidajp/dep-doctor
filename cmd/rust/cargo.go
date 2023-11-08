package rust

import (
	"net/http"

	"github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/rust/cargo"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type CargoDoctor struct {
	HTTPClient http.Client
}

func NewCargoDoctor() *CargoDoctor {
	client := &http.Client{}
	return &CargoDoctor{HTTPClient: *client}
}

func (d *CargoDoctor) Parse(r io.ReadSeekerAt) (types.Libraries, error) {
	p := &cargo.Parser{}
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *CargoDoctor) SourceCodeURL(lib types.Library) (string, error) {
	crate := Crate{name: lib.Name}
	url, err := crate.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
