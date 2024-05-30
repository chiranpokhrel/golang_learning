package leaderdetector

import "testing"

func TestCorrectLeaderSetAfterInit(t *testing.T) {
	tests := []struct {
		name       string
		nodes      []int
		wantLeader int
	}{
		{name: "Empty", nodes: []int{}, wantLeader: UnknownID},
		{name: "One", nodes: []int{0}, wantLeader: 0},
		{name: "Two", nodes: []int{0, 1}, wantLeader: 1},
		{name: "Three", nodes: []int{0, 1, 2}, wantLeader: 2},
		{name: "FourA", nodes: []int{0, 1, 2, 3, 4}, wantLeader: 4},
		{name: "FourB", nodes: []int{0, 1, 2, 3, 42}, wantLeader: 42},
		{name: "Five", nodes: []int{0, 1, 2, 3, 4, 5}, wantLeader: 5},
		{name: "SingleNegative", nodes: []int{-1}, wantLeader: UnknownID},
		{name: "TwoNegative", nodes: []int{-2, -3}, wantLeader: UnknownID},
		{name: "OneNegativeTwoPositive", nodes: []int{0, -1, 2}, wantLeader: 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ld := NewMonLeaderDetector(test.nodes)
			gotLeader := ld.Leader()
			if gotLeader != test.wantLeader {
				t.Errorf("Leader() = %d, want %d", gotLeader, test.wantLeader)
			}
		})
	}
}

func TestCorrectLeaderAfterSuspectAndRestoreSequence(t *testing.T) {
	tests := []struct {
		name       string
		nodes      []int
		events     []event
		wantLeader int
	}{
		{name: "Empty", nodes: []int{}, events: []event{}, wantLeader: UnknownID},
		{name: "One0", nodes: []int{0}, events: []event{}, wantLeader: 0},
		{name: "OneS0", nodes: []int{0}, events: []event{{S, 0}}, wantLeader: UnknownID},
		{name: "Two0", nodes: []int{0, 1}, events: []event{}, wantLeader: 1},
		{name: "TwoS0", nodes: []int{0, 1}, events: []event{{S, 0}}, wantLeader: 1},
		{name: "TwoS1", nodes: []int{0, 1}, events: []event{{S, 1}}, wantLeader: 0},
		{name: "TwoS01", nodes: []int{0, 1}, events: []event{{S, 0}, {S, 1}}, wantLeader: UnknownID},
		{name: "TwoS01R0", nodes: []int{0, 1}, events: []event{{S, 0}, {S, 1}, {R, 0}}, wantLeader: 0},
		{name: "TwoS01R1", nodes: []int{0, 1}, events: []event{{S, 0}, {S, 1}, {R, 1}}, wantLeader: 1},
		{name: "TwoS01R01", nodes: []int{0, 1}, events: []event{{S, 0}, {S, 1}, {R, 0}, {R, 1}}, wantLeader: 1},
		{name: "TwoS01R10", nodes: []int{0, 1}, events: []event{{S, 0}, {S, 1}, {R, 1}, {R, 0}}, wantLeader: 1},
		{name: "Three0", nodes: []int{0, 1, 2}, events: []event{}, wantLeader: 2},
		{name: "ThreeS0", nodes: []int{0, 1, 2}, events: []event{{S, 0}}, wantLeader: 2},
		{name: "ThreeS1", nodes: []int{0, 1, 2}, events: []event{{S, 1}}, wantLeader: 2},
		{name: "ThreeS2", nodes: []int{0, 1, 2}, events: []event{{S, 2}}, wantLeader: 1},
		{name: "ThreeS01", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}}, wantLeader: 2},
		{name: "ThreeS02", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 2}}, wantLeader: 1},
		{name: "ThreeS12", nodes: []int{0, 1, 2}, events: []event{{S, 1}, {S, 2}}, wantLeader: 0},
		{name: "ThreeS0R0", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {R, 0}}, wantLeader: 2},
		{name: "ThreeS1R1", nodes: []int{0, 1, 2}, events: []event{{S, 1}, {R, 1}}, wantLeader: 2},
		{name: "ThreeS2R2", nodes: []int{0, 1, 2}, events: []event{{S, 2}, {R, 2}}, wantLeader: 2},
		{name: "ThreeS012", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}}, wantLeader: UnknownID},
		{name: "ThreeS001", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 0}, {S, 1}}, wantLeader: 2},
		{name: "ThreeS2R2S0", nodes: []int{0, 1, 2}, events: []event{{S, 2}, {R, 2}, {S, 0}}, wantLeader: 2},
		{name: "ThreeR22S0", nodes: []int{0, 1, 2}, events: []event{{R, 2}, {R, 2}, {R, 0}}, wantLeader: 2},
		{name: "ThreeR222", nodes: []int{0, 1, 2}, events: []event{{R, 2}, {R, 2}, {R, 2}}, wantLeader: 2},
		{name: "ThreeS0011", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 0}, {S, 1}, {S, 1}}, wantLeader: 2},
		{name: "ThreeS00111", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 0}, {S, 1}, {S, 1}, {S, 1}}, wantLeader: 2},
		{name: "ThreeS22000", nodes: []int{0, 1, 2}, events: []event{{S, 2}, {S, 2}, {S, 0}, {S, 0}, {S, 0}}, wantLeader: 1},
		{name: "ThreeS012R0", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}, {R, 0}}, wantLeader: 0},
		{name: "ThreeS012R1", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}, {R, 1}}, wantLeader: 1},
		{name: "ThreeS012R2", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}, {R, 2}}, wantLeader: 2},
		{name: "ThreeS012R01", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}, {R, 0}, {R, 1}}, wantLeader: 1},
		{name: "ThreeS012R02", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}, {R, 0}, {R, 2}}, wantLeader: 2},
		{name: "ThreeS012R12", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}, {R, 1}, {R, 2}}, wantLeader: 2},
		{name: "ThreeS012R012", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}, {R, 0}, {R, 1}, {R, 2}}, wantLeader: 2},
		{name: "ThreeS012R021", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}, {R, 0}, {R, 2}, {R, 1}}, wantLeader: 2},
		{name: "ThreeS012R102", nodes: []int{0, 1, 2}, events: []event{{S, 0}, {S, 1}, {S, 2}, {R, 1}, {R, 0}, {R, 2}}, wantLeader: 2},
		{name: "ThreeR2S2R1S0", nodes: []int{0, 1, 2}, events: []event{{R, 2}, {S, 2}, {R, 1}, {S, 0}}, wantLeader: 1},
		{name: "ThreeR122_S42_R14_S0", nodes: []int{0, 1, 2}, events: []event{{R, 122}, {S, 42}, {R, 14}, {S, 0}}, wantLeader: 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ld := NewMonLeaderDetector(test.nodes)
			for _, event := range test.events {
				if event.eType == S {
					ld.Suspect(event.id)
				} else {
					ld.Restore(event.id)
				}
			}
			gotLeader := ld.Leader()
			if gotLeader != test.wantLeader {
				t.Errorf("Leader() = %d, want %d", gotLeader, test.wantLeader)
			}
		})
	}
}

func TestPublishSubscribe(t *testing.T) {
	// Sequence of suspect/restore events, and wanted leader publications.
	eventSequence := []eventPubSub{
		{desc: "Suspect 2, want publish for 1", eType: S, id: 2, wantOutput: Yes, wantLeader: 1},
		{desc: "Restore 2, want publish for 2", eType: R, id: 2, wantOutput: Yes, wantLeader: 2},
		{desc: "Restore 2, no change -> no output", eType: R, id: 2, wantOutput: No},
		{desc: "Suspect 1, no change -> no output", eType: S, id: 1, wantOutput: No},
		{desc: "Suspect 0, no change -> no output", eType: S, id: 0, wantOutput: No},
		{desc: "Suspect 2, all suspected, want publish for UnknownID", eType: S, id: 2, wantOutput: Yes, wantLeader: UnknownID},
		{desc: "Restore 0, want publish for 0", eType: R, id: 0, wantOutput: Yes, wantLeader: 0},
		{desc: "Restore 1, want publish for 1", eType: R, id: 1, wantOutput: Yes, wantLeader: 1},
		{desc: "Restore 2, want publish for 2", eType: R, id: 2, wantOutput: Yes, wantLeader: 2},
	}
	tests := []struct {
		name        string
		nodes       []int
		subscribers int
		events      []eventPubSub
	}{
		{name: "Empty", nodes: []int{}, subscribers: 0, events: []eventPubSub{}},
		{name: "ThreeNodes1Sub", nodes: []int{0, 1, 2}, subscribers: 1, events: eventSequence},
		{name: "ThreeNodes2Sub", nodes: []int{0, 1, 2}, subscribers: 2, events: eventSequence},
		{name: "ThreeNodes3Sub", nodes: []int{0, 1, 2}, subscribers: 3, events: eventSequence},
		{name: "ThreeNodes5Sub", nodes: []int{0, 1, 2}, subscribers: 5, events: eventSequence},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ld := NewMonLeaderDetector(test.nodes)
			subscribers := make([]<-chan int, test.subscribers)
			for n := 0; n < len(subscribers); n++ {
				subscribers[n] = ld.Subscribe()
				if subscribers[n] == nil {
					t.Fatalf("Subscriber[%d] = <nil> channel, want non-nil channel", n+1)
				}
			}
			for _, event := range test.events {
				if event.eType == S {
					ld.Suspect(event.id)
				} else {
					ld.Restore(event.id)
				}
				for k, subscriber := range subscribers {
					if event.wantOutput {
						select {
						case gotLeader := <-subscriber:
							// Got publication, check if leader is correct.
							if gotLeader != event.wantLeader {
								// Got publication for wrong leader.
								t.Errorf("Subscriber[%d]: got publication for leader %d, want leader %d (%s)", k+1, gotLeader, event.wantLeader, event.desc)
							}
						default:
							// We want publication, but got none.
							t.Errorf("Subscriber[%d]: got no publication, want one for leader %d (%s)", k+1, event.wantLeader, event.desc)
						}
					} else {
						select {
						case gotLeader := <-subscriber:
							// Got publication, want none.
							t.Errorf("Subscriber[%d]: got publication for leader %d, want no publication (%s)", k+1, gotLeader, event.desc)
						default:
							// Got no publication, and want none.
						}
					}
				}
			}
		})
	}
}
