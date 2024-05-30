#! /bin/bash
set -e

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083" -clientRequest="M1,M2,M3,M4" -clientId "1" &
./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083" -clientRequest="M5,M6,M7,M8" -clientId "2" &
./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083" -clientRequest="M9,M10,M11,M12" -clientId "3" &

read && killall paxosclient
