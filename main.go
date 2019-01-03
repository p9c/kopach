package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/parallelcointeam/pod/fork"
	"io/ioutil"
	"os"
)

func getBenchPath(cfg *config) string {
	return cfg.DataDir + "/benchmark.json"
}

func main() {
	var benches *[]Benchmark
	var benchmarkJSON string
	fmt.Println("Kopach CPU miner for Parallelcoin DUO")
	cfg, args, err := loadConfig()
	if err != nil {
		os.Exit(1)
	}
	_, _ = cfg, args
	if cfg.TestNet3 {
		fork.IsTestnet = true
	}
	if cfg.Bench {
		benchmarkJSON, *benches = Bench()
		f, err := os.Create(getBenchPath(cfg))
		if err != nil {
			fmt.Println("ERROR: unable to write benchmark results", err.Error())
			os.Exit(1)
		}
		f.WriteString(benchmarkJSON)
		fmt.Println("Benchmark data saved")
		f.Close()
		os.Exit(0)
	}
	if cfg.Algo == "easy" {
		if _, err := os.Stat(getBenchPath(cfg)); os.IsNotExist(err) {
			fmt.Println("Running benchmark for 'easy' mining mode")
			benchmarkJSON, *benches = Bench()
			f, err := os.Create(getBenchPath(cfg))
			if err != nil {
				fmt.Println("ERROR: unable to write benchmark results", err.Error())
				os.Exit(1)
			}
			f.WriteString(benchmarkJSON)
			fmt.Println("Benchmark data saved")
			f.Close()
		}

		fmt.Println("Mining in 'easy' mode, targeting lowest difficulty algorithm")
	}
	benchFile, err := os.Open(getBenchPath(cfg))
	if err != nil {
		fmt.Println("ERROR: unable to open benchmark file", err.Error())
		os.Exit(1)
	}
	defer benchFile.Close()
	reader := bufio.NewReader(benchFile)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("ERROR: could not read benchmark file", err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(b, &benches)
	if err != nil {
		fmt.Println("ERROR: unable to decode benchmark file, deleting", err.Error())
		err = os.Remove(getBenchPath(cfg))
		if err != nil {
			fmt.Println("ERROR: unable to delete file", err.Error())
		}
		os.Exit(1)
	}
	if benches != nil {
		fmt.Println("Loaded benchmark data")
	}
}
