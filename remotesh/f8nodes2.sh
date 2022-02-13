#!/bin/sh
pkill python3
pkill mp1_node
rm fnode2.txt
rm fnode2err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node2 1234 config2.txt 1>fnode2.txt 2>fnode2err.txt &
CR=$!
sleep 100
kill $CR

