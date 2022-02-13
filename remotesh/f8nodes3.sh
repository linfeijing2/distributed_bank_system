#!/bin/sh
pkill python3
pkill mp1_node
rm fnode3.txt
rm fnode3err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node3 1234 config3.txt 1>fnode3.txt 2>fnode3err.txt &
CR=$!
sleep 100
kill $CR

