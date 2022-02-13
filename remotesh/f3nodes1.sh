#!/bin/sh
pkill python3
pkill mp1_node
rm fnode1_small.txt
rm fnode1err_small.txt
sleep 5
python3 -u gentx.py 0.5 | ./mp1_node node1 1234 config1_small.txt 1>fnode1_small.txt 2>fnode1err_small.txt &
CR=$!
sleep 100
kill $CR

