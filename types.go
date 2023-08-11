package main

import "github.com/rivo/tview"

type FlexBox struct {
	*tview.Flex
	Tabs [3]*tview.Table
	Data chan any
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

type Setter interface {
	Set([]string)
}
