#!/bin/sh
pkill python3
pkill mp1_node
rm node6.txt
rm node6err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node6 1234 config6.txt 1>node6.txt 2>node6err.txt &
CR=$!
sleep 100
kill $CR

