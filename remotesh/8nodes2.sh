#!/bin/sh
pkill python3
pkill mp1_node
rm node2.txt
rm node2err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node2 1234 config2.txt 1>node2.txt 2>node2err.txt &
CR=$!
sleep 100
kill $CR

