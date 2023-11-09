package php

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

func (d *ComposerDoctor) Parse(r parser_io.ReadSeekerAt) (types.Libraries, error) {
	p := composer.NewParser()
	libs, _, err := p.Parse(r)
	return libs, err
}

func (d *ComposerDoctor) SourceCodeURL(lib types.Library) (string, error) {
	packagist := Packagist{lib: lib}
	url, err := packagist.fetchURLFromRegistry(d.HTTPClient)
	return url, err
}
