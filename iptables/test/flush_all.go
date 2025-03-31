package main

import (
	"fmt"

	"github.com/kocmo/go-xtables/iptables"
)

func FlushALL() {
	set()
	defer unset()

	err := iptables.NewIPTables().Flush()
	fmt.Println(err)
}
