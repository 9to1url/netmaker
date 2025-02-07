package local

import (
	"net"

	"github.com/gravitl/netmaker/logger"
	"github.com/gravitl/netmaker/netclient/ncutils"
	"github.com/seancfoley/ipaddress-go/ipaddr"
)

func setRoute(iface string, addr *net.IPNet, address string) error {
	var err error
	_, _ = ncutils.RunCmd("route add -net "+addr.String()+" -interface "+iface, false)
	return err
}

func deleteRoute(iface string, addr *net.IPNet, address string) error {
	var err error
	_, _ = ncutils.RunCmd("route delete -net "+addr.String()+" -interface "+iface, false)
	return err
}

func setCidr(iface, address string, addr *net.IPNet) {
	cidr := ipaddr.NewIPAddressString(addr.String()).GetAddress()
	if cidr.IsIPv4() {
		ncutils.RunCmd("route add -net "+addr.String()+" -interface "+iface, false)
	} else if cidr.IsIPv6() {
		ncutils.RunCmd("route add -net -inet6 "+addr.String()+" -interface "+iface, false)
	} else {
		logger.Log(1, "could not parse address: "+addr.String())
	}
	ncutils.RunCmd("route add -net "+addr.String()+" -interface "+iface, false)
}

func removeCidr(iface string, addr *net.IPNet, address string) {
	ncutils.RunCmd("route delete -net "+addr.String()+" -interface "+iface, false)
}
