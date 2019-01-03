package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/parallelcointeam/pod/btcjson"
	"github.com/parallelcointeam/pod/fork"
	"github.com/parallelcointeam/pod/rpcclient"
	"github.com/parallelcointeam/pod/wire"
)

const (
	showHelpMessage = "Specify -h to show available options"
)

// commandUsage display the usage for a specific command.
func commandUsage(method string) {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usage, err := btcjson.MethodUsageText(method)
	if err != nil {
		// This should never happen since the method was already checked before calling this function, but be safe.
		fmt.Fprintln(os.Stderr, "Failed to obtain command usage:", err)
		return
	}
	fmt.Fprintln(os.Stderr, "\n", appName, Version(), "\nUsage:")
	fmt.Fprintf(os.Stderr, "  %s\n", usage)
}

// usage displays the general usage when the help flag is not displayed and and an invalid command was specified.  The commandUsage function is used instead when a valid command was specified.
func usage(errorMessage string) {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	fmt.Fprintln(os.Stderr, errorMessage)
	fmt.Fprintln(os.Stderr, appName, Version())
	fmt.Fprintln(os.Stderr, "\n", appName, Version(), "\nUsage:")
	fmt.Fprintf(os.Stderr, "  %s [OPTIONS] \n\n",
		appName)
	fmt.Fprintln(os.Stderr, showHelpMessage)
}

func main() {
	cfg, _, err := loadConfig()
	if err != nil {
		os.Exit(1)
	}
	url := defs.ParseURL(cfg.URL)
	global.Endpoints = append(global.Endpoints, url)
	switch {
	case !cfg.OneOnly:
		p := url.Port + 1
		if p == 0 {
			p = 11048
			url.Port = p
		}
		for i := p; i < p+8; i++ {
			global.Endpoints = append(global.Endpoints, defs.URL{
				Username: url.Username,
				Password: url.Password,
				Protocol: url.Protocol,
				Address:  url.Address,
				Port:     i,
			})
		}

		fallthrough
	case len(cfg.OtherPorts) > 0:
		for i := range cfg.OtherPorts {
			global.Endpoints = append(global.Endpoints, defs.URL{
				Username: url.Username,
				Password: url.Password,
				Protocol: url.Protocol,
				Address:  url.Address,
				Port:     cfg.OtherPorts[i],
			})
		}
	}
	r := global.Endpoints
	var e []defs.URL
	for i := range r {
		if _, err := net.Dial("tcp", fmt.Sprintf("%s:%d", r[i].Address, r[i].Port)); err == nil {
			e = append(e, r[i])
		} else {
			fmt.Println("ERROR", err)
		}
	}
	global.Endpoints = e
	connCfg := &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", global.Endpoints[0].Address, global.Endpoints[0].Port),
		User:         global.Endpoints[0].Username,
		Pass:         global.Endpoints[0].Password,
		HTTPPostMode: true,                                    // Bitcoin core only supports HTTP POST mode
		TLS:          global.Endpoints[0].Protocol == "https", // Bitcoin core does not provide TLS by default
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	currentNet, _ := client.GetCurrentNet()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	fmt.Print("Network ")
	switch currentNet {
	case wire.MainNet:
		fmt.Println("MainNet")
	case wire.TestNet:
		fmt.Println("TestNet")
	case wire.TestNet3:
		fmt.Println("TestNet3")
		fork.IsTestnet = true
	case wire.SimNet:
		fmt.Println("SimNet")
	}
	client.Shutdown()

	// var blocktemplates []string
	var responses []*btcjson.GetWorkResult
	for i := range global.Endpoints {
		// fmt.Println(global.Endpoints[i].String())
		// Connect to local bitcoin core RPC server using HTTP POST mode.
		connCfg := &rpcclient.ConnConfig{
			Host:         fmt.Sprintf("%s:%d", global.Endpoints[i].Address, global.Endpoints[i].Port),
			User:         global.Endpoints[i].Username,
			Pass:         global.Endpoints[i].Password,
			HTTPPostMode: true,                                    // Bitcoin core only supports HTTP POST mode
			TLS:          global.Endpoints[i].Protocol == "https", // Bitcoin core does not provide TLS by default
		}
		// Notice the notification parameter is nil since notifications are
		// not supported in HTTP POST mode.
		client, err := rpcclient.New(connCfg, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Shutdown()
		res, _ := client.GetWork()
		if err != nil {
			fmt.Println("ERROR", err)
		}
		// bt, _ := hex.DecodeString(res.Data)
		// algo, _ := strconv.Atoi(res.Data[:8])
		// blocktemplates = append(blocktemplates, fmt.Sprintf("vers %08x algo %s\nprev %064x\nmerk %064x\ntime %08x bits %08x nonc %08x\n\n", bt[:4], fork.GetAlgoName(int32(algo), 10), bt[4:36], bt[36:68], bt[68:72], bt[72:76], bt[76:80]))
		responses = append(responses, res)
		fmt.Println(res.Data)
	}
	// sort.Strings(blocktemplates)
	// for i := range blocktemplates {
	// 	fmt.Print(blocktemplates[i])
	// }
	for i := range responses {
		g := getwork.ToBlockHeader(responses[i].Data)
		fmt.Printf("%08x ", g.Version)
		p, _ := hex.DecodeString(g.PrevBlock.String())
		fmt.Printf("%064x ", p)
		m, _ := hex.DecodeString(g.MerkleRoot.String())
		fmt.Printf("%064x ", m)
		fmt.Printf("%08x ", g.Timestamp.Unix())
		fmt.Printf("%08x ", g.Bits)
		fmt.Printf("%08x ", g.Nonce)
		fmt.Println()
	}
}
