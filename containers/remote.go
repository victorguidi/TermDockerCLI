package containers

import (
	"encoding/json"
	"os/user"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/victorguidi/TermDockerCLI/types"
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
func NewSSH(remoteHosts types.Config) {
	var hostKeyCallback ssh.HostKeyCallback
	hostKeyCallback, err := knownhosts.New("/home/" + u.Username + "/.ssh/known_hosts")
	if err != nil {
		panic(err)
	}

	for _, host := range remoteHosts.Hosts {
		command := make(chan string)
		response := make(chan []byte)
		ssh := SSH{
			Host: host.IP,
			Port: "22",
			config: &ssh.ClientConfig{
				User: host.User,
				Auth: []ssh.AuthMethod{
					ssh.Password(host.Password),
				},
				// HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
				HostKeyCallback:   hostKeyCallback,
				HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
			},
		}
		go ssh.OpenConnection(command, response)
	}
}

func (s *SSH) OpenConnection(command chan string, response chan []byte) {
	client, err := ssh.Dial("tcp", s.Host+":"+s.Port, s.config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	for {
		command := <-command
		var b []byte
		b, err = session.Output(command)
		if err != nil {
			panic("Failed to run: " + err.Error())
		}
		response <- b
	}
}

func (s *SSH) GetContainerFromRemote(c chan<- []DockerContainer) {
	defer close(c)

	cmdToRun := `docker ps -a --format '{"ContainerID":"{{.ID}}", "Image":"{{.Image}}"}' | jq -s .`
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
	var containers []DockerContainer
	decoder := json.NewDecoder(strings.NewReader(string(b)))
	if err := decoder.Decode(&containers); err != nil {
		panic("Error decoding JSON:" + err.Error())
	}

	c <- containers
}

func (s *SSH) GetRemoteLogs(c chan<- []byte, containerId string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmdToRun := `docker logs ` + containerId
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

	c <- b
}

func (s *SSH) GetImagesFromRemote() {}
