package main

import (
	"io"

	"gopkg.in/yaml.v2"
)

type MountDescription struct {
	Mount   string
	Command string
}

type Config struct {
	Mounts []MountDescription
}

func ReadMountsFromConfig(reader io.Reader) (map[string]string, error) {
	var config Config
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	mountsToCmds := make(map[string]string)
	for _, descr := range config.Mounts {
		mountsToCmds[descr.Mount] = descr.Command
	}

	return mountsToCmds, nil
}
