#!/bin/sh
pkill python3
pkill mp1_node
rm node3.txt
rm node3err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node3 1234 config3.txt 1>node3.txt 2>node3err.txt &
CR=$!
sleep 100
kill $CR

