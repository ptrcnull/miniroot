package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var rootPath = flag.String("root", "", "path to root directory")
var initCmd = flag.String("init", "/sbin/init", "init command")
var workdir = flag.String("workdir", "/", "work directory")
var pidns = flag.Bool("pidns", false, "create pid namespace")

func main() {
	flag.Parse()

	if *rootPath == "" {
		fmt.Println("root directory not set")
		os.Exit(1)
	}

	initParts := strings.Split(*initCmd, " ")
	cmd := exec.Command(initParts[0], initParts[1:]...)
	cmd.Dir = *workdir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot:       *rootPath,
		Cloneflags:   syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}
	if *pidns {
		cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWPID
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Printf("failed to start: %s\n", err)
	}

	sigs := make(chan os.Signal, 1)
	go func() {
		for {
			_ = <-sigs
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}()
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("process exited with error: %s\n", err)
	}
}
