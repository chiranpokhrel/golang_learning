# Changes to Lab 5

Below are the diffs for two of the files that you probably have edited already.
These are the important changes to note when merging the updated version of the lab.
There are many other changes to the tests and other files in the lab that should not impact your work.

Additionally, some of the tests have been renamed to better match the naming convention of the other labs.
The lab description have been updated to reflect this change.

Of particular note is that the tests (by default) no longer output debug information to the console.
However, it is easy to enable debug output by running the tests, e.g.,:

```console
cd lab5/gorumspaxos
LOG=1 go test -v -run TestFiveReplicas
LOG=1 go test -v -run TestLeaderFailure
```

The reason for this change is that the debug output can be quite verbose and make it harder to see the actual test results, especially on QuickFeed.
We recommend that you only enable debug output when you need it.

```diff
diff --git a/lab5/gorumspaxos/proposer.go b/lab5/gorumspaxos/proposer.go
index 03d3e1a..7185dba 100644
--- a/lab5/gorumspaxos/proposer.go
+++ b/lab5/gorumspaxos/proposer.go
@@ -98,6 +98,14 @@ func (p *Proposer) runMultiPaxos() {
 	// TODO(student): complete
 }

+// nextAcceptMsg returns the next accept message to be sent, if any.
+// If there are no pending accept messages or any client requests to process,
+// it returns nil.
+func (p *Proposer) nextAcceptMsg() (accept *pb.AcceptMsg) {
+	// TODO(student): complete
+	return accept
+}
+
 // Perform the accept quorum call on the replicas.
 //
 //  1. Check if any pending accept requests in the acceptReqQueue to process
@@ -105,7 +113,7 @@ func (p *Proposer) runMultiPaxos() {
 //  3. Increment the nextSlot and prepare an accept message for the pending request,
 //     using crnd and nextSlot.
 //  4. Perform accept quorum call on the configuration and return the learnMsg.
-func (p *Proposer) performAccept() (*pb.LearnMsg, error) {
+func (p *Proposer) performAccept(accept *pb.AcceptMsg) (*pb.LearnMsg, error) {
 	// TODO(student): complete
 	return nil, nil
 }
```

```diff
diff --git a/lab5/gorumspaxos/replica.go b/lab5/gorumspaxos/replica.go
index 715d21d..5489fe3 100644
--- a/lab5/gorumspaxos/replica.go
+++ b/lab5/gorumspaxos/replica.go
@@ -42,6 +42,7 @@ type PaxosReplica struct {
 	srv             *gorums.Server          // the gorums.Server that the replica is registered to
 	stop            chan struct{}           // channel for stopping the replica's run loop.
 	learntVal       map[uint32]*pb.LearnMsg // Stores all received learn messages
+	stopped         bool
 }

 // NewPaxosReplica returns a new Paxos replica with a nodeMap configuration.
@@ -91,6 +92,10 @@ func newTestReplicaLeader() *PaxosReplica {

 // Stops the failure detector, replica, and gorums server.
 func (r *PaxosReplica) Stop() {
+	if r.stopped {
+		return
+	}
+	r.stopped = true
 	r.failureDetector.Stop()
 	r.stop <- struct{}{} // stop the replica's run loop
 	r.fdManager.Close()
@@ -124,6 +129,7 @@ func (r *PaxosReplica) run() {
 		// TODO(student) create Paxos configuration and set the proposer's configuration
 		// TODO(student) implement Paxos replica run loop
 		_ = trustMsgs // TODO: remove this line when you start using trustMsgs
+		<-r.stop      // TODO: remove this line when you implement the method
 	}()

 	go func() {
@@ -193,3 +199,20 @@ func (r *PaxosReplica) ClientHandle(ctx gorums.ServerCtx, req *pb.Value) (rsp *p
 	// TODO(student) complete
 	return nil, errors.New("unable to get the response")
 }
+
+// remainingResponses returns the number of responses that are still pending.
+func (r *PaxosReplica) remainingResponses() int {
+	r.mu.Lock()
+	defer r.mu.Unlock()
+	// TODO(student) complete
+	return 0
+}
+
+// responseIDs returns the IDs of the responses that are still pending.
+func (r *PaxosReplica) responseIDs() []uint64 {
+	r.mu.Lock()
+	defer r.mu.Unlock()
+	ids := make([]uint64, 0)
+	// TODO(student) complete
+	return ids
+}
```
