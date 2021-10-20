package config

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var cf = flag.String("config", "../config.local.yml", "config path")

func TestInitConfig(t *testing.T) {
	flag.Parse()
	err := InitConfig(*cf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init config error: %v", err)
		os.Exit(1)
	}
	cnf := GetConfig()
	fmt.Println(cnf)
}
