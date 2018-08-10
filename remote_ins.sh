#!/bin/bash
scp /home/hlp/workspace/quorum/src/github.com/ethereum/go-ethereum/build/bin/geth root@192.168.3.30:/root/quorum/quorum
scp root@192.168.3.30:/root/quorum/quorum root@192.168.3.31:/root/quorum/quorum
scp root@192.168.3.30:/root/quorum/quorum root@192.168.3.32:/root/quorum/quorum


