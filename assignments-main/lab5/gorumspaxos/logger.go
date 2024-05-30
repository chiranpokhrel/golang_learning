package gorumspaxos

import (
	"fmt"
	"log"
	"os"
)

func (r *PaxosReplica) Logf(format string, a ...any) {
	if os.Getenv("LOG") != "" {
		log.Printf("Node %d\t%s", r.id, fmt.Sprintf(format, a...))
	}
}

func (p *Proposer) Logf(format string, a ...any) {
	if os.Getenv("LOG") != "" {
		log.Printf("Node %d\t%s", p.id, fmt.Sprintf(format, a...))
	}
}
