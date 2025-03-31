package main

import (
	"github.com/kocmo/go-xtables/iptables"
	"github.com/kocmo/go-xtables/pkg/network"
)

func AntiPING() {
	set()
	defer unset()

	iptables.NewIPTables().
		Table(iptables.TableTypeFilter).
		Chain(iptables.ChainTypeINPUT).
		MatchProtocol(false, network.ProtocolICMP).
		MatchICMP(false, network.ICMPType(network.EchoRequest)).
		TargetDrop().
		Append()
}
