package remote

type SSH struct {
	Host     string
	User     string
	Password string
	Port     int
}

func NewSSH() *SSH {
	return &SSH{}
}

func (s *SSH) GetContainerFromRemote() {
	// cmdToRun := `docker ps -a --format '{"ID":"{{.ID}}", "Image":"{{.Image}}"}' | jq -s .`
}

func (s *SSH) GetImagesFromRemote() {}
