package images

import (
	"encoding/json"
	"log"
	"os/exec"
)

type DockerImage struct {
	Repository string
	Tag        string
	ImageId    string
	Created    string
	Size       string
}

func NewImage() *DockerImage {
	return &DockerImage{}
}

func (i *DockerImage) GetImages(ci chan []DockerImage) {
	defer close(ci)

	cmd := exec.Command("docker", "images", "--format", `{
    "Repository": "{{.Repository}}",
    "Tag": "{{.Tag}}",
    "ImageId": "{{.ID}}",
    "Size": "{{.Size}}"
    }`)
	jqCmd := exec.Command("jq", "-s", ".")

	jqCmd.Stdin, _ = cmd.StdoutPipe()
	jqOutput, err := jqCmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating StdoutPipe for jqCmd", err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting cmd", err)
	}
	if err := jqCmd.Start(); err != nil {
		log.Fatal("Error starting jqCmd", err)
	}

	var images []DockerImage
	decoder := json.NewDecoder(jqOutput)
	if err := decoder.Decode(&images); err != nil {
		log.Fatal("Error decoding json", err)
	}

	if err := jqCmd.Wait(); err != nil {
		log.Fatal("Error waiting for jqCmd", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal("Error waiting for cmd", err)
	}

	ci <- images

}
