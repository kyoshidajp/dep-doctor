package erlang_elixir

import (
	"net/http"

	"github.com/aquasecurity/go-dep-parser/pkg/hex/mix"
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type MixDoctor struct {
	HTTPClient http.Client
}

func NewMixDoctor() *MixDoctor {
	client := &http.Client{}
	return &MixDoctor{HTTPClient: *client}
}

func (d *MixDoctor) Parse(r parser_io.ReadSeekerAt) (types.Libraries, error) {
	p := &mix.Parser{}
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *MixDoctor) SourceCodeURL(lib types.Library) (string, error) {
	hex := Hex{name: lib.Name}
	url, err := hex.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
