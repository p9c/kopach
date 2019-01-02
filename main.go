package main

import (
	"encoding/json"
	"fmt"
	"github.com/parallelcointeam/pod/fork"
	"os"
)

func main() {
	fmt.Println("Kopach CPU miner for Parallelcoin DUO")
	cfg, args, err := loadConfig()
	if err != nil {
		os.Exit(1)
	}
	_, _ = cfg, args
	if cfg.TestNet3 {
		fork.IsTestnet = true
	}
	var benchmark string
	if cfg.Bench {
		benchmark = Bench()
		fmt.Println(benchmark)
	}
	j, _ := json.MarshalIndent(cfg, "  ", "  ")
	fmt.Println(string(j))
}
