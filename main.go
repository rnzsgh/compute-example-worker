package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/golang/glog"
	"github.com/rnzsgh/compute-example-worker/listen"
)

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
}

func main() {

	sigs := make(chan os.Signal, 1)
	running := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	stopListening := listen.ListenForJobs()

	go func() {
		<-sigs
		running <- false
	}()

	<-running

	stopListening()

	log.Flush()
}
