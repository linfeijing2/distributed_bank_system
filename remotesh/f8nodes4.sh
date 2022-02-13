#!/bin/sh
pkill python3
pkill mp1_node
rm fnode4.txt
rm fnode4err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node4 1234 config4.txt 1>fnode4.txt 2>fnode4err.txt &
CR=$!
sleep 200
kill $CR

