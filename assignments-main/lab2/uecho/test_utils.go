package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/sync/errgroup"
)

const connectionTimeout = 3 * time.Second

func waitForServer(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("resolve error: %v", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("dial error: %v", err)
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(connectionTimeout))
	if err != nil {
		return fmt.Errorf("connection deadline error: %v", err)
	}

	log.Println("waiting connection...")
	g := new(errgroup.Group)
	g.Go(func() error {
		for {
			var buf [512]byte
			_, err := conn.Write([]byte("ping"))
			if err != nil {
				return fmt.Errorf("write error: %v", err)
			}
			if _, err = conn.Read(buf[0:]); err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			return nil
		}
	})
	return g.Wait()
}
