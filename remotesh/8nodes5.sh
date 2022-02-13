#!/bin/sh
pkill python3
pkill mp1_node
rm node5.txt
rm node5err.txt
sleep 5
python3 -u gentx.py 5 | ./mp1_node node5 1234 config5.txt 1>node5.txt 2>node5err.txt &
CR=$!
sleep 100
kill $CR

