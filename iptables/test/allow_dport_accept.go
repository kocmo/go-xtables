package main

import (
	"fmt"

	"github.com/kocmo/go-xtables/iptables"
	"github.com/kocmo/go-xtables/pkg/network"
)

func AllowDPortAccept() {
	set()
	defer unset()

	err := iptables.NewIPTables().
		Table(iptables.TableTypeFilter).
		Chain(iptables.ChainTypeINPUT).
		MatchProtocol(false, network.ProtocolTCP).
		MatchTCP(iptables.WithMatchTCPDstPort(false, 2432)).
		TargetAccept().
		Append()
	fmt.Println(err)
}
