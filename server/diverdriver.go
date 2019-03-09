package main

import (
	"encoding/json"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tylerw1369/iotago"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/tylerw1369/diverdriver/logs"
	"github.com/tylerw1369/diverdriver/server/ipc"
)

var config *viper.Viper

/*
PRECEDENCE (Higher number overrides the others):
1. default
2. key/value store
3. config
4. env
5. flag
6. explicit call to Set
*/
func loadConfig() *viper.Viper {
	// Setup Viper
	var config = viper.New()

	flag.StringP("pow.type", "t", "fpgago", "'giota-fpga', 'giota', 'giota-cl', 'giota-sse', 'giota-carm64', 'giota-c128', 'giota-c' or giota-go'")
	flag.IntP("pow.maxMinWeightMagnitude", "m", 14, "Maximum Min-Weight-Magnitude (Difficulty for PoW)")

	var logLevel = flag.StringP("log.level", "l", "INFO", "'DEBUG', 'INFO', 'NOTICE', 'WARNING', 'ERROR' or 'CRITICAL'")

	flag.StringP("server.diverDriverPath", "s", "/tmp/diverDriver.sock", "Unix socket path of diverDriver")

	config.BindPFlags(flag.CommandLine)

	var configPath = flag.StringP("config", "c", "diverDriver.config.json", "Config file path")
	flag.Parse()

	logs.SetLogLevel(*logLevel)

	// Bind environment vars
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvPrefix("FPGADIVER")
	config.SetEnvKeyReplacer(replacer)
	config.AutomaticEnv()

	// Load config
	if len(*configPath) > 0 {
		_, err := os.Stat(*configPath)
		if !flag.CommandLine.Changed("config") && os.IsNotExist(err) {
			// Standard config file not found => skip
			logs.Log.Info("Standard config file not found. Loading default settings.")
			return config
		}

		logs.Log.Infof("Loading config from: %s", *configPath)
		config.SetConfigFile(*configPath)
		err = config.ReadInConfig()
		if err != nil {
			logs.Log.Fatalf("Config could not be loaded from: %s, %v", *configPath, err.Error())
		}
	}

	return config
}

func init() {
	logs.Setup()
	config = loadConfig()
	logs.SetLogLevel(config.GetString("log.level"))

	cfg, _ := json.MarshalIndent(config.AllSettings(), "", "  ")
	logs.Log.Debugf("Following settings loaded: \n %+v", string(cfg))
}

func main() {
	flag.Parse() // Scan the arguments list

	var powFunc giota.PowFunc
	var powType string
	var powVersion string
	var err error

	switch strings.ToLower(config.GetString("pow.type")) {

	case "giota":
		powType, powFunc = giota.GetBestPoW()
		powVersion = ""

	case "giota-go":
		powFunc = giota.PowGo
		powType = "gIOTA-Go"

	case "giota-cl":
		powFunc, err = giota.GetPowFunc("PowCL")
		if err == nil {
			powType = "gIOTA-PowCL"
		} else {
			powType, powFunc = giota.GetBestPoW()
			logs.Log.Infof("POW type '%s' not available. Using '%s' instead", "PowCL", powType)
		}

	case "giota-sse":
		powFunc, err = giota.GetPowFunc("PowSSE")
		if err == nil {
			powType = "gIOTA-PowSSE"
		} else {
			powType, powFunc = giota.GetBestPoW()
			logs.Log.Infof("POW type '%s' not available. Using '%s' instead", "PowSSE", powType)
		}

	case "giota-carm64":
		powFunc, err = giota.GetPowFunc("PowCARM64")
		if err == nil {
			powType = "gIOTA-PowCARM64"
		} else {
			powType, powFunc = giota.GetBestPoW()
			logs.Log.Infof("POW type '%s' not available. Using '%s' instead", "PowCARM64", powType)
		}

	case "giota-c128":
		powFunc, err = giota.GetPowFunc("PowC128")
		if err == nil {
			powType = "gIOTA-PowC128"
		} else {
			powType, powFunc = giota.GetBestPoW()
			logs.Log.Infof("POW type '%s' not available. Using '%s' instead", "PowC128", powType)
		}

	case "giota-c":
		powFunc, err = giota.GetPowFunc("PowC")
		if err == nil {
			powType = "gIOTA-PowC"
		} else {
			powType, powFunc = giota.GetBestPoW()
			logs.Log.Infof("POW type '%s' not available. Using '%s' instead", "PowC", powType)
		}

	case "giota-fpga":
		powFunc, err = giota.GetPowFunc("PowFPGA")
		if err == nil {
			powType = "gIOTA-FPGA"
		} else {
			powType, powFunc = giota.GetBestPoW()
			logs.Log.Infof("POW type '%s' not available. Using '%s' instead", "PowFPGA", powType)
		}

	ipcserver.SetPowFunc(powFunc)

	// Servers should unlink the socket pathname prior to binding it.
	// https://troydhanson.github.io/network/Unix_domain_sockets.html
	syscall.Unlink(config.GetString("server.diverDriverPath"))

	logs.Log.Info("Starting diverDriver...")
	ln, err := net.Listen("unix", config.GetString("server.diverDriverPath"))
	if err != nil {
		logs.Log.Fatal("Listen error:", err)
	}

	exited := false
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		exited = true
		logs.Log.Infof("Caught signal %s: diverDriver shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	logs.Log.Info("diverDriver started. Waiting for connections...")
	logs.Log.Infof("Listening for connections on \"%v\"", config.GetString("server.diverDriverPath"))
	logs.Log.Infof("Using POW type: %v", powType)
	for !exited {
		fd, err := ln.Accept()
		if err != nil && !exited {
			logs.Log.Infof("Accept error: %v", err)
			continue
		} else {
			logs.Log.Debugf("New connection accepted from \"%v\"", fd.RemoteAddr)
		}

		go ipcserver.HandleClientConnection(fd, config, powType, powVersion)
	}
}
