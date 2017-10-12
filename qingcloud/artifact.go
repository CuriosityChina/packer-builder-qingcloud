package qingcloud

import (
	"fmt"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"log"
)

// Artifact ...
type Artifact struct {
	imageName string
	imageID   string
	zone      string
	// The client for making API calls
	service *qc.QingCloudService
}

// BuilderId ...
func (*Artifact) BuilderId() string {
	return BuilderId
}

// Files ...
func (*Artifact) Files() []string {
	return nil
}

// Id ...
func (a *Artifact) Id() string {
	return fmt.Sprintf("%s", a.imageID)
}

func (a *Artifact) String() string {
	return fmt.Sprintf("A image was created: '%v' (ID: %v)", a.imageName, a.imageID)
}

// State ...
func (a *Artifact) State(name string) interface{} {
	return nil
}

// Destroy ...
func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %s (%s)", a.imageID, a.imageName)
	imageClient, err := a.service.Image(a.zone)
	if err != nil {
		return err
	}

	_, err = imageClient.DeleteImages(&qc.DeleteImagesInput{
		Images: []*string{qc.String(a.imageID)},
	})

	return err
}
