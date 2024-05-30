# Multi-Paxos

## Background and Resources

Practical systems normally use Paxos as a building block to achieve consensus on a sequence of values, e.g., commands to be executed.
One way to achieve this would be to run a full instance of single-decree Paxos, both _Phase 1 and 2_, for every value.
This would require four message delays for every value to be decided.
With Multi-Paxos it is possible to reduce this overhead.

Multi-Paxos only perform _Phase 1_ once, when the leader changes.
A Proposer, thinking it is the new leader, issues a `Prepare` for every slot higher than the highest consecutively decided slot that it has seen.
In response to such a `Prepare`, each Acceptor sends a `Promise` message back to the Proposer, only if the `Prepare` message's round is higher than the Acceptor's current round.
The `Promise` message may contain a set of `(vrnd, vval)` tuples for every slot higher than or equal to the round from the `Prepare` message, if the Acceptor has already accepted any value for these slots.
Upon receiving a quorum of `Promise` messages, the Proposer is bound by (or locked in to) the highest `(vrnd, vval)` tuple reported for any slot higher than the slot from the corresponding `Prepare` message.

The Proposer can then perform _Phase 2_ (`Accept` and `Learn`) for every value to be decided.
Only two message delays are required to get a value accepted.

You are _strongly_ advised to read Section 3, _Implementing a State Machine_, from _Paxos Made Simple_ by Leslie Lamport for a more complete description of the Multi-Paxos optimization.
You may also find the other [resources](../resources/) listed here useful.

## Overview of the Skeleton Code

In this task you will implement the Multi-Paxos algorithm for each of the three Paxos roles.
The task is similar to what you did for single-decree Paxos, but is more complex since Multi-Paxos is able to choose multiple commands.
Both _Phase 1_ and _Phase 2_ of the Paxos protocol, as described in _Paxos Made Simple_, needs to be adjusted.
Especially the `Prepare`-`Promise` exchange needs to be modified.

The skeleton code, unit tests, and definitions for this assignment can be found in the `multipaxos` package.
Each of the three Paxos roles has a separate file for:

- skeleton code (e.g., `acceptor.go`) and
- unit tests (e.g., `acceptor_test.go`).

There is also a file called `defs.go` that contains `struct` definitions for the four Paxos messages and other related definitions.
You should not edit this file.

Similar to the single-decree Paxos, each of the three Paxos roles has a similar skeleton code structure.
They all have a constructor similar to:

```go
func NewAcceptor(id int) *Acceptor {
```

Each Paxos role also have a `handle` method for each message type they are expected to process.
For example, the Acceptor has a method for processing `Accept` messages with the following signature:

```go
func (a *Acceptor) handleAccept(accept Accept) Learn {
```

A `handle` method returns another Paxos message that should be sent as a result of handling the input message.
However, if the returned message is empty, the caller should not send any message, indicating that the Paxos algorithm ignored the input message.

For example, if an Acceptor handles an `Accept` message and should, according to the algorithm, reply with a `Learn` message, then the `handleAccept` should return the corresponding `Learn` message with its fields set to the correct values.
If handling the `Accept` resulted in no outgoing `Learn` message, then the empty `Learn{}` message should be returned.
In other words, the caller should _always_ check whether the returned message is empty, before deciding to send the returned message to other Paxos nodes.

The `handleLearn` method from `learner.go` has the following signature:

```go
func (l *Learner) handleLearn(learn Learn) (Value, Slot) {
```

This method does not output a Paxos message.
The return `Value` instead represents a value for a specific slot (`Slot`) that the Paxos nodes has reached consensus on (i.e. decided).
One a difference, however, from the single-decree Paxos Acceptor is the returned `Slot`.
The `Slot` indicates which slot the corresponding decided `Value` belongs to.

## Definitions

The `Value` type definition in `defs.go` has also changed from being a type alias for `string` to the following struct definition for this task:

```go
type Value struct {
	ClientID  string
	ClientSeq int
	Command   string
}
```

The `Value` type now carries information about the client that sent the command.

- `ClientID` is a unique client identifier.
- `ClientSeq` represents a client sequence number.
  This is used by clients to match a response, received from a Paxos system, to a corresponding request made by the client.
- `Command` is the state machine command to be executed by the Paxos system.

The `Round` type definition found in `defs.go` remain unchanged:

```go
type Round int

const NoRound Round = -1
```

However, the Paxos messages have changed slightly to the following:

```go
type Prepare struct {
	From int
	Slot Slot
	Crnd Round
}

type Promise struct {
	To, From int
	Rnd      Round
	Accepted []PValue
}

type Accept struct {
	From int
	Slot Slot
	Rnd  Round
	Val  Value
}

type Learn struct {
	From int
	Slot Slot
	Rnd  Round
	Val  Value
}
```

The `Prepare`, `Accept`, and `Learn` messages have all gotten a `Slot` field of type `Slot`.
This means that every `Accept` and `Learn` message now relates to a specific slot.

However, the `Slot` field in the `Prepare` message has a somewhat different meaning.
In Multi-Paxos, as explained [previously](#background-and-resources), a proposer only executes _Phase 1_ once on every leader change if it considers itself to be the leader.
Thus, the `Slot` field in the `Prepare` message represents the slot after the highest consecutive decided slot that the Proposer has seen.

This slot identifier is used by an Acceptor to construct a corresponding `Promise` as a reply.
That is, an Acceptor attaches information (`vrnd` and `vval`) for every slot that it has sent an `Accept` message for, that were equal to or higher than the one received in the `Prepare` message.
This information is included in the `Accepted` field of type `[]PValue`.
The slice should be sorted by increasing `Slot`.
The `PValue` struct is also defined in `defs.go`:

```go
type PValue struct {
	Slot Slot
	Vrnd Round
	Vval Value
}
```

To create and append the correct slots (if any) to the `[]PValue` slice, an Acceptor must keep track of the highest seen slot, for which it has sent an `Accept` message.
This can be done by maintaining a `highestSeen` variable of type `Slot`.

Once the Proposer receives a quorum of promises, it becomes locked to the value in `PValue.Vval` with the highest `Vrnd`, for each slot higher than the `Prepare` message's `Slot` field.

## Specification: Proposer

Below we focus on the Proposer role and its `handlePromise` method.
It has the following signature:

```go
func (p *Proposer) handlePromise(promise Promise) []Accept {
```

- _Input:_ A single `Promise` message.

  - The Proposer should ignore a `Promise` message if its round is different from the Proposer's current round.
    This means that it is not an answer to a previously sent `Prepare` from this Proposer.

  - The Proposer should ignore a `Promise` message if it has previously received a promise from the _same_ node for the _same_ round.

- _Output:_

  - If the input `Promise` message should be ignored, then the `[]Accept` slice should be `nil`.

  - If the input `Promise` result in a quorum for the current round, then the `[]Accept` slice should contain accept messages for the slots for which the Proposer is locked in.

  - If the Proposer is not locked in on any slots an empty `[]Accept{}` slice should be returned.

- _Additional Requirements:_

  - All `Accept` messages in the `[]Accept` slice must be in increasing consecutive slot order.

  - If there is a gap in the set of slots for which the Proposer is locked in, then the Proposer should create an `Accept` message with an empty `Value{}` to represent a no-op value for the gap slots.
    For example, if the proposer is locked in on slots 2 and 4, but not for slot 3, then the returned `[]Accept` slice should look like:

    ```go
		[]Accept{
			{Slot: 2, Val: valueOne},
			{Slot: 3, Val: Value{}},
			{Slot: 4, Val: valueTwo},
		},
    ```

  - If a `PValue` in a `Promise` message is for a slot lower than the Proposer's current `adu` (all-decided-up-to), then the `PValue` should be ignored.

  - A `Promise` message does not indicate whether or not a slot has been decided.
    That is, if the Proposer receives a `Promise` with `[]PValue.Slot > Prepare.Slot` it does not necessarily mean that the slot has been decided.
    The new leader may not have previously learnt it.
    In this case the slot will be proposed and decided again.
    The Paxos protocol ensures that the same value will be decided for the same slot.

### Important Notes

A few other important aspects of the Paxos roles are listed below:

- An Acceptor only need to maintain a single `rnd` variable (as for single-decree Paxos).
  The `rnd` variable spans across all slots.

- Only `vrnd` and `vval` must be stored for each specific slot.

- Similarly, the Proposer only need to maintain a single `crnd` variable.

- The Paxos roles share no slot history/storage in this implementation.
  Each role should maintain their own variables and data structures for keeping track of promises, accepts, and learns for each slot.

## Tasks

Summarized, you should for this task implement the following (all marked with `TODO(student)`):

- Any **unexported** field you may need in the `Proposer`, `Acceptor`, and `Learner` struct.

- The constructor for each of the Paxos roles: `NewAcceptor` and `NewLearner`.
  The `NewProposer` constructor is already implemented.
  Note that the `Proposer` also take its `adu` as an argument for testing purposes.

- The `handlePrepare` and `handleAccept` method in `acceptor.go`.

- The `handleLearn` method in `learner.go`.

- The `handlePromise` method in `proposer.go`.

> _Note:_ This task is solely the core message handling for each Multi-Paxos role.
> You may need to add fields to each Multi-Paxos role struct to maintain the state.

## Tests

Each of the three Paxos roles also have a separate `_test.go` file with unit tests.
You are free to add more test cases to these files or in separate test files; they will be ignored by QuickFeed.

Each `handle` method is tested separately.
The test cases contains a sequence of input messages, along with the expected output messages.

You can find a detailed description of some of the proposer test cases [here](#appendix---proposer-test-cases).

The test cases also provide a description of the actual invariant being tested.
You should take a look at the test code to get an understanding of what is going on.
An example of a failing proposer test case is shown below:

```go
--- FAIL: TestMultiPaxosProposer (0.00s)
    --- FAIL: TestMultiPaxosProposer/PValueScenario5 (0.00s)
        proposer_test.go:16: scenario 5 - message 1 - see figure in README.md
        proposer_test.go:17: handlePromise() mismatch (-want +got):
              []multipaxos.Accept(
            - 	nil,
            + 	{s"Accept{From: -1, Slot: -1, Rnd: -2, Val: No-op value}"},
              )
        proposer_test.go:16: scenario 5 - message 2 - see figure in README.md
        proposer_test.go:17: handlePromise() mismatch (-want +got):
              []multipaxos.Accept(
            - 	nil,
            + 	{s"Accept{From: -1, Slot: -1, Rnd: -2, Val: No-op value}"},
              )
        proposer_test.go:16: scenario 5 - message 3 - see figure in README.md
        proposer_test.go:17: handlePromise() mismatch (-want +got):
              []multipaxos.Accept{
            - 	s"Accept{From: 2, Slot: 2, Rnd: 2, Val: Value{ClientID: 1234, ClientSeq: 42, Command: ls}}",
              	{
            - 		From: 2,
            + 		From: -1,
            - 		Slot: 3,
            + 		Slot: -1,
            - 		Rnd:  2,
            + 		Rnd:  -2,
              		Val:  {},
              	},
            - 	s"Accept{From: 2, Slot: 4, Rnd: 2, Val: Value{ClientID: 5678, ClientSeq: 99, Command: rm}}",
            - 	s"Accept{From: 2, Slot: 5, Rnd: 2, Val: Value{ClientID: 5678, ClientSeq: 99, Command: rm}}",
              }
```

## Lab Approval

For this lab you should present your code and explain what you implemented, comparing the multi-paxos implementation with the single-decree Paxos.
You should demonstrate that your implementation fulfills the previously listed specification and that you understood the particularities of each protocol.

## Appendix - Proposer Test Cases

The following gives an explanation of some of the proposer test cases.
For all examples, we initialize the proposer with the following values:

| Field                                       | Value |
| :------------------------------------------ | ----: |
| Proposer ID                                 |     2 |
| Number of nodes                             |     3 |
| Highest seen consecutive decided slot (adu) |     1 |

Thus, in the following scenarios, the Proposer has sent a `Prepare` message to the acceptors:

```go
Prepare{From: 2, Slot: 2, Crnd: 2}
```

### Test Scenario 1

In this scenario, the Proposer first receives a `Promise` message from Acceptor 1:

```go
Promise{To: 2, From: 1, Rnd: 2}
```

In response to this `Promise` message, the Proposer should not send any `Accept` messages since it has not received a quorum of `Promise` messages.
Hence, the expected output is `nil`.

Next, the Proposer receives a `Promise` message from Acceptor 0:

```go
Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
  {Slot: 2, Vrnd: 1, Vval: valOne},
}}
```

In response to this `Promise` message, the Proposer should now send an `Accept` since it has received a quorum of `Promise` messages.
Since Acceptor 0 has already voted for a value (`valOne`) in slot 2, which is higher than the `Prepare` message's slot, the Proposer should adopt this value.
This is because another acceptor may also have voted for this value in slot 2, and thus it is safe to adopt this value.
Hence, the expected output should be:

```go
[]Accept{{From: 2, Slot: 2, Rnd: 2, Val: valOne}}
```

This scenario is covered by the following test case in `proposer_test.go`:

```go
name:     "PValueScenario1",
proposer: NewProposer(2, 3, 1, &mockLD{}),
msgs: []promiseAccepts{
	{
		promise: Promise{To: 2, From: 1, Rnd: 2},
    wantAccepts: nil,
	},
	{
		promise: Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
			{Slot: 2, Vrnd: 1, Vval: valOne},
		}},
		wantAccepts: []Accept{{From: 2, Slot: 2, Rnd: 2, Val: valOne}},
	},
},
```

### Test Scenario 2

In this scenario, the Proposer first receives a `Promise` message from Acceptor 1:

```go
Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
  {Slot: 2, Vrnd: 0, Vval: valOne},
}}
```

In response to this `Promise` message, the Proposer should not send any `Accept` messages since it has not received a quorum of `Promise` messages.
Hence, the expected output is `nil`.

Next, the Proposer receives a `Promise` message from Acceptor 0:

```go
Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
  {Slot: 2, Vrnd: 1, Vval: valTwo},
}}
```

In response to this `Promise` message, the Proposer should now send an `Accept` since it has received a quorum of `Promise` messages.
Since both acceptors have voted for a different value in slot 2, the Proposer should adopt the value with the highest `Vrnd`, which is Acceptor 0's vote for `valTwo`.
Hence, the expected output should be:

```go
[]Accept{{From: 2, Slot: 2, Rnd: 2, Val: valTwo}}
```

### Test Scenario 3

In this scenario, the Proposer first receives a `Promise` message from Acceptor 1:

```go
Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
  {Slot: 1, Vrnd: 0, Vval: valOne},
  {Slot: 2, Vrnd: 0, Vval: valTwo},
}}
```

In response to this `Promise` message, the Proposer should not send any `Accept` messages since it has not received a quorum of `Promise` messages.
Hence, the expected output is `nil`.

> Moreover, it is worth pointing out that the `PValue` for `Slot` 1 should be ignored.
> This is since it is lower than the `Prepare` message's `Slot`, which is 2, as specified at the start of the appendix.

Next, the Proposer receives a `Promise` message from Acceptor 0:

```go
Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
  {Slot: 2, Vrnd: 1, Vval: valTwo},
}}
```

In response to this `Promise` message, the Proposer should now send an `Accept` since it has received a quorum of `Promise` messages.
Since both acceptors have voted for a different value in slot 2, the Proposer should adopt the value with the highest `Vrnd`, which is Acceptor 0's vote for `valTwo`.
Hence, the expected output should be:

```go
[]Accept{{From: 2, Slot: 2, Rnd: 2, Val: valTwo}}
```

### Test Scenario 4

In this scenario, the Proposer first receives a `Promise` message from Acceptor 1 who has already voted for a value in slots 2 and 4.
Note the gap; there is no vote in slot 3:

```go
Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
  {Slot: 2, Vrnd: 1, Vval: valOne},
  {Slot: 4, Vrnd: 1, Vval: valThree},
}}
```

In response to this `Promise` message, the Proposer should not send any `Accept` messages since it has not received a quorum of `Promise` messages.
Hence, the expected output is `nil`.

Next, the Proposer receives a `Promise` message from Acceptor 0 who has only voted for a value in slot 2:

```go
Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
  {Slot: 2, Vrnd: 1, Vval: valOne},
}}
```

In response to this `Promise` message, the Proposer should now send an `Accept` since it has received a quorum of `Promise` messages.
With these `Promise` messages, the Proposer is locked in on slots 2 and 4.
To guarantee that the `Accept` messages are processed in increasing consecutive slot order, the gap in slot 3 must be filled with a no-op value.
This is accomplished by creating an `Accept` message with an empty `Value{}`.
Hence, the expected output should be:

```go
[]Accept{
  {From: 2, Slot: 2, Rnd: 2, Val: valOne},
  {From: 2, Slot: 3, Rnd: 2, Val: Value{}},
  {From: 2, Slot: 4, Rnd: 2, Val: valThree},
}
```

### Test Scenario 5

In this scenario, the Proposer first receives a `Promise` message from Acceptor 1 who has already voted for a value in slots 2, 4, and 5.
Note the gap; there is no vote in slot 3:

```go
Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
  {Slot: 2, Vrnd: 0, Vval: valOne},
  {Slot: 4, Vrnd: 0, Vval: valThree},
  {Slot: 5, Vrnd: 1, Vval: valTwo},
}}
```

In response to this `Promise` message, the Proposer should not send any `Accept` messages since it has not received a quorum of `Promise` messages.
Hence, the expected output is `nil`.

Next, the Proposer receives another identical `Promise` message from Acceptor 1.

```go
Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
  {Slot: 2, Vrnd: 0, Vval: valOne},
  {Slot: 4, Vrnd: 0, Vval: valThree},
  {Slot: 5, Vrnd: 1, Vval: valTwo},
}}
```

In response to this `Promise` message, the Proposer should not send any `Accept` messages since this is a duplicate input message, and we have not yet received a quorum of `Promise` messages from different acceptors.

Next, the Proposer receives a `Promise` message from Acceptor 0 who has voted for a value in slots 2 and 4:

```go
Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
  {Slot: 2, Vrnd: 0, Vval: valOne},
  {Slot: 4, Vrnd: 1, Vval: valTwo},
}}
```

We have now received a quorum of `Promise` messages, and the Proposer should now send `Accept` messages.
With these `Promise` messages, the Proposer is locked in on slots 2, 4, and 5.
To guarantee that the `Accept` messages are processed in increasing consecutive slot order, the gap in slot 3 must be filled with a no-op value.
This is accomplished by creating an `Accept` message with an empty `Value{}`.
Hence, the expected output should be:

```go
[]Accept{
  {From: 2, Slot: 2, Rnd: 2, Val: valOne},
  {From: 2, Slot: 3, Rnd: 2, Val: Value{}},
  {From: 2, Slot: 4, Rnd: 2, Val: valTwo},
  {From: 2, Slot: 5, Rnd: 2, Val: valTwo},
}
```

Note that `valTwo` was chosen for slot 4 since it has the highest `Vrnd` of the two votes for slot 4.
