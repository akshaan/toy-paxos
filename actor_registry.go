package main

// Actor registry is a global store for references to all proposers and acceptors
// This can be modified to add new acceptors and proposers. In a productions setting,
// this would be backed by something like Zookeeper with a watcher that periodically
// updates the lists of proposers and acceptors.
type ActorRegistry struct {
	proposer  *Proposer
	acceptors  []*Acceptor
}

var GlobalActorRegistry *ActorRegistry

// NewActorRegistry creates a new registry with new proposers and acceptors
func NewActorRegistry(numAcceptors int) {
	acceptors := make([]*Acceptor, numAcceptors)

	for i := 0; i < numAcceptors; i++ {
		acceptors[i] = NewAcceptor(uint64(i))
	}

	GlobalActorRegistry =  &ActorRegistry{
		proposer: NewProposer(),
		acceptors: acceptors,
	}
}

// getProposer fetches a reference to the proposer from the registry
func (r *ActorRegistry) getProposer() *Proposer {
	return r.proposer
}

// getAcceptors fetches references to all acceptors from the registry
func (r *ActorRegistry) getAcceptors() []*Acceptor {
	return r.acceptors
}

// getNumAcceptors gets the number of acceptors from the registry
func (r *ActorRegistry) getNumAcceptors() int {
	return len(r.acceptors)
}
