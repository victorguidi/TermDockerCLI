package main

type IGeneric interface {
	DockerContainer | DockerImage | DockerNetwork
}

type DockerContainer struct {
	ContainerId string
	Image       string
	Command     string
	Created     string
	Status      string
	Ports       string
	Names       string
}

type DockerImage struct {
	Repository string
	Tag        string
	ImageId    string
	Created    string
	Size       string
}

type DockerNetwork struct {
	NetworkId string
	Name      string
	Driver    string
	Scope     string
}

func (d *DockerContainer) Set(fields []string) {
	d.ContainerId = fields[0]
	d.Image = fields[1]
	d.Command = fields[2]
	d.Created = fields[3]
	d.Status = fields[4]
	d.Ports = fields[5]
	d.Names = fields[6]
}

func (d *DockerImage) Set(fields []string) {
	d.Repository = fields[0]
	d.Tag = fields[1]
	d.ImageId = fields[2]
	d.Created = fields[3]
	d.Size = fields[4]
}

func (d *DockerNetwork) Set(fields []string) {
	d.NetworkId = fields[0]
	d.Name = fields[1]
	d.Driver = fields[2]
	d.Scope = fields[3]
}

type Setter interface {
	Set([]string)
}
