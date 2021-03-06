package main

import (
	"flag"
	log "github.com/thinkboy/log4go"
	"goim/libs/perf"
	"runtime"
	"yf-im/yfgoim"
)

func main() {
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err)
	}
	runtime.GOMAXPROCS(Conf.MaxProc)
	log.LoadConfiguration(Conf.Log)
	defer log.Close()
	log.Info("logic[%s] start", Ver)
	perf.Init(Conf.PprofAddrs)
	// router rpc
	if err := InitRouter(); err != nil {
		log.Warn("router rpc current can't connect, retry")
	}
	MergeCount()
	go SyncCount()
	// logic rpc
	if err := InitRPC(yfgoim.NewYfAuther()); err != nil {
		panic(err)
	}
	if err := InitHTTP(); err != nil {
		panic(err)
	}
	if err := InitKafka(Conf.KafkaAddrs); err != nil {
		panic(err)
	}
	// block until a signal is received.
	InitSignal()
}
