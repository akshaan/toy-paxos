package main

import (
	"log"
	rand "math/rand"
)

// Acceptor represents an acceptor in the Paxos context.
// Acceptors accept values proposed by the
// proposer or reject them (in the case of a failure,
// for instance)
type Acceptor struct {
	maxId            uint64
	messageChannel   chan Message
	proposalAccepted bool
	acceptedId       uint64
	acceptedValue    string
	address          uint64
}

// NewAcceptor creates a new acceptor
func NewAcceptor(address uint64) *Acceptor {
	return &Acceptor{
		maxId:            0,
		messageChannel:   make(chan Message),
		proposalAccepted: false,
		address:          address,
	}
}

// sendMessage sends a new message to the acceptor a
func (a *Acceptor) sendMessage(message Message) {
	a.messageChannel <- message
}

// processPrepareMessage processes a prepare message from the proposer
// and returns a promise or a failure to the proposer
func (a *Acceptor) processPrepareMessage(prepareMessage Message) {
	if prepareMessage.messageType != PREPARE {
		return
	}
	log.Printf("Acceptor %d received prepare message %+v", a.address, prepareMessage)
	proposer := GlobalActorRegistry.getProposer()
	if prepareMessage.id <= a.maxId || *acceptorFailureProb > rand.Float64() {
		failMessage := Message{
			messageType: FAIL,
		}
		proposer.sendMessage(failMessage)
	} else {
		a.maxId = prepareMessage.id
		if a.proposalAccepted {
			proposer.sendMessage(
				Message{
					id:          a.maxId,
					value:       a.acceptedValue,
					acceptedId:  a.acceptedId,
					messageType: PROMISE,
				},
			)
		} else {
			proposer.sendMessage(
				Message{
					id:          a.maxId,
					messageType: PROMISE,
				},
			)
		}
	}
}

// processProposeMessage processes a proposal from the proposer
// and returns an accept message or a failure
func (a *Acceptor) processProposeMessage(proposeMessage Message) {
	proposer := GlobalActorRegistry.getProposer()
	if proposeMessage.messageType == TERMINATE {
		return
	}
	if proposeMessage.id == a.maxId && *acceptorFailureProb < rand.Float64() {
		a.proposalAccepted = true
		a.acceptedId = proposeMessage.id
		a.acceptedValue = proposeMessage.value
		log.Printf(
			"Acceptor %d accepted value %s, id %d",
			a.address,
			a.acceptedValue,
			a.acceptedId,
		)
		// In a production setting the accept messages might be sent
		// to additional services calles Learners in the original Paxos
		// paper. These learners might perform an operation that requires
		// consensus (e.g. updating a database). In this implementation,
		// we omit the learner and simply return the acceptance to the
		// proposer
		proposer.sendMessage(
			Message{
				id:          a.maxId,
				messageType: ACCEPT,
				value:       proposeMessage.value,
			},
		)
	} else {
		log.Printf(
			"Acceptor %d failed to accept value %s, id %d",
			a.address,
			proposeMessage.value,
			proposeMessage.id,
		)
		proposer.sendMessage(
			Message{
				messageType: FAIL,
			},
		)
	}
}

// runPaxos runs the Paxos algorithm on acceptor a
func (a *Acceptor) runPaxos() {
	// Wait for prepare message
	prepareMessage := <-a.messageChannel

	// Respond to prepare message
	a.processPrepareMessage(prepareMessage)

	// Wait for propose message
	proposeMessage := <-a.messageChannel

	// Response to proposal
	a.processProposeMessage(proposeMessage)
}

// Start starts the acceptor
func (a *Acceptor) Start() {
	go a.runPaxos()
}
