package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vishvananda/netlink"
)

func main() {
	fmt.Printf("Dummy-int-init\n")
	fmt.Printf("Contact: cyril.connan@gmail.com\n")
	fmt.Printf("Project: https://github.com/ccleouf66/dummy-net-init\n")
	fmt.Printf("---\n\n")

	if len(os.Args) != 3 {
		log.Fatalf("Bad arguments number. \nUsage:\n   %s interface-names interface-ip", os.Args[0])
	}

	var intName string = os.Args[1]
	var intAddr string = os.Args[2]

	// Check if the dummy link already exist
	nl, err := netlink.LinkList()
	if err != nil {
		log.Fatalf("Error getting link list: %s\n", err)
	}
	for _, l := range nl {
		if l.Type() == "dummy" {
			if l.Attrs().Name == intName {

				// Get link addresses
				laddr, err := netlink.AddrList(l, 0)
				if err != nil {
					log.Fatalf("Error getting addresses for link %s: %s\n", l.Attrs().Name, err)
				}
				// Check all addresses
				for _, addr := range laddr {
					a, err := netlink.ParseAddr(intAddr)
					if err != nil {
						log.Fatalf("Error parsing addresses %s: \n%s\n", intAddr, err)
					}
					if addr.IPNet.IP.To4().Equal(a.IP) {
						log.Printf("Interface %s with address %s already exist.\n", intName, intAddr)
						return

					}
				}
				// Interface exist, add new address
				err = AddDummyAddr(intAddr, l.(*netlink.Dummy))
				if err != nil {
					log.Fatalf("%v\n", err)
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
		log.Fatalf("Could not create %s interface, are you priviliged ?:\n%v\n", la.Name, err)
	}

	// Add address to dummy interface
	err = AddDummyAddr(intAddr, dummyInt)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

// AddDummyAddr add addr to dummy interface
func AddDummyAddr(addr string, dummy *netlink.Dummy) (err error) {
	a, err := netlink.ParseAddr(addr)
	if err != nil {
		log.Printf("Could not create addr %s\n", addr)
		return err
	}
	err = netlink.AddrAdd(dummy, a)
	if err != nil {
		log.Printf("Could not add address %s to interface %s\n", addr, dummy.Attrs().Name)
		return err
	}
	return nil
}
