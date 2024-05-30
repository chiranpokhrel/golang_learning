#! /bin/bash
set -e

./paxosserver -laddr="localhost:50081" -addrs="localhost:50083,localhost:50082" &
./paxosserver -laddr="localhost:50082" -addrs="localhost:50081,localhost:50083" &
./paxosserver -laddr="localhost:50083" -addrs="localhost:50081,localhost:50082" &

echo "running, enter to stop"
read && killall paxosserver
