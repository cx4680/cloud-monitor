package config

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var cf = flag.String("config.yml", "../config.yml.local.yml", "config.yml path")

func TestInitConfig(t *testing.T) {
	flag.Parse()
	err := InitConfig(*cf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init config.yml error: %v", err)
		os.Exit(1)
	}
}
