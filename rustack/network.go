package rustack

import (
	"fmt"
)

type Network struct {
	manager   *Manager
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	Vdc       struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"vdc"`
}

func NewNetwork(name string) Network {
	n := Network{Name: name}
	return n
}

func (m *Manager) GetNetworks(extraArgs ...Arguments) (networks []*Network, err error) {
	args := Defaults()
	args.merge(extraArgs)

	path := "v1/network"
	err = m.GetItems(path, args, &networks)
	for i := range networks {
		networks[i].manager = m
	}
	return
}

func (v *Vdc) GetNetworks(extraArgs ...Arguments) (networks []*Network, err error) {
	args := Arguments{
		"vdc": v.ID,
	}

	args.merge(extraArgs)
	networks, err = v.manager.GetNetworks(args)
	return
}

func (m *Manager) GetNetwork(id string) (network *Network, err error) {
	path := fmt.Sprintf("v1/network/%s", id)
	err = m.Get(path, Defaults(), &network)
	if err != nil {
		return
	}
	network.manager = m
	return
}

func (n *Network) CreateSubnet(subnet *Subnet) error {
	path := fmt.Sprintf("v1/network/%s/subnet", n.ID)
	err := n.manager.Post(path, subnet, &subnet)
	if err == nil {
		subnet.manager = n.manager
		subnet.network = n
	}

	return err
}

func (n *Network) Rename(name string) error {
	path := fmt.Sprintf("v1/network/%s", n.ID)
	return n.manager.Put(path, Arguments{"name": name}, n)
}

func (n *Network) GetSubnets() (subnets []*Subnet, err error) {
	path := fmt.Sprintf("v1/network/%s/subnet", n.ID)
	err = n.manager.GetItems(path, Arguments{}, &subnets)
	for i := range subnets {
		subnets[i].manager = n.manager
		subnets[i].network = n
	}

	return
}

func (n *Network) Delete() error {
	path := fmt.Sprintf("v1/network/%s", n.ID)
	return n.manager.Delete(path, Defaults(), n)
}
