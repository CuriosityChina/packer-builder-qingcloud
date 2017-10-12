package qingcloud

import (
	"fmt"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"time"
)

type stepCreateImage struct {
}

func (s *stepCreateImage) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)

	insID := state.Get("instance_id").(string)

	instanceClient, err := client.Instance(c.Zone)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Stop Instance...")
	_, err = instanceClient.StopInstances(&qc.StopInstancesInput{
		Instances: []*string{qc.String(insID)},
	})
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if err := waitForInstanceState("stopped", insID, instanceClient, 30*time.Second); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Capture Instance...")
	imageClient, err := client.Image(c.Zone)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	imageName := qc.String(fmt.Sprintf("packer-%s", time.Now().String()[:10]))
	img, err := imageClient.CaptureInstance(&qc.CaptureInstanceInput{
		Instance:  qc.String(insID),
		ImageName: imageName,
	})
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("final_image_id", *img.ImageID)
	state.Put("final_image_name", *imageName)

	return multistep.ActionContinue
}

func (s *stepCreateImage) Cleanup(state multistep.StateBag) {
}
