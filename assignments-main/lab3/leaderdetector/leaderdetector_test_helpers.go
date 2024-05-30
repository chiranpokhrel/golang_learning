// Test helper functions - DO NOT EDIT

package leaderdetector

type eventType bool

const (
	S   = true
	R   = false
	Yes = true
	No  = false
)

type event struct {
	eType eventType // S -> suspect, R -> restore
	id    int       // node id to suspect or restore
}

type eventPubSub struct {
	desc       string
	eType      eventType // S -> suspect, R -> restore
	id         int       // node id to suspect or restore
	wantOutput bool      // should the event produce output to subscribers
	wantLeader int       // what leader id should publication contain
}
