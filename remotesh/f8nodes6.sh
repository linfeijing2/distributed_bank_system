#!/bin/sh
pkill python3
pkill mp1_node
rm fnode6.txt
rm fnode6err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node6 1234 config6.txt 1>fnode6.txt 2>fnode6err.txt &
CR=$!
sleep 200
kill $CR

