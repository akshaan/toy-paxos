package main

type MessageType int

// This block replicated enum behavior so that each messageType
// is assigned a unique int value
const (

	// PREPARE is sent from proposer --> acceptor to propose
	// a new ID for the current paxos round
	PREPARE MessageType = iota

	// PROMISE is sent from acceptor --> proposer to accept 
	// a new ID from the proposer
	PROMISE MessageType = iota

	// PROPOSE is sent from proposer --> acceptor to propose
	// a new value for the current paxos round
	PROPOSE MessageType = iota

	// ACCEPT is sent from  acceptor -> proposer to accept
	// a new value in the current paxos round
	ACCEPT MessageType = iota

	// FAIL is sent from acceptor -> proposer to indicate
	// failure to accept a proposed ID value from proposer
	FAIL MessageType = iota

	// TERMINATE is sent from proposer -> acceptor to indicate
	// that a round of paxos has failed
	TERMINATE MessageType = iota
)


// Message represents the various messages exchanged 
// among actors in the Paxos context.
// In a real production scenario, these would likely be
// RPC request objects with an externally defined schema
// (e.g. Thrift, GRPC).
type Message struct {
	id uint64
	value string
	messageType MessageType
	acceptedId uint64
}
