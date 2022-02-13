#!/bin/sh
pkill python3
pkill mp1_node
rm node3_small.txt
rm node3err_small.txt
sleep 5
python3 -u gentx.py 0.5 | ./mp1_node node3 1234 config3_small.txt 1>node3_small.txt 2>node3err_small.txt &
CR=$!
sleep 100
kill $CR

