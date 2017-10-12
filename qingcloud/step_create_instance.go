package qingcloud

import (
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

type stepCreateInstance struct {
}

func (s *stepCreateInstance) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)

	ui.Say("Create Instance...")
	instanceClient, err := client.Instance(c.Zone)
	if err != nil {
		ui.Error(err.Error())
	}

	// Fix: 主动造成一个错误
	_, err = instanceClient.RunInstances(&qc.RunInstancesInput{})
	if err != nil {
		ui.Error(err.Error())
	}

	return multistep.ActionContinue
}

// Cleanup TODO:
func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {

}
