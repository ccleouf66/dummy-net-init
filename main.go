package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

var intName string = "s3-lb"
var intAddr string = "169.254.20.10/24"

func main() {

	// Check if the dummy link already exist
	nl, err := netlink.LinkList()
	if err != nil {
		fmt.Printf("Error getting link list: %s\n", err)
		return
	}
	for _, l := range nl {
		if l.Type() == "dummy" {
			if l.Attrs().Name == intName {

				// Get link addresses
				laddr, err := netlink.AddrList(l, 0)
				if err != nil {
					fmt.Printf("Error getting addresses for link %s: %s\n", l.Attrs().Name, err)
					return
				}
				// Check all addresses
				for _, addr := range laddr {
					a, err := netlink.ParseAddr(intAddr)
					if err != nil {
						fmt.Printf("Error parsing addresses %s: \n%s\n", intAddr, err)
					}
					if addr.IPNet.IP.To4().Equal(a.IP) {
						fmt.Printf("Interface %s with address %s already exist.\n", intName, intAddr)
						return
					}
				}
				// Interface exist, add new address
				err = AddDummyAddr(intAddr, l.(*netlink.Dummy))
				if err != nil {
					fmt.Printf("\n%v\n", err)
				}
				return
			}
		}
	}

	// Configure new dummy interface
	la := netlink.NewLinkAttrs()
	la.Name = intName
	dummyInt := &netlink.Dummy{LinkAttrs: la}

	// Create dummy interface
	err = netlink.LinkAdd(dummyInt)
	if err != nil {
		fmt.Printf("Could not create %s interface, are you priviliged ?:\n%v\n", la.Name, err)
		return
	}

	// Add address to dummy interface
	err = AddDummyAddr(intAddr, dummyInt)
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
}

// AddDummyAddr add addr to dummy interface
func AddDummyAddr(addr string, dummy *netlink.Dummy) (err error) {
	a, err := netlink.ParseAddr(addr)
	if err != nil {
		fmt.Printf("Could not create addr %s\n", addr)
		return err
	}
	err = netlink.AddrAdd(dummy, a)
	if err != nil {
		fmt.Printf("Could not add address %s to interface %s\n", addr, dummy.Attrs().Name)
		return err
	}
	return nil
}
