package images

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
