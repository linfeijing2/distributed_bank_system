#!/bin/sh
pkill python3
pkill mp1_node
rm node7.txt
rm node7err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node7 1234 config7.txt 1>node7.txt 2>node7err.txt &
CR=$!
sleep 100
kill $CR

