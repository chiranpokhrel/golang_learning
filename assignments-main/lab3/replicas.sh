#!/bin/bash

set -e

# The first IP:port is the address of the local replica, while the rest are the addresses of the other replicas.
bin/replica -addrs="localhost:50081,localhost:50083,localhost:50082" &
bin/replica -addrs="localhost:50082,localhost:50081,localhost:50083" &
bin/replica -addrs="localhost:50083,localhost:50081,localhost:50082" &

pid=$!
sleep 3
kill -9 $pid
sleep 5
killall replica
