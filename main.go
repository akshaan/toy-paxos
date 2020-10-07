package main

import "flag"

var numAcceptors = flag.Int("num-acceptors", 5, "Number of acceptors")
var acceptorFailureProb = flag.Float64("acceptor-failure-prob", 0.1, "Acceptor failure probability [0, 1)")

var exitChannel chan int = make(chan int)

func main() {
	// Parse CLI flags
	flag.Parse()

	for {
		// Initialize global ActorRegistry
		NewActorRegistry(*numAcceptors)

		// Run proposer
		GlobalActorRegistry.proposer.Start()

		// Run acceptor
		for _, acceptor := range GlobalActorRegistry.getAcceptors() {
			acceptor.Start()
		}

		_ = <-exitChannel
	}
}
