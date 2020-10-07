# Toy Paxos
A toy implementation of the [Paxos](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf) protocol in Golang.

## Getting started
- Install Go according to the instructions [here](https://www.google.com/search?q=install+golang&rlz=1C5CHFA_enUS847US847&oq=install+golang&aqs=chrome..69i57j0l7.3102j0j7&sourceid=chrome&ie=UTF-8)
- Run `go build` to build the compiled binary. The binary will be written to `./paxos`.
- Run the compiled binary using 
```
./paxos --num-acceptors <number-of-acceptors> --acceptor-failure-prob <probability-of-failure>
```
- Running the binary will keep running new Paxos rounds with new proposals. Use Ctrl-C to terminate.

## Files
```
.
├── `README.md` : this file
├── `acceptor.go` : Implementation of an Acceptor service
├── `actor.go` : Interface that both Acceptors and Proposers conform to
├── `actor_registry.go` : A registry that lists the Proposers and all Acceptors in the current Paxos run
├── `message.go` : Message struct that is exchanged between the Proposer and Acceptors
├── `main.go` : The top level file that runs rounds of Paxos in an infinite loop
└── `proposer.go` : Implementation of an Acceptor service
```

## Assumptions
- We assume there is only a single proposer at the moment
- The proposer does not fail
- We assume that message IDs from the proposer start at 1
- We assume that an empty string is not a valid value in each message. Currently all generated proposals are Emojis but any string 
  barring the empty string is a valid value 😬

## Potential features/improvements
- Tests, tests, tests
- Extend to multiple proposers
- Simulate proposer failures
- Simulate real failures more closely, using absent responses and timeouts instead of explicit FAIL messages
- Allow any string value to be a value proposal
- Replace global registry of proposer and acceptors with local Zookeeper instance
- Replace Message struct in message.go with true RPC framework (e.g. Thrift of GRPC)
