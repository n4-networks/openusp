package clitest

import (
	"testing"
	"time"
)

func TestIp_ShowIntf(t *testing.T) {
	cmd := "show ip intf"
	runAndCheckErr(t, cmd)
}

func TestIp_AddShowRemoveIntf(t *testing.T) {
	t.Log("Adding an interface to IP")
	cmd := "add ip intf alias testIntf type normal version v4 lowerlayer eth id 1"
	runAndCheckErr(t, cmd)

	time.Sleep(time.Second)

	t.Log("Show IP interface")
	cmd = "show ip intf"
	runAndCheckErr(t, cmd)

	t.Log("Removing the added IP Interface")
	cmd = "remove ip intf testIntf"
	runAndCheckErr(t, cmd)
}

func TestIp_ShowIntfAddr(t *testing.T) {
	cmd := "show ip intf ipv4addr"
	runAndCheckErr(t, cmd)
}

func TestIp_AddShowRemoveIntfAddr(t *testing.T) {

	// 1. Create an IP Interface
	t.Log("Adding an IP Interface")
	cmd := "add ip intf alias testIntf type normal version v4 lowerlayer eth id 1"
	runAndCheckErr(t, cmd)

	// 2. Add IPv4 address to the created interface
	t.Log("Adding address to IP Interface")
	cmd = "add ip addr alias myaddr intf testIntf type v4 mode static address 182.158.2.2 subnet \\24"
	runAndCheckErr(t, cmd)

	// 3. Remove the added address
	t.Log("Removing the added IP Address")
	cmd = "remove ip intf testIntf ipv4addr myaddr"
	runAndCheckErr(t, cmd)

	// 4. Remove the added interface
	cmd = "remove ip intf testIntf"
	runAndCheckErr(t, cmd)
}

func TestIp_CfgAddShowRemoveIntf(t *testing.T) {

	t.Log("Adding Interface to cfg table for IP")
	cmd := "addcfg ip intf alias testIntf type normal version v4 lowerlayer eth id 1"
	runAndCheckErr(t, cmd)

	cmd = "showcfg ip intf"
	runAndCheckErr(t, cmd)

	t.Log("Removing added IP Interface from cfg table")
	cmd = "removecfg ip intf testIntf"
	runAndCheckErr(t, cmd)
}

func TestIp_CfgAddShowRemoveIntfAddr(t *testing.T) {
	// Add an IP Interface
	t.Log("Adding IP Interface to cfg store")
	cmd := "addcfg ip intf alias testIntf type normal version v4 lowerlayer eth id 1"
	runAndCheckErr(t, cmd)

	// Add IPv4 address to the created interface
	t.Log("Adding address named IP Interface in the cfg store")
	cmd = "addcfg ip addr alias myaddr intf testIntf type v4 mode static address 182.158.2.2 subnet \\24"
	runAndCheckErr(t, cmd)

	cmd = "showcfg ip intf ipv4addr"
	runAndCheckErr(t, cmd)

	// Remove the added address
	t.Log("Removing the added IP Address from cfg store")
	cmd = "removecfg ip intf testIntf ipv4addr myaddr"
	runAndCheckErr(t, cmd)

	// 4. Remove the added interface
	t.Log("Removing the added IP interace from cfg store")
	cmd = "removecfg ip intf testIntf"
	runAndCheckErr(t, cmd)
}
