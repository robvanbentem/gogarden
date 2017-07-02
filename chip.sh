#!/bin/zsh


HOST=chip

echo "building arm binary"
env GOOS=linux GOARCH=arm go build

echo "removing old file"
ssh chip@$HOST "pkill -9 gogarden"

echo "copying new file"
scp ./gogarden chip@$HOST:/home/chip/gogarden

echo "running new version"
ssh chip@$HOST ./gogarden
