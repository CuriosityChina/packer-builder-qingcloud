package qingcloud

import (
	"fmt"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"time"
)

type stepAttachEIP struct {
}

func (s *stepAttachEIP) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)
	insID := state.Get("instance_id").(string)
	eipID := state.Get("eip_id").(string)

	ui.Say("Associate EIP...")

	ipClient, err := client.EIP(c.Zone)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if err := waitForEIPState("available", eipID, ipClient, 30*time.Second); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	_, err = ipClient.AssociateEIP(&qc.AssociateEIPInput{
		Instance: qc.String(insID),
		EIP:      qc.String(eipID),
	})
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if err := waitForEIPState("associated", eipID, ipClient, 30*time.Second); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	eips, err := ipClient.DescribeEIPs(&qc.DescribeEIPsInput{
		InstanceID: qc.String(insID),
	})
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if len(eips.EIPSet) != 1 {
		ui.Error("not found EIP")
		return multistep.ActionHalt
	}

	state.Put("eip_addr", *eips.EIPSet[0].EIPAddr)
	ui.Say(fmt.Sprintf("EIP Address is %s", *eips.EIPSet[0].EIPAddr))

	return multistep.ActionContinue
}

func (s *stepAttachEIP) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)
	eipID := state.Get("eip_id").(string)

	ui.Say("Dissociate EIPs...")
	ipClient, err := client.EIP(c.Zone)
	if err != nil {
		ui.Error(err.Error())
		return
	}

	_, err = ipClient.DissociateEIPs(&qc.DissociateEIPsInput{
		EIPs: []*string{qc.String(eipID)},
	})

	if err != nil {
		ui.Error(err.Error())
		return
	}

	if err := waitForEIPState("available", eipID, ipClient, 30*time.Second); err != nil {
		ui.Error(err.Error())
	}

	ui.Say("EIP dissociated!")
}
