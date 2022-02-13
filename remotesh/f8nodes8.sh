#!/bin/sh
pkill python3
pkill mp1_node
rm fnode8.txt
rm fnode8err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node8 1234 config8.txt 1>fnode8.txt 2>fnode8err.txt &
CR=$!
sleep 200
kill $CR

