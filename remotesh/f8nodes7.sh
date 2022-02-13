#!/bin/sh
pkill python3
pkill mp1_node
rm fnode7.txt
rm fnode7err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node7 1234 config7.txt 1>fnode7.txt 2>fnode7err.txt &
CR=$!
sleep 200
kill $CR

