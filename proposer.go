package main

import (
	"errors"
	"html"
	"log"
	rand "math/rand"
	"strconv"
	"time"
)

// Proposer represents an actor that proposes new values
// to acceptors and determines whether there is consensus
// on the new value among a majority of acceptors
type Proposer struct {
	id             uint64
	messageChannel chan Message
	valueToPropose string
}

// NewProposer creates a new proposer
func NewProposer() *Proposer {
	return &Proposer{
		id:             0,
		messageChannel: make(chan Message),
		valueToPropose: generateNewProposal(),
	}
}

// generateNewProposal generates a random Enoji to be used
// as a new proposal value
func generateNewProposal() string {
	low := 0x1F601
	high := 0x1F64F
	rand.Seed(time.Now().UnixNano())
	randomEmojiInt := rand.Intn(high-low+1) + low
	randomEmojiStr := html.UnescapeString("&#" + strconv.Itoa(randomEmojiInt) + ";")
	return randomEmojiStr

}

// sendMessage sends a message to this proposer
func (p *Proposer) sendMessage(message Message) {
	p.messageChannel <- message
}

// sendMessageToAcceptors sends a message to all acceptors
func (p *Proposer) sendMessageToAcceptors(message Message) {
	acceptors := GlobalActorRegistry.getAcceptors()
	for _, acceptor := range acceptors {
		acceptor.sendMessage(message)
	}
}

// didReceiveNonFailureFromMajority checks if the number of non-failure
// messages is greater than than 1/2 the number of acceptors.
func didReceiveNonFailureFromMajority(numNonFailures int) bool {
	numAcceptors := GlobalActorRegistry.getNumAcceptors()
	return numNonFailures > (numAcceptors / 2)
}

// fetchValueFromPromisesIfExists fetches a value from a promise message
// if it is present. Note that if the value is an empty string, it is
// considered to be absent. At the moment all proposals are Emoji
// but any string barring the empty string is considered to be a valid
// value.
func fetchValueFromPromisesIfExists(promiseMessages []Message) string {
	var value string = ""
	maxAcceptedId := promiseMessages[0].id
	for _, promise := range promiseMessages {
		if promise.value != "" && promise.acceptedId >= maxAcceptedId {
			value = promise.value
			maxAcceptedId = promise.acceptedId
		}
	}
	return value
}

// processPromiseOrFailMessages processes promise (or failure) messages
// from the acceptor. It returns a proposal to the acceptors in case
// of a promise and a terminate message otherwise
func (p *Proposer) processPromiseOrFailMessages() error {
	valueToPropose := p.valueToPropose
	promiseMessages := make([]Message, 0)
	for i := 0; i < GlobalActorRegistry.getNumAcceptors(); i++ {
		promiseOrFailMessage := <-p.messageChannel
		switch promiseOrFailMessage.messageType {
		case PROMISE:
			log.Printf("Proposer received promise %+v", promiseOrFailMessage)
			promiseMessages = append(promiseMessages, promiseOrFailMessage)
		case FAIL:
			log.Printf("Proposer received fail %+v", promiseOrFailMessage)
		default:
			log.Panicf("Expected PROMISE or FAIL message but got %+v", promiseOrFailMessage)
		}
	}

	if didReceiveNonFailureFromMajority(len(promiseMessages)) {
		acceptedValue := fetchValueFromPromisesIfExists(promiseMessages)
		if acceptedValue != "" {
			valueToPropose = acceptedValue
		}
		proposeMessage := Message{
			id:          p.id,
			value:       valueToPropose,
			messageType: PROPOSE,
		}
		p.sendMessageToAcceptors(proposeMessage)
		return nil
	} else {
		log.Printf("Proposer did not find quorum, Paxos failed :(")
		terminateMessage := Message{
			messageType: TERMINATE,
		}
		p.sendMessageToAcceptors(terminateMessage)
		return errors.New("Proposer did not find quorum, Paxos failed :(")
	}

}

// processAcceptMessages processes accept messages from the acceptors
// and determines whether consensus has been reached
func (p *Proposer) processAcceptMessages() {
	acceptMessages := make([]Message, 0)
	for i := 0; i < GlobalActorRegistry.getNumAcceptors(); i++ {
		acceptMessage := <-p.messageChannel
		switch acceptMessage.messageType {
		case ACCEPT:
			log.Printf("Proposer received accept %+v", acceptMessages)
			acceptMessages = append(acceptMessages, acceptMessage)
		case FAIL:
			log.Printf("Proposer received failure to accept %+v", acceptMessages)
		default:
			log.Panicf("Expected ACCEPT or FAIL message but got")
		}
	}

	if didReceiveNonFailureFromMajority(len(acceptMessages)) {
		// At this point, acceptMessages will have > 1 element because majority of acceptors accepted
		log.Printf("Consensus reached. Value %s accepted", acceptMessages[0].value)
	} else {
		log.Printf("Consensus could not be reached, Paxos failed :(")
	}

}

// runPaxos runs the Paxos algrithm on the proposer
func (p *Proposer) runPaxos() {
	p.id += 1
	log.Printf("Starting Paxos....")
	log.Printf("Generated new proposal: %s", p.valueToPropose)
	prepareMessage := Message{
		id:          p.id,
		messageType: PREPARE,
	}
	p.sendMessageToAcceptors(prepareMessage)
	err := p.processPromiseOrFailMessages()
	if err == nil {
		p.processAcceptMessages()
	}

	log.Println("")
	log.Println("--------------------------")
	log.Println("")
	exitChannel <- 0
}

// Start starts the proposer
func (p *Proposer) Start() {
	go p.runPaxos()
}
