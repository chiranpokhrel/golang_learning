package leaderdetector

// Definitions - DO NOT EDIT

// UnknownID represents an unknown id.
const UnknownID int = -1

// LeaderDetector is the interface implemented by a leader detector.
type LeaderDetector interface {
	Leader() int
	Subscribe() <-chan int
}
