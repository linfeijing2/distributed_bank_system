#!/bin/sh
pkill python3
pkill mp1_node
rm fnode5.txt
rm fnode5err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node5 1234 config5.txt 1>fnode5.txt 2>fnode5err.txt &
CR=$!
sleep 200
kill $CR

