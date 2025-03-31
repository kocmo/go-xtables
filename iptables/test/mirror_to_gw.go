package main

import (
	"fmt"
	"net"

	"github.com/kocmo/go-xtables/iptables"
	"github.com/kocmo/go-xtables/pkg/network"
)

func MirrorToGW() {
	set()
	defer unset()

	err := iptables.NewIPTables().
		Table(iptables.TableTypeMangle).
		Chain(iptables.ChainTypePREROUTING).
		MatchProtocol(false, network.ProtocolTCP).
		MatchTCP(iptables.WithMatchTCPDstPort(false, 2432)).
		TargetTEE(net.ParseIP("192.168.1.1")).
		Insert()
	fmt.Println(err)
}
