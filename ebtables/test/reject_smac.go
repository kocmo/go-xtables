package main

import (
	"fmt"

	"github.com/kocmo/go-xtables/ebtables"
)

func main() {
	set()
	defer unset()

	err := ebtables.NewEBTables().
		Table(ebtables.TableTypeFilter).
		Chain(ebtables.ChainTypeINPUT).
		MatchSource(false, "00:11:22:33:44:55").
		TargetDrop().
		Append()
	fmt.Println(err)
}
