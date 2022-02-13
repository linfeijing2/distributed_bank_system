#!/bin/sh
pkill python3
pkill mp1_node
rm node2_small.txt
rm node2err_small.txt
sleep 5
python3 -u gentx.py 0.5 | ./mp1_node node2 1234 config2_small.txt 1>node2_small.txt 2>node2err_small.txt &
CR=$!
sleep 100
kill $CR

