package qingcloud

import (
	"fmt"
	"github.com/hashicorp/packer/common/uuid"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

type stepKeypair struct {
}

func (s *stepKeypair) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)

	ui.Say("Create SSH key pair...")

	keypairClient, err := client.KeyPair(c.Zone)
	kp, err := keypairClient.CreateKeyPair(&qc.CreateKeyPairInput{
		Mode:          qc.String("system"),
		EncryptMethod: qc.String("ssh-rsa"),
		KeyPairName:   qc.String(fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())),
	})
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("keypair_id", *kp.KeyPairID)
	state.Put("ssh_private_key", *kp.PrivateKey)

	return multistep.ActionContinue
}

func (s *stepKeypair) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	client := state.Get("client").(*qc.QingCloudService)
	kpID := state.Get("keypair_id").(string)

	ui.Say("Destory SSH key pair...")
	keypairClient, err := client.KeyPair(c.Zone)
	keypairClient.DeleteKeyPairs(&qc.DeleteKeyPairsInput{
		KeyPairs: []*string{qc.String(kpID)},
	})

	if err != nil {
		ui.Error(err.Error())
		return
	}
	ui.Say("SSH keypair destoried!")
}
