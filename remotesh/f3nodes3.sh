#!/bin/sh
pkill python3
pkill mp1_node
rm fnode3_small.txt
rm fnode3err_small.txt
sleep 5
python3 -u gentx.py 0.5 | ./mp1_node node3 1234 config3_small.txt 1>fnode3_small.txt 2>fnode3err_small.txt &
CR=$!
sleep 200
kill $CR

