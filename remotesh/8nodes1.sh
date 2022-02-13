#!/bin/sh
pkill python3
pkill mp1_node
rm node1.txt
rm node1err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node1 1234 config1.txt 1>node1.txt 2>node1err.txt &
CR=$!
sleep 100
kill $CR

