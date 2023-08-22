package containers

import (
	"fmt"
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

// 1 go routine for each remote host
// then 1 go routine for each container in order to get the logs

func NewSSH(remoteHosts types.Config) []*SSH {
	assh := make([]*SSH, len(remoteHosts.Hosts))
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
		ssh.Channels[host.IP] = MapChannels{
			Host:     host.IP,
			Command:  make(chan string, 10),
			Response: make(chan []byte, 10),
		}
		go ssh.OpenConnection(ssh.Channels[host.IP])
		assh = append(assh, &ssh)
	}
	return assh
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

	in, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	out, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	err = session.Shell()

	// FIXME: Send the command and store each container and log in a map to be used later
	for {
		select {
		case command := <-channel.Command:
			in.Write([]byte(command + "\n"))

			buf := make([]byte, 1024*1024)
			n, err := out.Read(buf)
			if err != nil {
				panic(err)
			}

			fmt.Println("Hello", string(buf[:n]))

			channel.Response <- buf[:n]
		}
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
