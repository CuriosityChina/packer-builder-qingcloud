package qingcloud

import (
	"fmt"
	"github.com/hashicorp/packer/common/uuid"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"time"
)

type stepCreateInstance struct {
}

func (s *stepCreateInstance) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)

	kpID := state.Get("keypair_id").(string)

	ui.Say("Create Instance...")
	instanceClient, err := client.Instance(c.Zone)
	if err != nil {
		ui.Error(err.Error())
	}

	ins, err := instanceClient.RunInstances(&qc.RunInstancesInput{
		ImageID:      qc.String(c.ImageID),
		InstanceType: qc.String(c.InstanceType),
		InstanceName: qc.String(fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())),
		LoginKeyPair: qc.String(kpID),
		LoginMode:    qc.String("keypair"),
		VxNets:       []*string{qc.String("vxnet-0")},
	})
	if err != nil {
		ui.Error(err.Error())
	}

	if len(ins.Instances) != 1 {
		ui.Error("Create Instance Failed!")
		return multistep.ActionHalt
	}

	state.Put("instance_id", *ins.Instances[0])

	waitForInstanceState("running", *ins.Instances[0], instanceClient, 30*time.Second)

	return multistep.ActionContinue
}

func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)
	insID := state.Get("instance_id").(string)

	ui.Say("Terminate Instance...")
	instanceClient, err := client.Instance(c.Zone)
	if err != nil {
		ui.Error(err.Error())
	}

	_, err = instanceClient.TerminateInstances(&qc.TerminateInstancesInput{
		Instances: []*string{qc.String(insID)},
	})
	if err != nil {
		ui.Error(err.Error())
		return
	}
	waitForInstanceState("terminated", insID, instanceClient, 30*time.Second)

	ui.Say("Instance terminate!")
}
