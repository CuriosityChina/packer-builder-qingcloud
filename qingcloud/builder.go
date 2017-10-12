package qingcloud

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	"github.com/yunify/qingcloud-sdk-go/config"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"log"
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
	cfg, err := config.New(b.config.APIKey, b.config.APISecret)
	if err != nil {
		return nil, err
	}

	client, err := qc.Init(cfg)
	if err != nil {
		return nil, err
	}

	// Setup
	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("client", client)
	state.Put("ui", ui)

	// Build
	steps := []multistep.Step{
		new(stepCreateInstance),
		// Run
		// ...
	}

	// Run
	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(state)

	// Output
	return nil, nil
}

// Cancel ...
func (b *Builder) Cancel() {
	log.Printf("builder stop %s", "done")
}
