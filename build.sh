#!/bin/sh
export GOPATH=`pwd`
cd src
DIRS="centerserver
	gameserver
	misc"

for dir in $DIRS
do
    cd $dir
    echo building $dir...
    rm -f $dir
    go build
    cd ..
done
