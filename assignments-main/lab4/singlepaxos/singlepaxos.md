# Single-decree Paxos

## Overview of the Skeleton Code

In this task you will implement the single-decree Paxos algorithm for each of the three Paxos roles.
This task will be verified by QuickFeed.

The skeleton code, unit tests, and definitions for this assignment can be found in the `singlepaxos` package.
Each of the three Paxos roles have separate files for:

- skeleton code (e.g., `acceptor.go`) and
- unit tests (e.g., `acceptor_test.go`).

There is additionally a single file called `defs.go`.
This file contains `struct` definitions for the four Paxos messages.
You should not edit this file.

Each of the three Paxos roles have a similar skeleton code structure.
They all have a constructor similar to:

```go
func NewAcceptor(id int) *Acceptor {
```

Each Paxos role also have a `handle` method for each message type they are expected to process.
For example, the `Acceptor` has a method for processing `Accept` messages with the following signature:

```go
func (a *Acceptor) handleAccept(accept Accept) Learn {
```

A `handle` method returns another Paxos message that should be sent as a result of handling the input message.
However, if the returned message is empty, the caller should not send any message, indicating that the Paxos algorithm ignored the input message.

For example, if an `Acceptor` handles an `Accept` message and should, according to the algorithm, reply with a `Learn` message, then the `handleAccept` should return the corresponding `Learn` message with its fields set to the correct values.
If handling the `Accept` resulted in no outgoing `Learn` message, then the empty `Learn{}` message should be returned.
In other words, the caller should _always_ check whether the returned message is empty, before deciding to send the returned message to other Paxos nodes.

> In Go parlance, the zero value of a struct is the value of the struct when it is initialized using an empty struct literal, e.g., `Learn{}`.

The `handleLearn` method from `learner.go` does not output a Paxos message.
The return `Value` instead represents the value that the Paxos nodes have reached consensus on (i.e., decided).
This value is meant to be used internally on a node to indicate that a value was chosen.

## Definitions

The `Value` type is defined in `defs.go` as follows:

```go
type Value string

const ZeroValue Value = ""
```

The `Value` definition represents the type of value the Paxos nodes should agree on.
For simplicity, we define the `Value` as an alias for the `string` type.
In later labs, it will be represented by something more application-specific, e.g., a client request.
A constant named `ZeroValue` is also defined to represent the empty value.

The Paxos message definitions are found in `defs.go`, and shown below.
We are using the naming convention found in [this](../resources/paxos-insanely-simple.pdf) algorithm specification (slide 64 and 65).

```go
type Prepare struct {
	From int
	Crnd Round
}

type Promise struct {
	To, From int
	Rnd      Round
	Vrnd     Round
	Vval     Value
}

type Accept struct {
	From int
	Rnd  Round
	Val  Value
}

type Learn struct {
	From int
	Rnd  Round
	Val  Value
}
```

Note that _only_ the `Promise` message struct has a `To` field.
This is because the `Promise` should only be sent to the `Proposer` who sent the corresponding `Prepare` (unicast).
The other three messages should all be sent to every other Paxos node (broadcast).

The `Round` type definition is also found in `defs.go`, and is a type alias for an `int`:

```go
type Round int

const NoRound Round = -1
```

There is also an important constant named `NoRound`.
This constant should be used in `Promise` messages, specifically for the `Vrnd` field, to indicate that an `Acceptor` has not voted in any previous round.

## Tasks

In this task, you should implement the following (all marked with `TODO(student)`):

- Any **unexported** field you may need in the `Proposer`, `Acceptor`, and `Learner` struct.

- The constructor for each of the three Paxos roles: `NewProposer`, `NewAcceptor`, and `NewLearner`.

- The `increaseCrnd` method in `proposer.go`.

  Every `Proposer` must maintain a set of unique round numbers to use when issuing proposals.
  For this assignment the `Proposer` is defined to have its `crnd` field initially set to the same value as its `id`.
  The `increaseCrnd` method should increase the current `crnd` by the total number of Paxos nodes.
  This is one way to ensure that every proposer uses a disjoint set of round numbers for proposals.

  > Note that you need to do a type conversion (`Round(id)`) in the constructor to assign the `id` to the `crnd` field.

- The `handlePromise` method in `proposer.go`.

  > _Note:_ The `Proposer` has a field named `clientValue` of type `Value`.
  > The `Proposer.handlePromise()` method should use the `clientValue` as the value to be chosen in an eventual outgoing `Accept` message, unless another value has been locked in by a quorum of `Promise` messages.

- The `handlePrepare` and `handleAccept` method in `acceptor.go`.

- The `handleLearn` method in `learner.go`.

> _Note:_ This task is solely the core message handling for each Paxos role.
> You may need to add fields to each Paxos role's struct to maintain Paxos state.

## Tests

Each of the three Paxos roles also have a separate `_test.go` file with unit tests.
You are free to add more test cases to these files or in separate test files; they will be ignored by QuickFeed.

Each `handle` method is tested separately.
The test cases contains a sequence of input messages, along with the expected output messages.

The test cases also provide a description of the actual invariant being tested.
You should take a look at the test code to get an understanding of what is going on.
An example of a failing acceptor test case is shown below:

```go
--- FAIL: TestSinglePaxosProposer (0.00s)
    --- FAIL: TestSinglePaxosProposer/IgnoreDifferentRound (0.00s)
        proposer_test.go:18: promise for different round (1) than our current one (2) -> ignore promise
        proposer_test.go:19: handlePromise() mismatch (-want +got):
              singlepaxos.Accept{
            - 	From: 0,
            + 	From: -1,
            - 	Rnd:  0,
            + 	Rnd:  -2,
            - 	Val:  "",
            + 	Val:  "FooBar",
              }
        proposer_test.go:27: Accept.Val = "FooBar", want ""
```
