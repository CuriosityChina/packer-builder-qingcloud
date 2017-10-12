package qingcloud

import (
	"log"

	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

var _ packer.Builder = &Builder{}

// Builder ...
type Builder struct {
	config Config
	runner multistep.Runner
}

// Prepare ...
func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	c, warnings, errs := NewConfig(raws...)
	if errs != nil {
		return warnings, errs
	}
	b.config = *c
	return nil, nil
}

// Run ...
func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	return nil, nil
}

// Cancel ...
func (b *Builder) Cancel() {
	log.Printf("builder stop %s", "done")
}
