package qingcloud

import (
	"fmt"
	"github.com/hashicorp/packer/common/uuid"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"time"
)

type stepEIP struct {
}

func (s *stepEIP) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)

	ui.Say("Allocate EIP...")

	ipClient, err := client.EIP(c.Zone)
	// TODO: change to config
	ip, err := ipClient.AllocateEIPs(&qc.AllocateEIPsInput{
		BillingMode: qc.String("traffic"),
		Bandwidth:   qc.Int(4),
		EIPName:     qc.String(fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())),
	})
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if len(ip.EIPs) != 1 {
		ui.Error("can't AllocateEips")
		return multistep.ActionHalt
	}
	state.Put("eip_id", *ip.EIPs[0])

	waitForEIPState("available", *ip.EIPs[0], ipClient, 30*time.Second)

	return multistep.ActionContinue
}

func (s *stepEIP) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)
	ipID := state.Get("eip_id").(string)

	ui.Say("Destory IP...")
	ipClient, err := client.EIP(c.Zone)
	if err != nil {
		ui.Error(err.Error())
		return
	}

	waitForEIPState("available", ipID, ipClient, 30*time.Second)

	_, err = ipClient.ReleaseEIPs(&qc.ReleaseEIPsInput{
		EIPs: []*string{qc.String(ipID)},
	})

	if err != nil {
		ui.Error(err.Error())
		return
	}
	ui.Say("EIP released!")
}
