# Lab 3: Failure Detector and Leader Election

| Lab 3: | Failure Detector and Leader Election |
| ---------------------    | --------------------- |
| Subject:                 | DAT520 Distributed Systems |
| Deadline:                | **February 22, 2024 23:59** |
| Expected effort:         | 20-25 hours |
| Grading:                 | Pass/fail |
| Submission:              | Group |

## Table of Contents

1. [Introduction](#introduction)
2. [Prerequisites](#prerequisites)
3. [Failure Detector (20%)](#failure-detector-(20%))
4. [Specification](#specification)
5. [Leader Detector (20%)](#leader-detector-(20%))
6. [Gorums-based Failure Detector (50%)](#gorums-based-failure-detector-(50%))
7. [Gorums-based Replicas on Localhost](#gorums-based-replicas-on-localhost)
8. [Gorums-based Replicas in Docker Containers (10%)](#gorums-based-replicas-in-docker-containers-(10%))

## Introduction

The main objective of this lab assignment is to implement a failure detector and a leader detector module.
A failure detector provides information about which processes have crashed and which are correct.
A leader detector can use this information to identify a process that has _not_ failed, which may act as a _leader_, and coordinates certain steps of a distributed algorithm.
You will later use the leader detector that you implement in this lab assignment for exactly this purpose.
Specifically, the leader detector module will be used to elect a single node as the Proposer of the Paxos algorithm.
The Paxos algorithm needs a leader detection to trigger its Phase 1 exchange.

This lab consists of three parts.
Each part will be explained in more detail in their own sections.

1. **Failure Detector module:**
   Implement the eventually perfect failure detector algorithm from the textbook.
   Use the provided skeleton code and unit tests.
   This implementation will be verified by QuickFeed.
   The task represents 20% of the total lab.

2. **Leader Detector module:**
   Implement the eventual leader detector algorithm from the textbook.
   Use the provided skeleton code and unit tests.
   This implementation will also be verified by QuickFeed.
   The task represents 20% of the total lab.

3. **Gorums-based failure detector:**
   Implement a failure detector that uses Gorums to communicate between processes.
   The failure detector will be based on the eventually perfect failure detector implemented in part 1, but will use a multicast call to send heartbeat messages to the processes.
   This failure detector will be used when implementing Paxos in later labs.
   The task represents 60% of the total lab.

## Prerequisites

You need to register your group on [QuickFeed](https://uis.itest.run/app/home) before you begin this assignment as it constitutes a group project.
This can be done by creating a [new group](https://github.com/dat520-2024/info/blob/main/signup.md#group-signup-on-quickfeed) on the course's page on QuickFeed.
Select the students you are collaborating with and submit the group selection.
**Only one group member should do this.**
Please don't create a group unless you have agreed with the other member(s) up front.

If you don't have a group partner yet, you may take a look at the `group-maker` channel on discord.
If you see a person listed there that you wish to work with, please connect with him/her directly and agree to submit a group composition accordingly, following the above instructions.

**Important 1:**
One group consist of two or three students.
We only allow at most three students to collaborate on the group project, but only if there is a valid reason for this.
Four members will not be allowed.

**Important 2:**
Your group will only be approved when all members have passed both lab assignment 1 and 2.

All group members will get access to a shared group repository when your group has been approved.
The name of the repository will be the one selected when you create the group in QuickFeed.
You will receive an email notification when QuickFeed creates a new team on GitHub.
Refer to the procedure described [here](https://github.com/dat520-2024/info/blob/main/lab-submission.md#working-with-group-assignments) for instructions on how to setup and work with the group repository on your local machine.

## Failure Detector (20%)

A failure detector can be implemented using two approaches:

1. request/reply approach;
2. lease-based approach;

A failure detector that follows the request/reply approach is shown in Algorithm 2.7 in the textbook.
It sends a `HeartbeatRequest` to all other nodes, and if the request is not answered with a `HeartbeatReply` within a certain time, it suspects the silent process.

A failure detector using the lease-based approach is divided into a sending and receiving process.
The receiving process is essentially the same as the failure detector described in Algorithm 2.7, but with the crucial difference that it does not send `HeartbeatRequest` messages.
Instead, the sending process is a simple loop, that upon timeout sends a `HeartbeatReply` to all other processes.
Thus, this failure detector uses two timeouts, one for receiving and one for sending heartbeat messages.
You will implement such a lease-based failure detector in the Gorums-based failure detector in Part 3 of this lab.

### Specification

In this task you will implement an Eventually Perfect Failure Detector.
The specification for this failure detector is described on pages 53-56 in the textbook.
Your failure detector should use the Increasing Timeout algorithm.
See Algorithm 2.7 for more details.

You should use the provided skeleton to implement the failure detector.
All skeleton code and corresponding tests for this assignment can be found in the [failuredetector](./failuredetector) package.
The skeleton code is located in the `failuredetector.go` file, and is listed below.
Large parts of the failure detector is already implemented, but you will need to complete important remaining parts.
You should complete the implementation by extending the parts of the code marked with the `TODO(student)` label.
The failure detector specification is documented using code comments.
You should refer to these comments for a detailed specification of the implementation.

The unit tests for the failure detector is located in the file `failuredetector_test.go`.
You can run all the tests in the detector package using the command `go test -v`.
As described in previous labs, you can also use the `-run` flag to only run a specific test.
You are also encouraged to take a close look at the test code to see what is actually being tested.
This may help you when writing and debugging your code.

The initial skeleton code for the failure detector in `failuredetector.go` is listed below:

```go
package failuredetector

import "time"

// EvtFailureDetector represents a Eventually Perfect Failure Detector as
// described at page 53 in:
// Christian Cachin, Rachid Guerraoui, and Luís Rodrigues: "Introduction to
// Reliable and Secure Distributed Programming" Springer, 2nd edition, 2011.
type EvtFailureDetector struct {
	id        int          // the id of this node
	nodeIDs   []int        // node ids for every node in cluster
	alive     map[int]bool // map of node ids considered alive
	suspected map[int]bool // map of node ids considered suspected

	sr SuspectRestorer // Provided SuspectRestorer implementation

	delay         time.Duration // the current delay for the timeout procedure
	delta         time.Duration // the delta value to be used when increasing delay
	timeoutSignal *time.Ticker  // the timeout procedure ticker

	hbSend chan Heartbeat // channel for sending outgoing heartbeat messages
	hbRecv chan Heartbeat // channel for receiving incoming heartbeat messages
	stop   chan struct{}  // channel for signaling a stop request to the main run loop

	testHook func() // DO NOT REMOVE THIS LINE. A no-op when not testing.
}

// NewEvtFailureDetector returns a new Eventual Failure Detector. It takes the
// following arguments:
//
// id: The id of the node running this instance of the failure detector.
//
// nodeIDs: A list of ids for every node in the cluster (including the node
// running this instance of the failure detector).
//
// sr: A leader detector implementing the SuspectRestorer interface.
//
// delta: The initial value for the timeout interval. Also the value to be used
// when increasing delay.
func NewEvtFailureDetector(id int, nodeIDs []int, sr SuspectRestorer, delta time.Duration) *EvtFailureDetector {
	suspected := make(map[int]bool)
	alive := make(map[int]bool)

	// TODO(student): perform any initialization necessary

	return &EvtFailureDetector{
		id:        id,
		nodeIDs:   nodeIDs,
		alive:     alive,
		suspected: suspected,

		sr: sr,

		delay: delta,
		delta: delta,

		hbSend: make(chan Heartbeat, 16),
		hbRecv: make(chan Heartbeat, 8),
		stop:   make(chan struct{}),

		testHook: func() {}, // DO NOT REMOVE THIS LINE. A no-op when not testing.
	}
}

// Start starts e's main run loop as a separate goroutine. The main run loop
// handles incoming heartbeat requests and responses. The loop also trigger e's
// timeout procedure at an interval corresponding to e's internal delay
// duration variable.
func (e *EvtFailureDetector) Start() {
	e.timeoutSignal = time.NewTicker(e.delay)
	go func(timeout <-chan time.Time) {
		for {
			e.testHook() // DO NOT REMOVE THIS LINE. A no-op when not testing.
			select {
			case <-e.hbRecv:
				// TODO(student): Handle incoming heartbeat
			case <-timeout:
				e.timeout()
			case <-e.stop:
				return
			}
		}
	}(e.timeoutSignal.C)
}

// Stop stops the failure detector's main run loop.
//
// Calling Stop multiple times without a Start call in between them will panic.
func (e *EvtFailureDetector) Stop() {
	e.stop <- struct{}{}
}

// Deliver delivers heartbeat to the failure detector.
//
// After stopping the failure detector, Deliver will panic if called.
func (e *EvtFailureDetector) Deliver(heartbeat Heartbeat) {
	e.hbRecv <- heartbeat
}

// Heartbeats returns a receive-only channel on which the failure detector sends
// outgoing heartbeats.
//
//   - Heartbeat replies are sent in response to incoming heartbeat requests.
//   - Heartbeat requests are sent periodically to all other nodes in the cluster
//     as part of the timeout procedure.
//
// The channel is closed when the failure detector is stopped, after which no
// more heartbeats will be sent and trying to read from the channel will panic.
func (e *EvtFailureDetector) Heartbeats() <-chan Heartbeat {
	return e.hbSend
}

// timeout updates the failure detector's internal state and sends out
// heartbeat requests to all other nodes in the cluster. If the intersection
// between the set of alive nodes and the set of suspected nodes is not empty,
// the delay is increased by delta and the timeout procedure is restarted.
//
// After stopping the failure detector, timeout will panic if called.
func (e *EvtFailureDetector) timeout() {
	// TODO(student): Implement timeout procedure
}

// TODO(student): Add other unexported functions or methods if needed.

```

The failure detector uses a `Heartbeat` struct to send heartbeat request and replies.
Outgoing heartbeats from the failure detector should be sent on the `hbSend` channel.
The `Heartbeat` struct is defined in the `defs.go` file:

```go
// A Heartbeat is the basic message used by failure detectors to communicate
// with other nodes. A heartbeat message can be of type request or reply.
type Heartbeat struct {
	From int
	To   int
	Type hbType // true -> hbRequest, false -> hbReply
}
```

You should complete the following parts of the code in `failuredetector.go`:

- Perform any initialization necessary in the `NewEvtFailureDetector(...)` function.

- Implement handling of incoming heartbeat messages.
  More specifically:
  `case <-e.hbRecv:` in the `select` statement within the `Start()` method.
  Your code should check if the heartbeat is a request or response, and act accordingly.

- Implement the failure detector's timeout procedure:
  `func (e *EvtFailureDetector) timeout()`.
  During the timeout procedure the failure detector should inform the provided `SuspectRestorer` object if it thinks a node is suspected of being faulty or seem to have been restored.
  The failure detector has a reference to a `SuspectRestorer` entity via the `sr` field in the `EvtFailureDetector` struct.
  This field is of the type `SuspectRestorer` which is an interface.
  It defines two methods, `Suspect(id)` and `Restore(id)`, which is available to the failure detector.
  If the failure detector determines that a node should be considered suspected or restored during the timeout procedure, then it should inform the `SuspectRestorer` object using these two available methods.
  You will in the next task implement a leader detector that satisfies this interface.
  Take a look [here](https://golang.org/doc/effective_go.html#interfaces_and_types) to learn more about how interfaces work in Go.

Other information:

- A node identifier is defined to be an integer (type `int`).

- A node should not treat itself as special during the timeout procedure.
  This means that a node should not set itself alive without receiving a heartbeat reply from itself.
  This is done to simplify testing, keep the code general and to let it mirror the algorithm in the textbook as closely as possible.

- You may add your own _unexported_ functions or methods if needed.

## Leader Detector (20%)

In this task you will implement the Monarchical Eventual Leader Detector.
The description of this type of leader detector can be found on pages 56-60 in the textbook.
See Algorithm 2.8 for a complete description.

All code and corresponding tests for the leader detector module can be found in the [leaderdetector](./leaderdetector) package.
The skeleton code for the leader detector is located in the `leaderdetector.go` file, and is listed below.
As you did for the failure detector, you need to implement the relevant parts of the code marked with `TODO(student)`.
Refer to the code comments for a complete specification of each type, function, and method.
The corresponding tests can be found in `leaderdetector_test.go`.

```go
package leaderdetector

// A MonLeaderDetector represents a Monarchical Eventual Leader Detector as
// described at page 53 in:
// Christian Cachin, Rachid Guerraoui, and Luís Rodrigues: "Introduction to
// Reliable and Secure Distributed Programming" Springer, 2nd edition, 2011.
type MonLeaderDetector struct {
	// TODO(student): Add needed fields
}

// NewMonLeaderDetector returns a new Monarchical Eventual Leader Detector
// given a list of node ids.
func NewMonLeaderDetector(nodeIDs []int) *MonLeaderDetector {
	// TODO(student): Add needed implementation
	m := &MonLeaderDetector{}
	return m
}

// NodeIDs returns the list of node ids.
func (m *MonLeaderDetector) NodeIDs() []uint32 {
	// TODO(student): Implement
	return []uint32{}
}

// Leader returns the current leader. Leader will return UnknownID if all nodes
// are suspected.
func (m *MonLeaderDetector) Leader() int {
	// TODO(student): Implement
	return UnknownID
}

// Suspect instructs the leader detector to consider the node with matching
// id as suspected. If the suspect indication result in a leader change
// the leader detector should publish this change to its subscribers.
func (m *MonLeaderDetector) Suspect(id int) {
	// TODO(student): Implement
}

// Restore instructs the leader detector to consider the node with matching
// id as restored. If the restore indication result in a leader change
// the leader detector should publish this change to its subscribers.
func (m *MonLeaderDetector) Restore(id int) {
	// TODO(student): Implement
}

// Subscribe returns a buffered channel which will be used by the leader
// detector to publish the id of the highest ranking node.
// The leader detector will publish UnknownID if all nodes become suspected.
// Subscribe will drop publications to slow subscribers.
// Note: Subscribe returns a unique channel to every subscriber;
// it is not meant to be shared.
func (m *MonLeaderDetector) Subscribe() <-chan int {
	// TODO(student): Implement
	return nil
}

// TODO(student): Add other unexported functions or methods if needed.
```

Other information:

- The _Monarchical_ Eventual Leader Detection uses a node ranking.
  The ranking for this implementation is defined using the node identifiers.
  The highest ranking is defined to be the node with the highest id (i.e. the highest integer).

- Negative node IDs should be ignored by the leader detector.
  An exception is the special identifier constant `UnknownID` in `defs.go` with value `-1`.

- The leader detector algorithm from the textbook send a `< TRUST | leader >` indication event when a new leader is detected.
  This behavior is modeled in the code by the `Subscribe()` method.
  Any part of a application can call `Subscribe()` to be notified about leader change events.
  Each caller will receive a unique channel where they may receive `< TRUST >` indications in the form of node IDs (`int`).

- The `NodeIDs` method returns the slice of node IDs as `[]uint32`.
  This allows us to get the node IDs from the leader detector for use with Gorums; since Gorums uses `uint32` for its IDs.

- You may add your own _unexported_ functions or methods if needed.

## Gorums-based Failure Detector (50%)

In this task you will implement a lease-based failure detector that uses Gorums to communiate; recall the description from [Part 1](#failure-detector-25) above.
The skeleton for this version can be found in the [`gorumsfd`](./gorumsfd/) folder.
The protobuf definition used in this task can be found in the [`proto`](./gorumsfd/proto) folder; they have been precompiled.

All nodes send heartbeat to each other using a multicast method at some interval.
A node receiving a heartbeat can update its local state since it now knows that the sender of the heartbeat is operational.

Periodically, each node runs the `timeout()` method to identify suspected and restored nodes.
The `timeout()` method is responsible for calling the leader detector's `Suspect` method for each suspected node.
Similarly, the `timeout()` method must also call the `Restore()` method for each node that gets a restored status;
that is, a node that received a heartbeat from a previously suspected node.

You need to implement the following methods:

- `Start(hbSender func(*pb.HeartBeat))`:
  This is the failure detector's run loop; it should start at least one goroutine, to perform the following functionalities:

  1. Periodically send heartbeats to all nodes in the configuration;
     the provided `hbSender` function should be used for this purpose.
  2. Periodically update the status of the nodes to the `SuspectRestorer`.
  3. Wait to receive the signal to stop the goroutines.

  The `Start` method may perform the periodic tasks in a single goroutine or two goroutines.
  However, care must be taken to configure the periods between the sending of heartbeats and reacting to missing heartbeats.

- `timeout()`:
  This method is described above and is responsible for the main logic of the failure detector, and for updating the leader detector using the `Suspect` and `Restore` methods.
  The logic to update the `delay`, the `alive` and `suspected` sets is similar to the failure detector in [Part 1](#failure-detector-25).

The [`failuredetector_test.go`](./gorumsfd/failuredetector_test.go) file contains tests that can be used to verify your implementation.

### Gorums-based Replicas on Localhost

We provide an implementation of a Gorums replica that uses your failure detector.
You can use this to validate your failure detector's code.
Please do inspect the code in the `cmd/replica` folder.

To test your implementation, you can build it with the provided `Makefile` and the provided `replicas.sh` script.

```console
cd lab3
make
./replicas.sh
```

This script creates three replicas and after 3 seconds kills one of them to verify the failure detector and leader detector responds accordinly.

> [!WARNING]
> **Please disable logging of heartbeats.**
>
> During testing it can be useful to enable logging of receive heartbeats.
> However, this may generate a lot of noise when running tests on QuickFeed.
> Therefore, please comment out logging before pushing your code to GitHub.
> Optionally, you can push your intermediate changes to a separate branch, since QuickFeed only runs tests against the default branch (typically, main or master).

### Gorums-based Replicas in Docker Containers (10%)

In this part, you will prepare a `Dockerfile` that contains the `replica` binary.
You should be able to run multiple docker containers that can connect via a network with each other, such that they each can detect the failure of another replica.

Please add the `Dockerfile` to the `lab3` folder.
Instructions explaining how to run your containers should be added in a file `docker.md` also added to the `lab3` folder.

This assignment will be assessed by the TAs during approval.
