package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const configPath = "/etc/mountmond.yaml"

func main() {
	configPathOverride := flag.String("config-path", configPath, "path to a YAML file containing configuration")
	flag.Parse()

	configFile, err := os.Open(*configPathOverride)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	mounts, err := ReadMountsFromConfig(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	watchDog := NewWatchDog(mounts)
	watchDog.Start()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT)

	<-signalChan
	watchDog.Shutdown()
}
