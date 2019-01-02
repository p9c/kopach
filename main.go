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
	var benches []Benchmark
	if cfg.Bench {
		benchmark, benches = Bench()
		f, err := os.Create(cfg.DataDir + "/benchmark.json")
		if err != nil {
			fmt.Println("ERROR: unable to write benchmark results")
			os.Exit(1)
		}
		f.WriteString(benchmark)
		fmt.Println("Benchmark data saved")
		f.Close()
	}
	j, _ := json.MarshalIndent(cfg, "  ", "  ")
	fmt.Println(string(j))
	fmt.Println(benches)
}
