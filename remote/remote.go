package remote

import (
	"encoding/json"
	"os"
	"os/user"
	"strings"

	"github.com/joho/godotenv"
	"github.com/victorguidi/TermDockerCLI/containers"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var (
	u, _ = user.Current()
)

type SSH struct {
	Host   string
	Port   string
	config *ssh.ClientConfig
}

func init() {
	godotenv.Load()
}

// TODO: Adapt this to read from a list of hosts, maybe on yml
func NewSSH() *SSH {
	var hostKeyCallback ssh.HostKeyCallback
	hostKeyCallback, err := knownhosts.New("/home/" + u.Username + "/.ssh/known_hosts")
	if err != nil {
		panic(err)
	}
	return &SSH{
		Host: os.Getenv("REMOTE_HOST"),
		Port: "22",
		config: &ssh.ClientConfig{
			User: os.Getenv("REMOTE_USER"),
			Auth: []ssh.AuthMethod{
				ssh.Password(os.Getenv("REMOTE_PASSWORD")),
			},
			// HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
			HostKeyCallback:   hostKeyCallback,
			HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
		},
	}
}

func (s *SSH) GetContainerFromRemote(c chan<- []containers.DockerContainer) {
	defer close(c)

	cmdToRun := `docker ps -a --format '{"ID":"{{.ID}}", "Image":"{{.Image}}"}' | jq -s .`
	client, err := ssh.Dial("tcp", s.Host+":"+s.Port, s.config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}

	defer session.Close()

	var b []byte
	b, err = session.Output(cmdToRun)
	if err != nil {
		panic("Failed to run: " + err.Error())
	}
	var containers []containers.DockerContainer
	decoder := json.NewDecoder(strings.NewReader(string(b)))
	if err := decoder.Decode(&containers); err != nil {
		panic("Error decoding JSON:" + err.Error())
	}

	c <- containers

}

func (s *SSH) GetImagesFromRemote() {}
