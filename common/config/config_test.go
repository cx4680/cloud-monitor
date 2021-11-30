package config

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var cf = flag.String("cfg.yml", "../cfg.yml.local.yml", "cfg.yml path")

func TestInitConfig(t *testing.T) {
	flag.Parse()
	err := InitConfig(*cf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init cfg.yml error: %v", err)
		os.Exit(1)
	}
}
