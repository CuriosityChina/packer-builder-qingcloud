package qingcloud

import (
	"github.com/mitchellh/multistep"
)

type stepCreateImage struct {
}

func (s *stepCreateImage) Run(state multistep.StateBag) multistep.StepAction {
	return multistep.ActionContinue
}

func (s *stepCreateImage) Cleanup(state multistep.StateBag) {
}
