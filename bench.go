package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/parallelcointeam/pod/fork"
	"time"
)

var (
	sha256Reps = int(1 << 24)
	scryptReps = int(1 << 16)
	hf1Reps    = int(1 << 9)
)

// Bench runs benchmarks on all algorithms for each hardfork level
func Bench() {
	fmt.Println("Benchmark requested")
	fmt.Println("Please turn off any high cpu processes for a more accurate benchmark")
	fmt.Println("Pre-HF1 benchmarks:")
	fork.IsTestnet = false
	for a := range fork.List[0].AlgoVers {
		fmt.Println("Benchmarking algo", fork.List[0].AlgoVers[a])
		var speed int64
		switch a {
		case 2:
			speed = bench(0, fork.List[0].AlgoVers[a], sha256Reps)
		case 514:
			speed = bench(0, fork.List[0].AlgoVers[a], scryptReps)
		}
		fmt.Println(speed, "ns/hash", speed)
	}
	fmt.Println("HF1 benchmarks:")
	fork.IsTestnet = true
	for a := range fork.List[1].AlgoVers {
		fmt.Println("Benchmarking algo", fork.List[1].AlgoVers[a])
		speed := bench(1, fork.List[1].AlgoVers[a], hf1Reps)
		fmt.Println(speed/1000, "μs/hash", speed)
	}
}

func bench(hf int, algo string, reps int) int64 {
	startTime := time.Now()
	b := make([]byte, 80)
	rand.Read(b)
	// Zero out the nonce value
	for i := 76; i < 80; i++ {
		b[i] = 0
	}
	var height int32
	if hf == 1 {
		height = fork.List[1].ActivationHeight
	}
	var done bool
	var i int
	for i = 0; i < reps && !done; i++ {
		b, done = updateNonce(b)
		fork.Hash(b, algo, height)
	}
	endTime := time.Now()
	return int64(endTime.Sub(startTime)) / int64(i)
}

func updateNonce(b []byte) (out []byte, done bool) {
	if len(b) < 80 {
		return
	}
	nonce := binary.LittleEndian.Uint32(b[76:80])
	nonce++
	if nonce == 1<<31 {
		done = true
		return
	}
	nonceBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonceBytes, nonce)
	for i := range b[76:80] {
		b[76+i] = nonceBytes[i]
	}
	return b, done
}
