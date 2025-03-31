package main

import (
	"runtime"

	"github.com/vishvananda/netns"
)

var (
	originns netns.NsHandle
	newns    netns.NsHandle
)

func set() {
	// sandboxAddr := os.Getenv("SANDBOX_ADDR")
	// if sandboxAddr != "" {
	// 	sandboxUser := os.Getenv("SANDBOX_USER")
	// 	sandboxPassword := os.Getenv("SANDBOX_PASSWORD")

	// 	monkey.Patch(cmd.Cmd, func(name string, arg ...string) ([]byte, []byte, error) {
	// 		return cmd.SSHCmdPassword(sandboxAddr, sandboxUser, sandboxPassword,
	// 			name, arg...)
	// 	})
	// } else {
	runtime.LockOSThread()
	originns, _ = netns.Get()
	newns, _ = netns.New()
	// }
}

func unset() {
	// sandboxAddr := os.Getenv("SANDBOX_ADDR")
	// if sandboxAddr != "" {
	// 	monkey.UnpatchAll()
	// } else {
	runtime.UnlockOSThread()
	newns.Close()
	originns.Close()
	// }
}
