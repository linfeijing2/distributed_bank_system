#!/bin/sh
pkill python3
pkill mp1_node
rm node4.txt
rm node4err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node4 1234 config4.txt 1>node4.txt 2>node4err.txt &
CR=$!
sleep 100
kill $CR

