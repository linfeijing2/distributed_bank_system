#!/bin/sh
pkill python3
pkill mp1_node
rm node8.txt
rm node8err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node8 1234 config8.txt 1>node8.txt 2>node8err.txt &
CR=$!
sleep 100
kill $CR

