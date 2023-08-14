package containers

import (
	"os/user"
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
	Host     string
	Port     string
	config   *ssh.ClientConfig
	Channels map[string]MapChannels
}

type MapChannels struct {
	Host     string
	Command  chan string
	Response chan []byte
}

func init() {
	godotenv.Load()
}

func NewSSH(remoteHosts types.Config) {
	var hostKeyCallback ssh.HostKeyCallback
	hostKeyCallback, err := knownhosts.New("/home/" + u.Username + "/.ssh/known_hosts")
	if err != nil {
		panic(err)
	}

	for _, host := range remoteHosts.Hosts {
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
			Channels: make(map[string]MapChannels),
		}
		go ssh.OpenConnection(ssh.Channels[host.IP])
	}
}

func (s *SSH) OpenConnection(channel MapChannels) {
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
		command := <-channel.Command
		var b []byte
		b, err = session.Output(command)
		if err != nil {
			panic("Failed to run: " + err.Error())
		}
		channel.Response <- b
	}
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
