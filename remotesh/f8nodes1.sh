#!/bin/sh
pkill python3
pkill mp1_node
rm fnode1.txt
rm fnode1err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node1 1234 config1.txt 1>fnode1.txt 2>fnode1err.txt &
CR=$!
sleep 100
kill $CR

