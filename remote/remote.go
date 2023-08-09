package remote

import (
	"os"

	"github.com/joho/godotenv"
)

type SSH struct {
	Host     string
	User     string
	Password string
	Port     int
}

func init() {
	godotenv.Load()
}

func NewSSH() *SSH {
	return &SSH{
		Host:     os.Getenv("REMOTE_HOST"),
		User:     os.Getenv("REMOTE_USER"),
		Password: os.Getenv("REMOTE_PASSWORD"),
		Port:     22,
	}
}

func (s *SSH) GetContainerFromRemote() {
	// cmdToRun := `docker ps -a --format '{"ID":"{{.ID}}", "Image":"{{.Image}}"}' | jq -s .`
}

func (s *SSH) GetImagesFromRemote() {}
