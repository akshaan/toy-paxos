package main

// Actor is an interface for an actor in the distributed context
// where Paxos is running. Proposers and Acceptors conform to this
// interface
type Actor interface {
	Start()                      // Start the actor
	SendMessage(message Message) // Send a message to this actor
}
