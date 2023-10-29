package cmd

import (
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/php/composer"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
)

type ComposerDoctor struct {
	HTTPClient http.Client
}

func NewComposerDoctor() *ComposerDoctor {
	client := &http.Client{}
	return &ComposerDoctor{HTTPClient: *client}
}

func (d *ComposerDoctor) Deps(r parser_io.ReadSeekerAt) []types.Library {
	p := composer.NewParser()
	deps, _, _ := p.Parse(r)
	return deps
}

func (d *ComposerDoctor) SourceCodeURL(name string) (string, error) {
	packagist := Packagist{name: name}
	url, err := packagist.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
