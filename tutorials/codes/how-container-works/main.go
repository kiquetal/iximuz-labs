package main

import "os"
import "os/exec"
import "syscall"

func main() {
	if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
		panic(err)
	}

	cmd := exec.Command("sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	cmd.Run()
}
