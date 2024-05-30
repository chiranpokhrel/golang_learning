# Answers to Paxos Questions

Please answer the questions below by editing this file.

1. Is it possible that Paxos enters an infinite loop? Explain.

   > Your answer here.

2. Is the value to agree on included in the `Prepare` message?

   > Your answer here.

3. Does Paxos rely on an increasing proposal/round number in order to work? Explain.

   > Your answer here.

4. Consider this description for Phase 1B:
   If the proposal number _N_ is higher than any previous proposals, then each Acceptor promises not to accept any proposals less than _N_.
   The Acceptor does this by sending the value it last accepted for this instance to the Proposer.

   What is meant by "the value it last accepted"?
   And what is an "instance" in this case?

   > Your answer here.
   > Your answer here.

5. Explain, with an example, what will happen if there are multiple proposers.

   > Your answer here.

6. What happens if two proposers both believe themselves to be the leader and send `Prepare` messages simultaneously?

   > Your answer here.

7. What can we say about system synchrony if there are multiple proposers (or leaders)?

   > Your answer here.

8. Can an acceptor accept one value in round 1 and another value in round 2? Explain.

   > Your answer here.

9. What must happen for a value to be "chosen"?
   What is the connection between chosen values and learned values?

   > Your answer here.
   > Your answer here.
