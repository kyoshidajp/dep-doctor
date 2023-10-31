package cmd

import (
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
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

func (d *CargoDoctor) Libraries(r parser_io.ReadSeekerAt) []types.Library {
	p := &cargo.Parser{}
	libs, _, _ := p.Parse(r)
	return libs
}

func (d *CargoDoctor) SourceCodeURL(lib types.Library) (string, error) {
	crate := Crate{name: lib.Name}
	url, err := crate.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
