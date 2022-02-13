#!/bin/sh
pkill python3
pkill mp1_node
rm node1_small.txt
rm node1err_small.txt
sleep 5
python3 -u gentx.py 0.5 | ./mp1_node node1 1234 config1_small.txt 1>node1_small.txt 2>node1err_small.txt &
CR=$!
sleep 100
kill $CR

